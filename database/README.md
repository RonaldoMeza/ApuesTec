# Base de datos ApuesTec

PostgreSQL es la fuente oficial de verdad de ApuesTec. Redis se usa solo como cache y optimizacion para rankings, rate limiting e invitaciones temporales; cualquier dato critico debe poder reconstruirse desde PostgreSQL.

## Estructura

```text
database/
  migrations/
    001_init_schema.up.sql
    001_init_schema.down.sql
  seeds/
    001_seed_roles.sql
    002_seed_sample_teams_matches.sql
  scripts/
    apply-migrations.ps1
    apply-seeds.ps1
    check-tables.ps1
    reset-database.ps1
  .env.example
  README.md
```

## Flujo recomendado para inicializar la base local

Ejecuta estos comandos desde la raiz del proyecto. Los scripts requieren que los contenedores esten levantados y no se ejecutan automaticamente con `docker compose up`.

1. Levantar servicios:

```powershell
docker compose up --build -d
```

2. Aplicar migraciones:

```powershell
.\database\scripts\apply-migrations.ps1
```

3. Aplicar seeds:

```powershell
.\database\scripts\apply-seeds.ps1
```

4. Validar tablas y roles:

```powershell
.\database\scripts\check-tables.ps1
```

5. Abrir pgAdmin y refrescar el arbol de objetos:

- Host: `localhost`
- Port: `5433`
- Database: valor de `POSTGRES_DB` en `database/.env`
- User: valor de `POSTGRES_USER` en `database/.env`
- Password: valor de `POSTGRES_PASSWORD` en `database/.env`

## Tablas iniciales

- `users`
- `auth_accounts`
- `refresh_tokens`
- `roles`
- `user_roles`
- `rooms`
- `room_members`
- `room_invites`
- `teams`
- `matches`
- `predictions`
- `score_events`
- `audit_logs`

## Scripts manuales

- `database/scripts/apply-migrations.ps1`: aplica `001_init_schema.up.sql` sobre la base configurada en el contenedor `postgres`; si ya existen tablas publicas, no reejecuta la migracion para evitar modificar datos locales.
- `database/scripts/apply-seeds.ps1`: aplica los seeds minimos idempotentes de roles, equipos y partidos.
- `database/scripts/check-tables.ps1`: lista tablas y roles iniciales para confirmar el estado local.
- `database/scripts/reset-database.ps1`: elimina y recrea tablas locales usando `down`, `up` y seeds; pide confirmacion porque borra datos locales.

## Archivos SQL

- `database/migrations/001_init_schema.up.sql`: crea las tablas iniciales.
- `database/migrations/001_init_schema.down.sql`: elimina las tablas iniciales.
- `database/seeds/001_seed_roles.sql`: inserta roles globales minimos.
- `database/seeds/002_seed_sample_teams_matches.sql`: inserta equipos y partidos de prueba local.

## Ejecutar migraciones manualmente sin script

Con los servicios levantados:

```bash
docker compose up -d postgres
docker compose exec -T postgres psql -U apuestec_user -d apuestec < database/migrations/001_init_schema.up.sql
```

En Windows PowerShell:

```powershell
Get-Content -Raw .\database\migrations\001_init_schema.up.sql | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'
```

## Ejecutar seeds manualmente sin script

```bash
docker compose exec -T postgres psql -U apuestec_user -d apuestec < database/seeds/001_seed_roles.sql
docker compose exec -T postgres psql -U apuestec_user -d apuestec < database/seeds/002_seed_sample_teams_matches.sql
```

En Windows PowerShell:

```powershell
Get-Content -Raw .\database\seeds\001_seed_roles.sql | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'
Get-Content -Raw .\database\seeds\002_seed_sample_teams_matches.sql | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'
```

## Revertir migraciones manualmente sin script

```bash
docker compose exec -T postgres psql -U apuestec_user -d apuestec < database/migrations/001_init_schema.down.sql
```

En Windows PowerShell:

```powershell
Get-Content -Raw .\database\migrations\001_init_schema.down.sql | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'
```

## Validaciones basicas

Listar tablas:

```bash
docker compose exec postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "\dt"'
```

Validar constraints principales:

```bash
docker compose exec postgres psql -U apuestec_user -d apuestec -c "SELECT conname, contype FROM pg_constraint WHERE conrelid IN ('users'::regclass, 'rooms'::regclass, 'room_members'::regclass, 'room_invites'::regclass, 'matches'::regclass, 'predictions'::regclass) ORDER BY conname;"
```

Validar seeds:

```bash
docker compose exec postgres psql -U apuestec_user -d apuestec -c "SELECT name FROM roles ORDER BY name;"
docker compose exec postgres psql -U apuestec_user -d apuestec -c "SELECT COUNT(*) AS teams_count FROM teams; SELECT COUNT(*) AS matches_count FROM matches;"
```

## Notas

- Los archivos `.env` reales no deben commitearse.
- Los scripts usan las variables del contenedor `postgres`; no contienen credenciales reales.
- Las migraciones y seeds son manuales para evitar modificar o reinicializar accidentalmente una base local con datos.
- `reset-database.ps1` elimina y recrea la base local; usar solo en desarrollo local.
- Los codigos de invitacion se almacenan como hash en `room_invites.code_hash`.
- `score_events` no incluye `room_id`, `scope` ni puntajes independientes por sala.
- `room_members` permite eliminar la pertenencia activa cuando un usuario abandona una sala; el evento `ROOM_LEFT` debe registrarse en `audit_logs` en fases de backend.
