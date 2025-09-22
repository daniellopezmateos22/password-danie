# password-danie 🔐

Gestor de contraseñas estilo LastPass / 1Password.  
Implementado como prueba técnica usando **Go (Gin)** en el backend y **React (Vite)** en el frontend, todo dentro de **Docker Compose**.

---

## 🚀 Tecnologías

- **Backend**: Go 1.23 + Gin
  - SQLite (persistencia)
  - JWT (autenticación)
  - Bcrypt (hash de contraseñas de usuario)
  - AES-256 (cifrado de contraseñas en el vault)
  - Migraciones automáticas
  - Tests unitarios e integración (E2E con modernc.org/sqlite)
- **Frontend**: React + Vite + TypeScript
  - Pantalla de Login/Registro
  - Pantalla de Vault con CRUD (Create, Read, Update, Delete)
- **Infraestructura**:
  - Docker + Docker Compose
  - Variables de entorno en `.env`
  - Despliegue en **GitHub Codespaces**

---

## 📂 Estructura del proyecto

```plaintext
password-danie/
├── backend/                 # Backend en Go
│   ├── cmd/server/main.go   # Entry point (composition root)
│   ├── internal/
│   │   ├── http/            # Rutas y controladores (Gin)
│   │   ├── dto/             # DTOs de requests/responses
│   │   ├── middleware/      # Middleware (AuthRequired, JWT)
│   │   ├── repository/      # Interfaces de repositorios
│   │   │   └── sqlite/      # Implementaciones SQLite
│   │   ├── security/        # JWT utils, AES helpers
│   │   └── usecase/         # Casos de uso (Auth, Vault, ResetPassword)
│   ├── pkg/db/              # Conexión SQLite + migraciones
│   ├── migrations/          # Migraciones SQL
│   ├── go.mod / go.sum
│   └── Dockerfile
│
├── frontend/                # Frontend en React + Vite
│   ├── src/
│   │   ├── pages/           # AuthPage, VaultPage
│   │   ├── api.ts           # Cliente HTTP con fetch
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── vite.config.ts
│   ├── package.json
│   └── Dockerfile
│
├── docker-compose.yml       # Orquesta backend + frontend
├── .env.example             # Variables de entorno
└── README.md
```

---

## ⚙️ Variables de entorno

Ejemplo `.env` en la raíz del repo:

```env
PORT=8080
SQLITE_DSN=data/app.db
JWT_SECRET=change-me
AES_KEY=0123456789abcdef0123456789abcdef
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h
```
```Ejemplo .env.local en frontend/ (solo para Codespaces/local dev):
VITE_API_BASE_URL=http://localhost:8080
En Codespaces se recomienda generar dinámicamente esta variable:
VITE_API_BASE_URL=https://${CODESPACE_NAME}-8080.app.github.dev
```
---

## ▶️ Levantar el proyecto

Clonar repo y arrancar con Docker Compose:

```bash
git clone https://github.com/daniellopezmateos22/password-danie.git
cd password-danie
cp .env.example .env

docker compose build --no-cache
docker compose up -d
```

- **Backend**: http://localhost:8080  
- **Frontend**: http://localhost:5173  

---
## 🚀 Despliegue en GitHub Codespaces

Este proyecto está preparado para ejecutarse directamente en Codespaces.

Abre el repo en Codespaces (Open with Codespaces en GitHub).

Una vez dentro, crea los archivos de entorno:

backend/.env → usa el ejemplo de arriba (PORT, JWT_SECRET, AES_KEY, etc).

frontend/.env.local → apunta a tu API pública de Codespaces:

---

## 🎥 Demo en Video

Mira la demo completa en YouTube, donde se prueban los **tests E2E** y luego el uso del **frontend (login, registro y CRUD del vault)**:

👉 [Ver demo en YouTube](https://www.youtube.com/watch?v=cIQzVgFrfSk)

---

## 🔑 Endpoints principales (API REST)

### Auth
- `POST /api/v1/auth/register` → Crear usuario
- `POST /api/v1/auth/login` → Login y obtener JWT
- `POST /api/v1/auth/reset/request` → Solicitar reset password
- `POST /api/v1/auth/reset/confirm` → Confirmar reset password

### Users
- `GET /api/v1/users/me` → Info del usuario (JWT requerido)

### Vault
- `GET /api/v1/vault/entries` → Listar contraseñas (con búsqueda `q`, filtrado por dominio `domain`, paginación)
- `GET /api/v1/vault/entries/:id` → Obtener por ID (**vista detallada de una contraseña**)
- `POST /api/v1/vault/entries` → Crear nueva entrada
- `PUT /api/v1/vault/entries/:id` → Actualizar entrada
- `DELETE /api/v1/vault/entries/:id` → Eliminar entrada

---

## 🧪 Tests

Tests unitarios + integración.

Ejecutar test E2E (Go + modernc.org/sqlite):

```bash
cd backend
go test ./internal/integration -run Test_FullAPI_HappyPath -v
```

Esto recorre **todos los endpoints**: health, ready, register, login, users/me, CRUD del vault, reset password.

---

## 🖥️ Frontend

  - Tras login → **VaultPage** con CRUD literal:
  - Crear secreto
  - Buscar/Listar
  - Filtrar por dominio
  - Ver detalle
  - Update
  - Delete


Configura el frontend con:

```env
# frontend/.env
VITE_API_BASE_URL=http://localhost:8080
```
---

## ✅ Checklist de requisitos del enunciado

| Requisito                             | Endpoint / Funcionalidad          | Estado |
|---------------------------------------|-----------------------------------|--------|
| Registro de usuario                   | POST /api/v1/auth/register        | ✅     |
| Login + JWT                           | POST /api/v1/auth/login           | ✅     |
| Ver usuario actual                    | GET /api/v1/users/me              | ✅     |
| Crear contraseña                      | POST /api/v1/vault/entries        | ✅     |
| Listar contraseñas                    | GET /api/v1/vault/entries         | ✅     |
| Vista detallada de una contraseña     | GET /api/v1/vault/entries/:id     | ✅     |
| Actualizar contraseña                 | PUT /api/v1/vault/entries/:id     | ✅     |
| Eliminar contraseña                   | DELETE /api/v1/vault/entries/:id  | ✅     |
| Búsqueda en contraseñas               | GET /api/v1/vault/entries?q=...   | ✅     |
| Filtrado por dominio                  | GET /api/v1/vault/entries?domain= | ✅     |
| Reset password (request + confirm)    | POST /api/v1/auth/reset/*         | ✅     |
| Seguridad (bcrypt + AES + JWT)        | Backend                           | ✅     |
| Frontend básico con CRUD              | React/Vite                        | ✅     |
| Infraestructura con Docker Compose    | docker-compose.yml                | ✅     |
| Tests automáticos end-to-end          | Test_FullAPI_HappyPath            | ✅     |

---

## 👨‍💻 Autor
**Daniel López Mateos**  



