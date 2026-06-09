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

## Fases implementadas

| Fase | Descripción                          |
|------|--------------------------------------|
| 1    | Estructura base del monorepo         |
| 2    | Infraestructura ejecutable (Docker)  |
| 3    | Base de datos (PostgreSQL, migraciones, seeds) |
| 4    | Backend base (Go/Gin, health, middlewares, Nginx) |
| 5    | Autenticación backend (registro, login, JWT, refresh, roles, bloqueo, Google OAuth, auditoría) |
| 6    | Frontend auth base (login, register, dashboard, perfil, landing pública, modo oscuro, protección de rutas, toast, avatar menú) |

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

## Ejecucion

### Opcion 1: Todo con Docker

Todos los servicios (frontend, backend, postgres, redis, nginx) corren dentro de Docker.

```powershell
# Copiar plantillas de entorno
Copy-Item database/.env.example database/.env
Copy-Item backend/.env.example backend/.env
Copy-Item frontend/.env.example frontend/.env.local

# Migraciones y seeds
docker compose up -d postgres
.\database\scripts\apply-migrations.ps1
.\database\scripts\apply-seeds.ps1

# Levantar todo
docker compose up --build -d

# Abrir http://localhost:8081
```

```bash
docker compose down   # Detener todo
docker compose ps     # Estado de servicios
```

### Opcion 2: Infraestructura en Docker + app en local (desarrollo rapido)

PostgreSQL, Redis y Nginx corren en Docker. Backend y frontend se ejecutan localmente
para iterar cambios rapidamente con hot reload.

```powershell
# 1. Iniciar infraestructura
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# 2. Backend local (PowerShell)
cd backend
.\run.ps1

# 3. Frontend local (otra terminal)
cd frontend
npm run dev

# Abrir http://localhost:8081 (via Nginx) o http://localhost:3000 (directo Next.js)
```

## Puertos

- `http://localhost:8081`: punto de entrada principal por Nginx.
- `http://localhost:3000`: frontend Next.js en modo local (sin Docker).
- `frontend:3000`: puerto interno del frontend Next.js (Docker).
- `backend:8080`: puerto interno de la API Go/Gin.
- `postgres:5432`: puerto interno de PostgreSQL.
- `localhost:5433`: puerto local de administracion de PostgreSQL para desarrollo.
- `redis:6379`: puerto de Redis (interno de Docker y expuesto al host).

## Servicios

- `postgres`: base de datos principal con volumen persistente.
- `redis`: cache y capa de optimizacion con volumen persistente.
- `nginx`: reverse proxy, unico punto de entrada publico en `http://localhost:8081`.
- `backend`: API Go/Gin con autenticacion JWT, registro, login, roles y auditoria.
- `frontend`: Next.js con auth, landing, dashboard, perfil y proteccion de rutas.
- `k6`: pruebas de estres bajo perfil `tools`.

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

Cada capa tiene su archivo de entorno basado en `.env.example`:

- `database/.env` — Credenciales de PostgreSQL
- `backend/.env` — Configuracion de API, JWT, OAuth, conexiones a PostgreSQL/Redis
- `frontend/.env.local` — Variables publicas de Next.js

Variables clave por capa:

| Capa     | Variable                          | Descripcion                            |
|----------|-----------------------------------|----------------------------------------|
| Database | `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` | Credenciales de base de datos |
| Backend  | `DATABASE_URL`                    | Conexion a PostgreSQL                  |
| Backend  | `REDIS_URL`                       | Conexion a Redis                       |
| Backend  | `JWT_ACCESS_SECRET`               | Secreto para firmar access tokens      |
| Backend  | `JWT_REFRESH_SECRET`              | Secreto para firmar refresh tokens     |
| Backend  | `CORS_ALLOWED_ORIGINS`            | Origenes permitidos para CORS          |
| Backend  | `LOGIN_MAX_ATTEMPTS`              | Intentos fallidos antes de bloqueo     |
| Backend  | `LOGIN_LOCK_MINUTES`              | Duracion del bloqueo                   |
| Frontend | `NEXT_PUBLIC_API_URL`             | URL base de la API                     |
| Frontend | `NEXT_PUBLIC_GOOGLE_CLIENT_ID`    | Client ID de Google OAuth              |

No se deben commitear archivos `.env` reales ni secretos. Las plantillas `.env.example`
son la referencia documentada para crear credenciales locales.

## k6

La carpeta `tests/k6/` queda preparada para scripts de pruebas de estres en fases posteriores. El servicio puede invocarse con perfil de herramientas cuando existan scripts:

```bash
docker compose --profile tools run --rm k6 version
```

## Nota de desarrollo

Todo desarrollo debe seguir `PLAN_PROYECTO_APUESTEC.md` y las instrucciones de `AGENTS.md`. Antes de avanzar de fase, deben validarse los criterios de aceptacion definidos en el plan.
