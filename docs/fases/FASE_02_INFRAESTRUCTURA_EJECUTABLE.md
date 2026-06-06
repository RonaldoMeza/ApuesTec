# Fase 2 - Infraestructura ejecutable

## Objetivo de la fase

Dejar el monorepo de ApuesTec ejecutable localmente con Docker Compose, preparando los servicios base para el desarrollo posterior sin implementar autenticacion, salas, predicciones, puntuacion ni pantallas funcionales completas.

## Servicios configurados

- `frontend`: servicio Next.js base ejecutado en modo desarrollo dentro de Docker.
- `backend`: servicio Go/Gin minimo con health check `GET /api/v1/health`.
- `postgres`: PostgreSQL 16 con base inicial `apuestec`, variables desde `database/.env`, volumen persistente y healthcheck.
- `redis`: Redis 7 con append-only file, volumen persistente y healthcheck.
- `nginx`: reverse proxy y unico punto de entrada externo por `localhost:8081`.
- `k6`: servicio preparado con perfil `tools` para ejecucion posterior de pruebas de estres.

## Archivos creados o modificados

- `docker-compose.yml`: definicion de servicios, `env_file`, healthchecks, red interna implicita y volumenes persistentes.
- `backend/Dockerfile`: imagen de desarrollo para ejecutar la API Go.
- `backend/.dockerignore`: exclusiones para evitar archivos locales innecesarios en la imagen.
- `backend/cmd/api/main.go`: API minima con endpoint de health.
- `frontend/Dockerfile`: imagen de desarrollo para ejecutar Next.js.
- `frontend/.dockerignore`: exclusiones para evitar `node_modules`, `.next` y caches locales en la imagen.
- `nginx/nginx.conf`: configuracion inicial de reverse proxy hacia frontend y backend.
- `README.md`: instrucciones de ejecucion local, puertos, servicios y persistencia/cache.
- `docs/fases/FASE_02_INFRAESTRUCTURA_EJECUTABLE.md`: documentacion de esta fase.

## Variables de entorno usadas

Docker Compose carga los archivos locales `database/.env`, `backend/.env` y `frontend/.env.local`. Cada archivo real debe crearse desde su plantilla `.env.example` y no debe commitearse.

PostgreSQL, plantilla `database/.env.example` y archivo local `database/.env`:

```env
POSTGRES_DB=apuestec
POSTGRES_USER=apuestec_user
POSTGRES_PASSWORD=apuestec_password
POSTGRES_PORT=5432
```

Backend, plantilla `backend/.env.example` y archivo local `backend/.env`:

```env
APP_NAME=apuestec-backend
APP_ENV=development
APP_PORT=8080
API_PREFIX=/api/v1
APP_PUBLIC_URL=http://localhost:8081
FRONTEND_URL=http://localhost:8081
CORS_ALLOWED_ORIGINS=http://localhost:8081
DATABASE_URL=postgres://apuestec_user:apuestec_password@postgres:5432/apuestec?sslmode=disable
REDIS_URL=redis://redis:6379/0
JWT_ACCESS_SECRET=change_me_access_secret
JWT_REFRESH_SECRET=change_me_refresh_secret
JWT_ACCESS_TTL=1h
JWT_REFRESH_TTL=720h
PASSWORD_HASH_ALGO=bcrypt
BCRYPT_COST=12
GOOGLE_CLIENT_ID=change_me_google_client_id
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8081/api/v1/auth/google/callback
COOKIE_SECURE=false
COOKIE_SAME_SITE=Lax
LOG_LEVEL=debug
```

Frontend, plantilla `frontend/.env.example` y archivo local `frontend/.env.local`:

```env
NEXT_PUBLIC_APP_NAME=ApuesTec
NEXT_PUBLIC_APP_URL=http://localhost:8081
NEXT_PUBLIC_API_BASE_URL=/api/v1
NEXT_PUBLIC_GOOGLE_CLIENT_ID=change_me_google_client_id
```

`NODE_ENV=development` se define directamente en el servicio `frontend` de `docker-compose.yml`.

## Comandos para ejecutar

Levantar servicios base:

```bash
docker compose up --build
```

Detener servicios:

```bash
docker compose down
```

Consultar estado:

```bash
docker compose ps
```

Validar escalamiento del backend:

```bash
docker compose up --scale backend=3
```

Ejecutar k6 preparado para fases posteriores:

```bash
docker compose --profile tools run --rm k6 version
```

## Comandos para verificar

Verificar Nginx:

```bash
curl http://localhost:8081/health
```

Verificar backend a traves de Nginx:

```bash
curl http://localhost:8081/api/v1/health
```

Verificar estado de contenedores:

```bash
docker compose ps
```

## Problemas encontrados y solucion aplicada

- `docker-compose.yml`, `backend/Dockerfile` y `nginx/nginx.conf` estaban como placeholders de Fase 1. Se reemplazaron por configuraciones ejecutables de Fase 2.
- `docs/fases/FASE_02_INFRAESTRUCTURA_EJECUTABLE.md` existia vacio. Se completo con la documentacion operativa de la fase.
- El backend no tenia punto de entrada ejecutable. Se agrego un `main.go` minimo con Gin y un health check basico permitido por el plan.
- El frontend no tenia Dockerfile. Se agrego un Dockerfile de desarrollo local sin crear pantallas funcionales nuevas.
- El puerto `8080` del host no esta disponible en la maquina local. Se publico Nginx en `localhost:8081` manteniendo el puerto interno `80` del contenedor.
- El contexto de build del frontend incluia archivos locales innecesarios. Se agrego `frontend/.dockerignore` y el contexto bajo de cientos de MB a unos KB en la reconstruccion.
- Docker Compose fue alineado para cargar variables reales desde `database/.env`, `backend/.env` y `frontend/.env.local`, dejando los `.env.example` como plantillas sin secretos.

## Resultado de validacion

- `docker compose config`: correcto.
- `docker compose --profile tools config`: correcto, incluye servicio `k6` preparado.
- `go test ./...` en `backend/`: correcto, sin archivos de test todavia.
- `npm run lint` en `frontend/`: correcto.
- `docker compose up --build -d`: correcto.
- `docker compose ps`: correcto, servicios `frontend`, `backend`, `postgres`, `redis` y `nginx` en estado healthy.
- `docker compose up -d --scale backend=3`: correcto, backend escalo a 3 replicas healthy.
- `http://localhost:8081/health`: correcto.
- `http://localhost:8081/api/v1/health`: correcto.

## Checklist de validacion

- [x] Existe `backend/`.
- [x] Existe `frontend/`.
- [x] Existe `database/`.
- [x] Existe `nginx/`.
- [x] Existe `tests/k6/`.
- [x] Existe `docker-compose.yml`.
- [x] Existe `README.md`.
- [x] Existe `AGENTS.md`.
- [x] Existe `PLAN_PROYECTO_APUESTEC.md`.
- [x] Docker Compose define `frontend`, `backend`, `postgres`, `redis`, `nginx` y `k6`.
- [x] Docker Compose carga `env_file` por capa sin commitear secretos reales.
- [x] PostgreSQL tiene base inicial, volumen persistente y healthcheck.
- [x] Redis tiene volumen persistente y healthcheck.
- [x] Backend tiene Dockerfile funcional y health check minimo.
- [x] Frontend tiene Dockerfile funcional.
- [x] Nginx es el unico punto de entrada externo.
- [x] Backend no se expone directamente al host.
- [x] README.md queda actualizado.
- [x] No se implementa logica funcional compleja.
- [x] No se cambia el stack tecnologico.
- [x] No se usan secretos reales.
