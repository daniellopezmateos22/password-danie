# password-danie 🔐

Este es mi proyecto para la **prueba técnica**.  
La idea es construir una aplicación completa con:

- **Backend en Go** (con Gin + Gorm + JWT + Postgres).
- **Frontend en React + Vite** servido con **Nginx**.
- **Base de datos Postgres** en contenedor.
- Todo orquestado con **Docker Compose**.

---

## 🚀 Cómo levantar la app

Primero clonar el repositorio:

```bash
git clone https://github.com/daniellopezmateos22/password-danie.git
cd password-danie
```

### 1. Configurar variables de entorno
Creo un archivo `.env` en la raíz con algo así:

```env
JWT_SECRET=super-secret-change-me
DB_DSN=host=db user=postgres password=postgres dbname=vault sslmode=disable port=5432 TimeZone=UTC
ENC_KEY_BASE64=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

> ⚠️ La clave `ENC_KEY_BASE64` debe ser de **32 bytes en base64**.

### 2. Levantar con Docker Compose

```bash
docker compose up --build
```

Esto levanta 3 servicios:
- `db` → PostgreSQL
- `api` → servidor Go en http://localhost:8080
- `web` → frontend en http://localhost:3000

---

## 📡 Endpoints principales (API)

- `GET /health` → chequeo de estado.
- `POST /auth/register` → registrar usuario.
- `POST /auth/login` → login, devuelve token JWT.
- `POST /auth/forgot` y `POST /auth/reset` → flujo básico de recuperación de contraseña.
- CRUD completo de **vault** en `/api/vault`:
  - `GET /api/vault`
  - `POST /api/vault`
  - `GET /api/vault/:id`
  - `PATCH /api/vault/:id`
  - `DELETE /api/vault/:id`

---

## 🛠️ Stack

- Go 1.23 + Gin
- Gorm (ORM) + PostgreSQL
- JWT para autenticación
- React + Vite en frontend
- Nginx para servir estáticos
- Docker + Docker Compose

---

## ✨ Autor

Hecho por **Daniel López Mateos** ✌️
