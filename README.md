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

## Docker Compose

El comando esperado para levantar el proyecto en fases futuras sera:

```bash
docker compose up --build
```

Para validar escalamiento horizontal del backend en fases posteriores:

```bash
docker compose up --scale backend=3
```

## Nota de desarrollo

Todo desarrollo debe seguir `PLAN_PROYECTO_APUESTEC.md` y las instrucciones de `AGENTS.md`. Antes de avanzar de fase, deben validarse los criterios de aceptacion definidos en el plan.
