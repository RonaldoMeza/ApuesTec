# ApuesTec

ApuesTec es una plataforma educativa de predicciones deportivas del Mundial. El sistema permite registrar predicciones de marcadores, acumular puntos, consultar rankings, crear salas privadas y usar invitaciones temporales, sin dinero real ni logica de apuestas monetarias.

## Stack tecnologico

Frontend:

- Next.js.
- React.
- TypeScript.
- TailwindCSS.
- shadcn/ui.
- Recharts.

Backend:

- Go.
- Gin Gonic.

Persistencia e infraestructura:

- PostgreSQL.
- Redis.
- Docker y Docker Compose.
- Nginx.
- k6.

## Estructura del monorepo

```text
ApuesTec/
  backend/
  frontend/
  database/
  nginx/
  tests/
    k6/
  docker-compose.yml
  README.md
  AGENTS.md
  PLAN_PROYECTO_APUESTEC.md
```

## Reglas principales

- No usar dinero real, pagos, cuotas reales, casas de apuestas ni logica de apuestas monetarias.
- No usar Firebase Auth.
- No usar API externa de partidos en el MVP.
- El backend siempre debe emitir sus propios JWT.
- Google OAuth solo valida identidad.
- El refresh token no debe guardarse en `localStorage`.
- PostgreSQL es la fuente oficial de verdad.
- Redis es una capa de cache y optimizacion.
- Mantener arquitectura feature-based en backend y frontend.
- No colocar logica de negocio en handlers.
- No acceder directamente a PostgreSQL ni Redis desde handlers.
- No subir secretos ni archivos `.env` reales.

## Ejecucion local con Docker Compose

Docker Compose levanta el monorepo con los servicios base de infraestructura para desarrollo local. Nginx es el unico punto de entrada publico; el backend y Redis quedan disponibles solo dentro de la red interna de Docker. PostgreSQL mantiene su comunicacion interna por `postgres:5432` y expone `localhost:5433` unicamente para administracion local en desarrollo.

Antes de levantar el entorno, crea los archivos locales de variables desde las plantillas:

```bash
cp database/.env.example database/.env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env.local
```

En Windows PowerShell:

```powershell
Copy-Item database/.env.example database/.env
Copy-Item backend/.env.example backend/.env
Copy-Item frontend/.env.example frontend/.env.local
```

Los archivos `.env` reales no deben commitearse. Las plantillas `.env.example` son la referencia documentada para crear credenciales locales.

Levantar el entorno:

```bash
docker compose up --build
```

Detener y eliminar contenedores de la ejecucion local:

```bash
docker compose down
```

Consultar estado de servicios:

```bash
docker compose ps
```

Validar escalamiento horizontal del backend:

```bash
docker compose up --scale backend=3
```

## Puertos

- `http://localhost:8081`: punto de entrada principal por Nginx.
- `frontend:3000`: puerto interno del frontend Next.js.
- `backend:8080`: puerto interno de la API Go/Gin.
- `postgres:5432`: puerto interno de PostgreSQL.
- `localhost:5433`: puerto local de administracion de PostgreSQL para desarrollo, por ejemplo pgAdmin instalado en la maquina host.
- `redis:6379`: puerto interno de Redis.

## Servicios levantados

- `frontend`: aplicacion base Next.js para desarrollo local.
- `backend`: API Go/Gin base con health checks y validacion de dependencias internas.
- `postgres`: base de datos principal de ApuesTec con volumen persistente.
- `redis`: cache y capa de optimizacion con volumen persistente.
- `nginx`: reverse proxy y punto de entrada principal.
- `k6`: servicio preparado bajo perfil `tools` para pruebas posteriores.

## Persistencia y cache

PostgreSQL es la fuente oficial de verdad de ApuesTec. Redis se usa como cache y optimizacion para futuras funciones como rankings, rate limiting e invitaciones temporales; los datos criticos deben poder reconstruirse desde PostgreSQL.

## Health checks del backend

Nginx es el unico punto de entrada publico. Valida el backend desde la raiz del proyecto con:

```powershell
docker compose up --build -d
Invoke-WebRequest -UseBasicParsing -Uri http://localhost:8081/api/v1/health
Invoke-WebRequest -UseBasicParsing -Uri http://localhost:8081/api/v1/health/dependencies
docker compose logs backend
```

Endpoints disponibles:

- `GET http://localhost:8081/api/v1/health`: estado basico del servicio.
- `GET http://localhost:8081/api/v1/health/dependencies`: valida conexion interna a PostgreSQL y Redis.

El backend no ejecuta migraciones ni crea tablas. Las migraciones y seeds siguen siendo manuales desde `database/scripts/`.

## Inicializacion de base de datos local

Ejecuta estos comandos desde la raiz del proyecto. Los contenedores deben estar levantados y los scripts requieren ejecucion manual; Docker Compose no aplica migraciones ni seeds automaticamente.

```powershell
docker compose up --build -d
.\database\scripts\apply-migrations.ps1
.\database\scripts\apply-seeds.ps1
.\database\scripts\check-tables.ps1
```

No se deben commitear archivos `.env` reales. Nginx sigue siendo el punto de entrada publico en `http://localhost:8081`.

## Migraciones y seeds de base de datos

La estructura inicial de PostgreSQL vive en `database/`:

- `database/migrations/001_init_schema.up.sql`: crea el esquema inicial.
- `database/migrations/001_init_schema.down.sql`: revierte el esquema inicial.
- `database/seeds/001_seed_roles.sql`: inserta roles globales minimos.
- `database/seeds/002_seed_sample_teams_matches.sql`: inserta equipos y partidos de prueba local.

Ejecutar migracion y seeds en PowerShell:

```powershell
docker compose up -d postgres
.\database\scripts\apply-migrations.ps1
.\database\scripts\apply-seeds.ps1
```

Validar tablas:

```powershell
.\database\scripts\check-tables.ps1
```

Resetear la base local de desarrollo, eliminando y recreando tablas y datos locales:

```powershell
.\database\scripts\reset-database.ps1
```

## Conexion local desde pgAdmin

Para inspeccionar PostgreSQL desde pgAdmin instalado en la maquina host, crea un servidor con estos parametros:

- Host: `localhost`
- Port: `5433`
- Database: valor de `POSTGRES_DB` en `database/.env`
- User: valor de `POSTGRES_USER` en `database/.env`
- Password: valor de `POSTGRES_PASSWORD` en `database/.env`

El puerto `5433` existe solo para administracion local en desarrollo. La aplicacion dentro de Docker no lo usa: el backend debe seguir conectandose a PostgreSQL mediante `DATABASE_URL` con `postgres:5432`. Redis no se expone al host y Nginx sigue siendo el unico punto de entrada publico por `http://localhost:8081`.

Revertir migracion:

```powershell
Get-Content -Raw database/migrations/001_init_schema.down.sql | docker compose exec -T postgres psql -U apuestec_user -d apuestec
```

Mas detalles operativos estan en `database/README.md` y `docs/fases/FASE_03_BASE_DATOS.md`.

## Variables de entorno

Docker Compose carga variables desde estos archivos locales:

- `database/.env`: credenciales y base inicial de PostgreSQL, basado en `database/.env.example`.
- `backend/.env`: configuracion de la API, URLs internas de PostgreSQL/Redis, JWT y OAuth, basado en `backend/.env.example`.
- `frontend/.env.local`: variables publicas de Next.js, basado en `frontend/.env.example`.

Variables clave por capa:

- Database: `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_PORT`.
- Backend: `APP_NAME`, `APP_ENV`, `APP_PORT`, `API_PREFIX`, `APP_PUBLIC_URL`, `FRONTEND_URL`, `CORS_ALLOWED_ORIGINS`, `DATABASE_URL`, `REDIS_URL`, `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `PASSWORD_HASH_ALGO`, `BCRYPT_COST`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL`, `COOKIE_SECURE`, `COOKIE_SAME_SITE`, `LOG_LEVEL`.
- Frontend: `NEXT_PUBLIC_APP_NAME`, `NEXT_PUBLIC_APP_URL`, `NEXT_PUBLIC_API_BASE_URL`, `NEXT_PUBLIC_GOOGLE_CLIENT_ID`.

No se deben commitear archivos `.env` reales ni secretos. El backend usa `postgres` y `redis` como hosts internos de Docker en `DATABASE_URL` y `REDIS_URL`; `DATABASE_URL` debe mantenerse apuntando a `postgres:5432` dentro de la red interna.

## k6

La carpeta `tests/k6/` queda preparada para scripts de pruebas de estres en fases posteriores. El servicio puede invocarse con perfil de herramientas cuando existan scripts:

```bash
docker compose --profile tools run --rm k6 version
```

## Nota de desarrollo

Todo desarrollo debe seguir `PLAN_PROYECTO_APUESTEC.md` y las instrucciones de `AGENTS.md`. Antes de avanzar de fase, deben validarse los criterios de aceptacion definidos en el plan.
