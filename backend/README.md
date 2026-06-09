# ApuesTec — Backend

API REST de ApuesTec escrita en Go con Gin Gonic.

## Stack

- Go + Gin Gonic
- PostgreSQL (pgx v5)
- Redis (go-redis)
- JWT (golang-jwt v5, HS256)
- bcrypt para hash de contraseñas

## Estructura (feature-based)

```
backend/
  cmd/api/main.go                  Punto de entrada
  internal/
    config/                        Carga de variables de entorno
    database/                      Conexión PostgreSQL
    redis/                         Conexión Redis
    response/                      Formato estándar de respuesta
    routes/                        Router y middlewares globales
    auth/                          Autenticación (handler, service, repository, middleware)
    users/                         CRUD de usuarios
    roles/                         Asignación de roles
    audit/                         Registro de eventos de auditoría
    health/                        Health checks
```

## Endpoints

| Método | Ruta                          | Auth     | Descripción                     |
|--------|-------------------------------|----------|---------------------------------|
| POST   | /api/v1/auth/register         | No       | Registro local                  |
| POST   | /api/v1/auth/login            | No       | Inicio de sesión                |
| POST   | /api/v1/auth/refresh          | No       | Rotación de refresh token       |
| POST   | /api/v1/auth/logout           | No       | Revocación de refresh token     |
| GET    | /api/v1/auth/me               | Bearer   | Información del usuario actual  |
| POST   | /api/v1/auth/google           | No       | Autenticación con Google        |
| POST   | /api/v1/auth/change-password  | Bearer   | Cambio de contraseña            |
| GET    | /api/v1/health                | No       | Health check básico             |
| GET    | /api/v1/health/dependencies   | No       | Health check con dependencias   |

## Formato de respuesta

Todas las respuestas usan un envoltorio estándar:

```json
// Éxito
{ "success": true, "data": { ... }, "message": "..." }

// Error
{ "success": false, "error": { "code": "...", "message": "..." } }
```

## Ejecución

### Con Docker

```bash
docker compose up --build -d
```

### Local (desarrollo rápido)

Requiere PostgreSQL, Redis y Nginx corriendo en Docker:

```powershell
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
cd backend
.\run.ps1
```

El script `run.ps1` carga variables desde `.env` y sobrescribe conexiones para apuntar a localhost.

## Variables de entorno

Ver `backend/.env` y `backend/.env.example`. Variables principales:

- `DATABASE_URL` — Conexión a PostgreSQL
- `REDIS_URL` — Conexión a Redis
- `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` — Secretos JWT
- `CORS_ALLOWED_ORIGINS` — Orígenes permitidos
- `LOGIN_MAX_ATTEMPTS` / `LOGIN_LOCK_MINUTES` — Bloqueo por intentos
