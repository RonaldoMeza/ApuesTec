# Fase 4 - Backend base real

## Objetivo de la fase

Construir la base tecnica real del backend de ApuesTec en Go + Gin, conectandolo con PostgreSQL y Redis mediante variables de entorno, respuestas estandar, middlewares iniciales, rutas base y apagado ordenado.

No se implemento autenticacion funcional, salas, invitaciones, predicciones, puntajes, rankings ni frontend funcional.

## Estructura backend creada o ajustada

```text
backend/
  cmd/
    api/
      main.go
  internal/
    config/
      config.go
    database/
      postgres.go
    redis/
      redis.go
    middleware/
      cors.go
      logging.go
      recovery.go
      request_id.go
      security.go
    response/
      response.go
    routes/
      router.go
  Dockerfile
  .env.example
```

La estructura mantiene espacio para arquitectura feature-based en fases posteriores, sin crear modulos funcionales completos todavia.

## Variables de entorno usadas

Variables requeridas al iniciar:

- `APP_NAME`
- `APP_ENV`
- `APP_PORT`
- `API_PREFIX`
- `APP_PUBLIC_URL`
- `FRONTEND_URL`
- `DATABASE_URL`
- `REDIS_URL`
- `JWT_ACCESS_SECRET`
- `JWT_REFRESH_SECRET`

Variables leidas con valores por defecto seguros si faltan:

- `CORS_ALLOWED_ORIGINS`: usa `FRONTEND_URL` si esta vacia.
- `JWT_ACCESS_TTL`: `1h`.
- `JWT_REFRESH_TTL`: `720h`.
- `PASSWORD_HASH_ALGO`: `bcrypt`.
- `BCRYPT_COST`: `12`.
- `GOOGLE_CLIENT_ID`: vacia.
- `GOOGLE_CLIENT_SECRET`: vacia.
- `GOOGLE_REDIRECT_URL`: vacia.
- `COOKIE_SECURE`: `false`.
- `COOKIE_SAME_SITE`: `Lax`.
- `LOG_LEVEL`: `info`.

No se hardcodearon valores sensibles en codigo. `DATABASE_URL` debe seguir usando `postgres:5432` dentro de Docker y `REDIS_URL` debe seguir usando `redis:6379`.

## Librerias agregadas

- `github.com/jackc/pgx/v5/pgxpool`: pool y conexion PostgreSQL.
- `github.com/redis/go-redis/v9`: cliente Redis y validacion con `PING`.

## Conexion a PostgreSQL

El paquete `internal/database` usa `DATABASE_URL`, configura un pool basico con `pgxpool`, valida conexion con `Ping` al iniciar y expone cierre mediante `pool.Close()` desde `main.go`.

El backend no ejecuta migraciones, no crea tablas y no modifica el modelo de datos.

## Conexion a Redis

El paquete `internal/redis` usa `REDIS_URL`, crea un cliente `go-redis`, valida conexion con `PING` al iniciar y se cierra durante el apagado ordenado.

No se implemento cache, rankings ni rate limiting funcional.

## Respuestas estandar

El paquete `internal/response` define helpers reutilizables:

- Respuesta exitosa: `success`, `data`, `message`.
- Respuesta de error: `success`, `error.code`, `error.message`.

Los errores de dependencias no exponen errores crudos de PostgreSQL ni Redis al cliente.

## Middlewares implementados

- Request ID con header `X-Request-ID`.
- Logging basico por request.
- Recovery para panics con respuesta estandar controlada.
- CORS usando `CORS_ALLOWED_ORIGINS`.
- Headers basicos de seguridad: `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`.

No se implemento autenticacion JWT, autorizacion por roles ni rate limiting funcional.

## Rutas implementadas

- `GET /api/v1/health`: responde estado simple del servicio.
- `GET /api/v1/health/dependencies`: valida PostgreSQL y Redis.

Respuesta esperada de dependencias:

```json
{
  "success": true,
  "data": {
    "postgres": "ok",
    "redis": "ok"
  },
  "message": "dependencies ok"
}
```

Si una dependencia falla, responde `503` con error controlado `DEPENDENCY_UNAVAILABLE`.

## Graceful shutdown

`cmd/api/main.go` escucha `SIGINT` y `SIGTERM`, apaga el servidor HTTP con timeout y cierra conexiones de PostgreSQL y Redis.

## Como validar health checks

```powershell
Invoke-WebRequest -UseBasicParsing -Uri http://localhost:8081/api/v1/health
```

## Como validar dependencias

```powershell
Invoke-WebRequest -UseBasicParsing -Uri http://localhost:8081/api/v1/health/dependencies
```

## Comandos ejecutados

```powershell
go get github.com/jackc/pgx/v5/pgxpool github.com/redis/go-redis/v9
gofmt -w .\cmd\api\main.go .\internal\config\config.go .\internal\database\postgres.go .\internal\redis\redis.go .\internal\response\response.go .\internal\middleware\request_id.go .\internal\middleware\logging.go .\internal\middleware\recovery.go .\internal\middleware\cors.go .\internal\middleware\security.go .\internal\routes\router.go
go mod tidy
go test ./...
docker compose config --quiet
docker compose up --build -d
docker compose ps
docker compose logs backend
```

## Problemas encontrados y solucion

- El backend anterior solo tenia un `main.go` minimo con `GET /api/v1/health`. Se separo la base tecnica en paquetes internos sin crear modulos funcionales fuera de fase.
- La fase requiere conectividad real con PostgreSQL y Redis. Se agregaron librerias mantenibles y validacion de conexion al inicio.
- Para evitar exponer errores internos, el health de dependencias responde estados controlados.

## Resultado de validacion

- `go test ./...`: correcto, todos los paquetes compilan.
- `docker compose config --quiet`: correcto, sin salida.
- `docker compose up --build -d`: correcto, backend y frontend reconstruidos; servicios iniciados.
- `docker compose ps`: correcto; backend, frontend, Nginx, PostgreSQL y Redis quedaron `healthy`.
- `GET http://localhost:8081/api/v1/health`: correcto, respondio `200` con `{"service":"apuestec-backend","status":"ok"}`.
- `GET http://localhost:8081/api/v1/health/dependencies`: correcto, respondio `200` con PostgreSQL y Redis en `ok`.
- `docker compose logs backend`: correcto; muestra rutas registradas, servidor escuchando en `8080` y logs de requests con `request_id`.
- Nginx sigue siendo el unico puerto publico de la aplicacion en `http://localhost:8081`.
- Backend sigue sin publicarse al host, solo `8080/tcp` interno.

## Checklist final

- [x] Backend compila correctamente.
- [x] Backend inicia en Docker.
- [x] Backend lee variables desde `backend/.env`.
- [x] Backend conecta correctamente con PostgreSQL usando `DATABASE_URL`.
- [x] Backend conecta correctamente con Redis usando `REDIS_URL`.
- [x] `GET /api/v1/health` responde `ok`.
- [x] `GET /api/v1/health/dependencies` responde estado de PostgreSQL y Redis.
- [x] Nginx sigue siendo el unico punto publico en `http://localhost:8081`.
- [x] Backend sigue interno en puerto `8080`.
- [x] No se implemento autenticacion funcional.
- [x] No se implementaron salas, invitaciones, predicciones, puntajes ni rankings.
- [x] No se cambiaron reglas funcionales del plan.
- [x] No se agregaron secretos reales.
- [x] No se ejecutan migraciones desde el backend.
