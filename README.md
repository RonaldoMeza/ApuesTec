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

Docker Compose levanta el monorepo con los servicios base de infraestructura para desarrollo local. Nginx es el unico punto de entrada externo; el backend, PostgreSQL y Redis quedan disponibles dentro de la red interna de Docker.

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
- `redis:6379`: puerto interno de Redis.

## Servicios levantados

- `frontend`: aplicacion base Next.js para desarrollo local.
- `backend`: API Go/Gin minima con `GET /api/v1/health`.
- `postgres`: base de datos principal de ApuesTec con volumen persistente.
- `redis`: cache y capa de optimizacion con volumen persistente.
- `nginx`: reverse proxy y punto de entrada principal.
- `k6`: servicio preparado bajo perfil `tools` para pruebas posteriores.

## Persistencia y cache

PostgreSQL es la fuente oficial de verdad de ApuesTec. Redis se usa como cache y optimizacion para futuras funciones como rankings, rate limiting e invitaciones temporales; los datos criticos deben poder reconstruirse desde PostgreSQL.

## Variables de entorno

Docker Compose carga variables desde estos archivos locales:

- `database/.env`: credenciales y base inicial de PostgreSQL, basado en `database/.env.example`.
- `backend/.env`: configuracion de la API, URLs internas de PostgreSQL/Redis, JWT y OAuth, basado en `backend/.env.example`.
- `frontend/.env.local`: variables publicas de Next.js, basado en `frontend/.env.example`.

Variables clave por capa:

- Database: `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_PORT`.
- Backend: `APP_NAME`, `APP_ENV`, `APP_PORT`, `API_PREFIX`, `APP_PUBLIC_URL`, `FRONTEND_URL`, `CORS_ALLOWED_ORIGINS`, `DATABASE_URL`, `REDIS_URL`, `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`.
- Frontend: `NEXT_PUBLIC_APP_NAME`, `NEXT_PUBLIC_APP_URL`, `NEXT_PUBLIC_API_BASE_URL`, `NEXT_PUBLIC_GOOGLE_CLIENT_ID`.

No se deben commitear archivos `.env` reales ni secretos. El backend usa `postgres` y `redis` como hosts internos de Docker en `DATABASE_URL` y `REDIS_URL`.

## k6

La carpeta `tests/k6/` queda preparada para scripts de pruebas de estres en fases posteriores. El servicio puede invocarse con perfil de herramientas cuando existan scripts:

```bash
docker compose --profile tools run --rm k6 version
```

## Nota de desarrollo

Todo desarrollo debe seguir `PLAN_PROYECTO_APUESTEC.md` y las instrucciones de `AGENTS.md`. Antes de avanzar de fase, deben validarse los criterios de aceptacion definidos en el plan.
