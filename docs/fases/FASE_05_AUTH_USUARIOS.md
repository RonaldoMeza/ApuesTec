# Fase 5: Autenticacion y Usuarios

## Objetivo

Implementar autenticacion completa con registro local, login, JWT, refresh tokens,
bloqueo por intentos fallidos, roles globales, integracion con Google OAuth y auditoria.

## Stack

- Backend: Go + Gin + pgx v5 + golang-jwt v5
- Paswords: bcrypt via golang.org/x/crypto
- Base de datos: PostgreSQL (tablas existentes)
- Cache: Redis (disponible para rate limiting futuro)

## Arquitectura

Se siguio el patron Feature-Based Architecture con capas:

```
internal/auth/
  ├── dto.go          # Request/Response DTOs y errores del dominio
  ├── repository.go   # Operaciones DB (refresh_tokens, auth_accounts)
  ├── service.go      # Logica de negocio (registro, login, refresh, etc.)
  ├── handler.go      # Handlers HTTP
  ├── middleware.go    # JWT auth middleware + role guard
  └── auth_test.go    # Tests unitarios

internal/users/
  ├── dto.go          # UserResponse DTO
  └── repository.go   # CRUD usuarios (create, findByEmail, etc.)

internal/roles/
  └── repository.go   # CRUD roles y asignaciones

internal/audit/
  └── repository.go   # Registro de auditoria
```

## Endpoints implementados

| Method | Path                          | Auth     | Descripcion                          |
|--------|-------------------------------|----------|--------------------------------------|
| POST   | /api/v1/auth/register         | No       | Registro local de usuario            |
| POST   | /api/v1/auth/login            | No       | Login con email y password           |
| POST   | /api/v1/auth/refresh          | No       | Rotacion de refresh token            |
| POST   | /api/v1/auth/logout           | No       | Revocacion de refresh token          |
| GET    | /api/v1/auth/me               | Bearer   | Informacion del usuario autenticado  |
| POST   | /api/v1/auth/google           | No       | Autenticacion con Google ID token    |
| POST   | /api/v1/auth/change-password  | Bearer   | Cambio de contrasena                 |

## Flujo de autenticacion

1. **Registro**: El usuario envia `fullName`, `email`, `password`. Se hashea la password con bcrypt,
   se crea el usuario en `users`, se asigna rol `USER`, se genera par de tokens JWT.
2. **Login**: Verifica credenciales, chequea bloqueo, incrementa intentos fallidos si es necesario,
   genera tokens JWT al exito.
3. **Refresh**: Recibe refresh token, lo hashea (SHA-256), busca en DB, verifica vigencia,
   revoca el anterior (rotacion), emite nuevo par.
4. **Logout**: Hashea el refresh token recibido, lo revoca en DB.
5. **Google Auth**: Recibe `idToken` de Google, valida claims (aud, iss, exp), busca o crea
   auth_account y user, emite tokens JWT del backend.
6. **Change Password**: Verifica password actual, hashea la nueva, actualiza en DB,
   revoca TODOS los refresh tokens activos del usuario.

## Tokens

### Access Token (JWT)
- Duracion: configurable via `JWT_ACCESS_TTL` (default 1h)
- Claims: `sub` (userID), `email`, `roles`, `iat`, `exp`
- Algoritmo: HS256
- Secreto: `JWT_ACCESS_SECRET`

### Refresh Token
- Token opaco de 32 bytes aleatorios en hex (64 caracteres)
- Se almacena hasheado con SHA-256 en `refresh_tokens.token_hash`
- Duracion: configurable via `JWT_REFRESH_TTL` (default 30 dias)
- Rotacion: cada refresh revoca el anterior y emite uno nuevo
- Revocacion manual via logout o cambio de password

## Bloqueo por intentos fallidos

- Configurable via `LOGIN_MAX_ATTEMPTS` (default 5) y `LOGIN_LOCK_MINUTES` (default 15)
- Al alcanzar el maximo de intentos, se setea `users.locked_until`
- El login verifica `locked_until > NOW()` antes de procesar
- Si el bloqueo expiro, se resetean los intentos automaticamente

## Roles

Tres roles globales predefinidos en la semilla (`001_seed_roles.sql`):

| Rol          | Descripcion                                   |
|--------------|-----------------------------------------------|
| SUPER_ADMIN  | Permisos administrativos completos            |
| ADMIN        | Administracion operativa del MVP              |
| USER         | Rol base para usuarios registrados (default)  |

Al registrarse, todo usuario recibe automaticamente el rol `USER`.

## Auditoria

Eventos registrados en `audit_logs`:

- `USER_REGISTERED` - registro exitoso
- `USER_LOGGED_IN` - login exitoso
- `LOGIN_FAILED` - intento fallido de login
- `USER_LOCKED` - usuario bloqueado por intentos
- `TOKEN_REFRESHED` - refresh token rotado
- `USER_LOGGED_OUT` - logout
- `PASSWORD_CHANGED` - cambio de contrasena
- `GOOGLE_AUTH_SUCCESS` - autenticacion google exitosa
- `GOOGLE_AUTH_FAILED` - autenticacion google fallida

## Variables de entorno nuevas

```
LOGIN_MAX_ATTEMPTS=5        # Intentos fallidos antes de bloquear
LOGIN_LOCK_MINUTES=15       # Duracion del bloqueo en minutos
```

Las variables `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `BCRYPT_COST`,
`JWT_ACCESS_TTL`, `JWT_REFRESH_TTL` y `GOOGLE_CLIENT_ID` ya existian
en la configuracion de Fase 4.

## Archivos creados

```
backend/internal/auth/dto.go
backend/internal/auth/repository.go
backend/internal/auth/service.go
backend/internal/auth/handler.go
backend/internal/auth/middleware.go
backend/internal/auth/auth_test.go
backend/internal/roles/repository.go
backend/internal/audit/repository.go
backend/internal/users/dto.go
backend/internal/users/repository.go
```

## Archivos modificados

```
backend/internal/config/config.go          # LoginMaxAttempts, LoginLockMinutes
backend/internal/routes/router.go          # Registro de rutas de auth
backend/.env                               # LOGIN_MAX_ATTEMPTS, LOGIN_LOCK_MINUTES
backend/.env.example                       # Nuevas variables documentadas
backend/go.mod                             # golang-jwt/jwt/v5 agregado
backend/go.sum                             # Actualizado
```

## Validacion

```bash
go test ./...                    # Tests unitarios
go build ./...                   # Compilacion
docker compose config --quiet    # Configuracion valida
```
