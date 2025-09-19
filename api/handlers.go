package api

import (
	"strings"
	"time"
	"strconv" 

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"crypto/rand"
    "encoding/base64"
	
	"github.com/daniellopezmateos22/password-danie/internal/auth"
	"github.com/daniellopezmateos22/password-danie/internal/crypto"
	"github.com/daniellopezmateos22/password-danie/internal/models"
)

type vaultOut struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Username  string    `json:"username"`
    URL       string    `json:"url,omitempty"`
    Icon      string    `json:"icon,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func toVaultOut(v models.VaultItem) vaultOut {
    return vaultOut{
        ID:        v.ID,
        Title:     v.Title,
        Username:  v.Username,
        URL:       v.URL,
        Icon:      v.Icon,
        CreatedAt: v.CreatedAt,
        UpdatedAt: v.UpdatedAt,
    }
}

func Routes(r *gin.Engine, db *gorm.DB) {
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	authG := r.Group("/auth")
	{
		authG.POST("/register", registerHandler(db))
		authG.POST("/login", loginHandler(db))
		authG.POST("/forgot", forgotHandler(db))
		authG.POST("/reset", resetHandler(db))
	}

	api := r.Group("/api")
	api.Use(auth.Middleware())
	{
		api.GET("/vault", listVault(db))
		api.POST("/vault", createVault(db))
		api.GET("/vault/:id", getVault(db))
		api.GET("/vault/:id/reveal", revealVault(db)) 
		api.PATCH("/vault/:id", updateVault(db))
		api.DELETE("/vault/:id", deleteVault(db))
	}
}

// ======= AUTH =======

func registerHandler(db *gorm.DB) gin.HandlerFunc {
	type In struct{ Email, Password string }
	return func(c *gin.Context) {
		var in In
		if err := c.ShouldBindJSON(&in); err != nil || in.Email == "" || in.Password == "" {
			c.JSON(400, gin.H{"error": "email y password requeridos"})
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		u := models.User{Email: strings.ToLower(in.Email), PasswordHash: string(hash)}
		if err := db.Create(&u).Error; err != nil {
			c.JSON(409, gin.H{"error": "email ya registrado"})
			return
		}
		c.JSON(201, gin.H{"id": u.ID, "email": u.Email})
	}
}

func loginHandler(db *gorm.DB) gin.HandlerFunc {
	type In struct{ Email, Password string }
	return func(c *gin.Context) {
		var in In
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(400, gin.H{"error": "payload inválido"})
			return
		}
		var u models.User
		if err := db.Where("email = ?", strings.ToLower(in.Email)).First(&u).Error; err != nil {
			c.JSON(401, gin.H{"error": "credenciales inválidas"})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
			c.JSON(401, gin.H{"error": "credenciales inválidas"})
			return
		}
		token, _ := auth.MakeToken(u.ID)
		c.JSON(200, gin.H{"token": token})
	}
}

func forgotHandler(db *gorm.DB) gin.HandlerFunc {
	type In struct{ Email string }
	return func(c *gin.Context) {
		var in In
		if err := c.ShouldBindJSON(&in); err != nil || in.Email == "" {
			c.JSON(400, gin.H{"error": "email requerido"})
			return
		}
		var u models.User
		if err := db.Where("email = ?", strings.ToLower(in.Email)).First(&u).Error; err == nil {
			tok := models.ResetToken{
				UserID:    u.ID,
				Token:     RandToken(24),
				ExpiresAt: time.Now().Add(15 * time.Minute),
			}
			db.Create(&tok)
			c.JSON(200, gin.H{"message": "token generado", "reset_token": tok.Token})
			return
		}
		c.JSON(200, gin.H{"message": "si existe, se enviará un email"})
	}
}

func resetHandler(db *gorm.DB) gin.HandlerFunc {
	type In struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	return func(c *gin.Context) {
		var in In
		if err := c.ShouldBindJSON(&in); err != nil || in.Token == "" || in.NewPassword == "" {
			c.JSON(400, gin.H{"error": "token y password requeridos"})
			return
		}
		var rt models.ResetToken
		if err := db.Where("token = ? AND used = false", in.Token).First(&rt).Error; err != nil || time.Now().After(rt.ExpiresAt) {
			c.JSON(400, gin.H{"error": "token inválido o expirado"})
			return
		}
		var u models.User
		if err := db.First(&u, rt.UserID).Error; err != nil {
			c.JSON(400, gin.H{"error": "usuario no existe"})
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
		db.Model(&u).Update("password_hash", string(hash))
		db.Model(&rt).Update("used", true)
		c.JSON(200, gin.H{"message": "password actualizada"})
	}
}

// ======= VAULT =======

type vaultIn struct {
	Title    string `json:"title" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	URL      string `json:"url"`
	Notes    string `json:"notes"`
	Icon     string `json:"icon"`
}

func listVault(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        uid := c.GetUint("user_id")
        var items []models.VaultItem

        q := db.Where("user_id = ?", uid)

        if s := c.Query("q"); s != "" {
            like := "%" + s + "%"
            q = q.Where("title ILIKE ? OR url ILIKE ?", like, like)
        }
        if domain := c.Query("domain"); domain != "" {
            like := "%://" + domain + "%"
            q = q.Where("url ILIKE ?", like)
        }

        // paginación y orden
        limitStr := c.DefaultQuery("limit", "50")
        offsetStr := c.DefaultQuery("offset", "0")
        limit, _ := strconv.Atoi(limitStr)
        offset, _ := strconv.Atoi(offsetStr)
        if limit <= 0 || limit > 200 { limit = 50 }
        if offset < 0 { offset = 0 }

        sort := c.DefaultQuery("sort", "updated_desc")
        switch sort {
        case "updated_desc":
            q = q.Order("updated_at desc")
        case "created_desc":
            q = q.Order("created_at desc")
        default:
            q = q.Order("id desc")
        }

        if err := q.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
            c.JSON(500, gin.H{"error": "db error"})
            return
        }
        c.JSON(200, items)
    }
}



func createVault(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        uid := c.GetUint("user_id")
        var in vaultIn
        if err := c.ShouldBindJSON(&in); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }

        pc, err := crypto.Encrypt([]byte(in.Password))
        if err != nil {
            c.JSON(500, gin.H{"error": "encrypt password error"})
            return
        }
        var nc []byte
        if in.Notes != "" {
            nc, err = crypto.Encrypt([]byte(in.Notes))
            if err != nil {
                c.JSON(500, gin.H{"error": "encrypt notes error"})
                return
            }
        }

        item := models.VaultItem{
            UserID:    uid,
            Title:     in.Title,
            Username:  in.Username,
            PasswordC: pc,
            URL:       in.URL,
            NotesC:    nc,
            Icon:      in.Icon,
        }
        if err := db.Create(&item).Error; err != nil {
            c.JSON(500, gin.H{"error": "db error"})
            return
        }
        c.JSON(201, toVaultOut(item))
    }
}


func getVault(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        uid := c.GetUint("user_id")
        var item models.VaultItem
        if err := db.Where("user_id = ? AND id = ?", uid, c.Param("id")).First(&item).Error; err != nil {
            c.JSON(404, gin.H{"error": "no encontrado"})
            return
        }
        c.JSON(200, toVaultOut(item))
    }
}


func updateVault(db *gorm.DB) gin.HandlerFunc {
    type inPatch struct {
        Title    *string `json:"title"`
        Username *string `json:"username"`
        Password *string `json:"password"`
        URL      *string `json:"url"`
        Notes    *string `json:"notes"`
        Icon     *string `json:"icon"`
    }
    return func(c *gin.Context) {
        uid := c.GetUint("user_id")
        var item models.VaultItem
        if err := db.Where("user_id = ? AND id = ?", uid, c.Param("id")).First(&item).Error; err != nil {
            c.JSON(404, gin.H{"error": "no encontrado"})
            return
        }
        var in inPatch
        if err := c.ShouldBindJSON(&in); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        if in.Title != nil {
            item.Title = *in.Title
        }
        if in.Username != nil {
            item.Username = *in.Username
        }
        if in.Password != nil {
            pc, err := crypto.Encrypt([]byte(*in.Password))
            if err != nil {
                c.JSON(500, gin.H{"error": "encrypt password error"})
                return
            }
            item.PasswordC = pc
        }
        if in.URL != nil {
            item.URL = *in.URL
        }
        if in.Notes != nil {
            var nc []byte
            if *in.Notes != "" {
                var err error
                nc, err = crypto.Encrypt([]byte(*in.Notes))
                if err != nil {
                    c.JSON(500, gin.H{"error": "encrypt notes error"})
                    return
                }
            }
            item.NotesC = nc
        }
        if in.Icon != nil {
            item.Icon = *in.Icon
        }
        if err := db.Save(&item).Error; err != nil {
            c.JSON(500, gin.H{"error": "db error"})
            return
        }
        c.JSON(200, toVaultOut(item))
    }
}


func deleteVault(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		res := db.Where("user_id = ? AND id = ?", uid, c.Param("id")).Delete(&models.VaultItem{})
		if res.Error != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}
		if res.RowsAffected == 0 {
			c.JSON(404, gin.H{"error": "no encontrado"})
			return
		}
		c.JSON(200, gin.H{"deleted": c.Param("id")})
	}
}

func revealVault(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        uid := c.GetUint("user_id")
        var item models.VaultItem
        if err := db.Where("user_id = ? AND id = ?", uid, c.Param("id")).First(&item).Error; err != nil {
            c.JSON(404, gin.H{"error": "no encontrado"})
            return
        }
        pass, err := crypto.Decrypt(item.PasswordC)
        if err != nil {
            c.JSON(500, gin.H{"error": "decrypt error"})
            return
        }
        c.JSON(200, gin.H{"password": string(pass)})
    }
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandToken(n int) string {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        return base64.RawURLEncoding.EncodeToString([]byte(time.Now().String()))[:n]
    }
    return base64.RawURLEncoding.EncodeToString(b)[:n]
}
