# Fase 3 - Base de datos

## Objetivo de la fase

Definir el modelo inicial de datos de ApuesTec en PostgreSQL mediante migraciones SQL, seeds minimos y documentacion tecnica, sin implementar autenticacion completa, salas, predicciones, puntajes, rankings ni pantallas funcionales.

Nota: `PLAN_PROYECTO_APUESTEC.md` enumera Base de Datos como Fase 4, mientras que `prompt3.md` la solicita como Fase 3. Se implemento el alcance de Base de Datos porque coincide con el modelo de datos propuesto por el plan y fue solicitado explicitamente despues de la Fase 2.

## Tablas creadas

- `users`: usuarios registrados y estado operativo.
- `auth_accounts`: identidades externas o proveedores de autenticacion asociados a usuarios.
- `refresh_tokens`: refresh tokens revocables almacenados como hash.
- `roles`: roles globales `SUPER_ADMIN`, `ADMIN` y `USER`.
- `user_roles`: relacion entre usuarios y roles globales.
- `rooms`: salas privadas con estados `ACTIVE` y `CLOSED`.
- `room_members`: pertenencia activa a salas con roles `OWNER`, `MODERATOR` y `MEMBER`.
- `room_invites`: invitaciones temporales con codigo hasheado y estado operativo.
- `teams`: equipos gestionados manualmente para el MVP.
- `matches`: partidos gestionados manualmente para el MVP.
- `predictions`: predicciones de marcador por usuario y partido.
- `score_events`: eventos reales de puntaje global oficial.
- `audit_logs`: auditoria de eventos importantes, incluido `ROOM_LEFT`.

## Relaciones principales

- `auth_accounts.user_id` referencia `users.id`.
- `refresh_tokens.user_id` referencia `users.id`.
- `user_roles.user_id` referencia `users.id` y `user_roles.role_id` referencia `roles.id`.
- `rooms.created_by` referencia `users.id`.
- `room_members.room_id` referencia `rooms.id` y `room_members.user_id` referencia `users.id`.
- `room_invites.room_id` referencia `rooms.id` y `room_invites.created_by` referencia `users.id`.
- `matches.home_team_id` y `matches.away_team_id` referencian `teams.id`.
- `predictions.user_id` referencia `users.id` y `predictions.match_id` referencia `matches.id`.
- `score_events.user_id` referencia `users.id`, `score_events.match_id` referencia `matches.id` y `score_events.prediction_id` referencia `predictions.id`.
- `audit_logs.user_id` referencia `users.id` y permite `NULL` si el usuario se elimina.

## Restricciones aplicadas

- `UNIQUE(users.email)` mediante `users_email_unique`.
- `UNIQUE(room_members.room_id, room_members.user_id)` mediante `room_members_room_user_unique`.
- `UNIQUE(predictions.user_id, predictions.match_id)` mediante `predictions_user_match_unique`.
- `roles.name` es unico y restringido a roles globales del MVP.
- `room_invites.code_hash` es unico.
- Indice unico parcial `room_invites_one_active_per_room_idx` para permitir solo una invitacion activa por sala.
- Checks de puntajes y marcadores no negativos donde corresponde.
- Check para impedir partidos con el mismo equipo local y visitante.
- Check para exigir marcadores cuando un partido esta `FINISHED`.
- Check para duraciones de invitacion permitidas: 1, 3, 5, 10, 15 y 20 minutos.

## Estados definidos

- `users.status`: `ACTIVE`, `BLOCKED`, `DISABLED`.
- `rooms.status`: `ACTIVE`, `CLOSED`.
- `room_invites.status`: `ACTIVE`, `EXPIRED`, `REVOKED`.
- `matches.status`: `SCHEDULED`, `LOCKED`, `FINISHED`, `CANCELLED`.

## Seeds creados

- `database/seeds/001_seed_roles.sql`: inserta `SUPER_ADMIN`, `ADMIN` y `USER` sin crear usuarios reales.
- `database/seeds/002_seed_sample_teams_matches.sql`: inserta equipos y partidos de prueba para desarrollo local.

## Variables de entorno relacionadas

PostgreSQL usa `database/.env`, creado localmente desde `database/.env.example`:

```env
POSTGRES_DB=apuestec
POSTGRES_USER=apuestec_user
POSTGRES_PASSWORD=apuestec_password
POSTGRES_PORT=5432
```

El backend usara `DATABASE_URL` en fases posteriores. No se deben commitear archivos `.env` reales ni secretos.

## Conexion local desde pgAdmin

PostgreSQL expone el puerto `5433` del host solo para administracion local en desarrollo. Esta exposicion permite inspeccionar la base desde pgAdmin instalado en la maquina host sin cambiar la comunicacion interna de Docker.

Parametros para registrar el servidor en pgAdmin:

- Host: `localhost`
- Port: `5433`
- Database: valor de `POSTGRES_DB` en `database/.env`
- User: valor de `POSTGRES_USER` en `database/.env`
- Password: valor de `POSTGRES_PASSWORD` en `database/.env`

El backend debe seguir usando `DATABASE_URL` con host interno `postgres` y puerto `5432`. No se expone Redis al host, no se expone el backend directamente al host y Nginx sigue siendo el unico punto de entrada publico por `http://localhost:8081`.

## Fase 3.1 - Scripts manuales de base de datos

Se agregaron scripts PowerShell en `database/scripts/` para que el usuario ejecute migraciones, seeds, validaciones y resets locales bajo demanda.

- `apply-migrations.ps1`: aplica `database/migrations/001_init_schema.up.sql`; si ya existen tablas publicas, no reejecuta la migracion para evitar modificar datos locales.
- `apply-seeds.ps1`: aplica `database/seeds/001_seed_roles.sql` y `database/seeds/002_seed_sample_teams_matches.sql`.
- `check-tables.ps1`: lista tablas y roles iniciales.
- `reset-database.ps1`: ejecuta `down`, `up` y seeds despues de pedir confirmacion; elimina datos locales.

Las migraciones y seeds no se ejecutan automaticamente con `docker compose up`. Esta decision evita modificar o reinicializar accidentalmente una base local que ya tenga datos. Los scripts usan `POSTGRES_USER` y `POSTGRES_DB` dentro del contenedor `postgres`, por lo que no contienen credenciales reales.

Flujo recomendado desde la raiz del proyecto:

```powershell
docker compose up --build -d
.\database\scripts\apply-migrations.ps1
.\database\scripts\apply-seeds.ps1
.\database\scripts\check-tables.ps1
```

Para pgAdmin, refrescar el arbol de objetos despues de ejecutar migraciones y seeds. La conexion local usa `localhost:5433`, mientras la aplicacion mantiene la conexion interna `postgres:5432`.

## Comandos para ejecutar migraciones o scripts SQL

Levantar PostgreSQL:

```bash
docker compose up -d postgres
```

Ejecutar migracion up:

```powershell
.\database\scripts\apply-migrations.ps1
```

Ejecutar seeds:

```powershell
.\database\scripts\apply-seeds.ps1
```

Revertir migracion:

```powershell
Get-Content -Raw database/migrations/001_init_schema.down.sql | docker compose exec -T postgres psql -U apuestec_user -d apuestec
```

## Comandos para validar tablas en PostgreSQL

Listar tablas:

```powershell
.\database\scripts\check-tables.ps1
```

Validar constraints:

```bash
docker compose exec postgres psql -U apuestec_user -d apuestec -c "SELECT conrelid::regclass AS table_name, conname, contype FROM pg_constraint WHERE connamespace = 'public'::regnamespace ORDER BY table_name, conname;"
```

Validar que `score_events` no tenga scope de sala:

```bash
docker compose exec postgres psql -U apuestec_user -d apuestec -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'score_events' ORDER BY ordinal_position;"
```

Validar seeds:

```bash
docker compose exec postgres psql -U apuestec_user -d apuestec -c "SELECT name FROM roles ORDER BY name;"
docker compose exec postgres psql -U apuestec_user -d apuestec -c "SELECT COUNT(*) AS teams_count FROM teams; SELECT COUNT(*) AS matches_count FROM matches;"
```

## Decisiones tecnicas tomadas

- Se uso `UUID` con `gen_random_uuid()` mediante `pgcrypto` para claves primarias.
- Se usaron `CHECK constraints` para estados principales, roles por sala y duraciones de invitacion.
- Se agrego un indice unico parcial para reforzar desde PostgreSQL que solo exista una invitacion activa por sala.
- `score_events` no contiene `room_id`, `scope` ni campos de puntaje por sala para evitar duplicar puntos por pertenencia a salas.
- `room_members` representa solo pertenencias activas; en el MVP se puede eliminar el registro cuando un usuario abandona una sala.
- `audit_logs.action` no tiene check restrictivo para permitir registrar eventos futuros como `ROOM_LEFT` sin nueva migracion inmediata.
- Los seeds no crean usuarios ni contrasenas reales.
- Los seeds son idempotentes: roles usan `ON CONFLICT`, equipos usan `ON CONFLICT (country_code)` y partidos evitan duplicados con `WHERE NOT EXISTS`.
- Las migraciones y seeds se ejecutan manualmente mediante scripts para evitar cambios destructivos accidentales al levantar Docker.

## Problemas encontrados y soluciones

- La numeracion de fase difiere entre `prompt3.md` y el plan. Se documento la diferencia y se implemento solo el alcance de Base de Datos solicitado.
- No existia `database/README.md`, `database/migrations/` ni `database/seeds/`. Se crearon con la estructura requerida.
- No existia `database/scripts/`. Se creo para Fase 3.1 con scripts manuales de administracion local.
- No se agrego logica de backend ni pantallas para respetar el alcance de la fase.

## Resultado de validacion

La validacion se ejecuto dentro del contenedor PostgreSQL usando una base temporal `apuestec_phase3_validation`, para no modificar datos locales existentes.

- `docker compose up -d postgres`: correcto, PostgreSQL quedo healthy.
- Creacion de base temporal `apuestec_phase3_validation`: correcta.
- `database/migrations/001_init_schema.up.sql`: correcto.
- `database/seeds/001_seed_roles.sql`: correcto, inserto 3 roles.
- `database/seeds/002_seed_sample_teams_matches.sql`: correcto, inserto 6 equipos y 3 partidos.
- Validacion de tablas: correcto, existen 13 tablas publicas.
- Validacion de constraints obligatorias: correcto, existen `users_email_unique`, `room_members_room_user_unique`, `predictions_user_match_unique` y checks de estados principales.
- Validacion de `score_events`: correcto, contiene `id`, `user_id`, `match_id`, `prediction_id`, `points`, `reason` y `created_at`; no contiene `room_id`, `scope` ni campos de puntaje por sala.
- `database/migrations/001_init_schema.down.sql`: correcto, dejo 0 tablas publicas en la base temporal.
- Eliminacion de base temporal: correcta.

Validacion Fase 3.1 sobre la base local principal `apuestec`:

- `docker compose config --quiet`: correcto, sin salida.
- `docker compose up -d`: correcto, PostgreSQL, Redis, backend, frontend y Nginx quedaron healthy.
- `.\database\scripts\apply-migrations.ps1`: correcto; detecto que la base local ya tenia tablas publicas y no reejecuto la migracion para evitar modificar datos existentes.
- `.\database\scripts\apply-seeds.ps1`: correcto; los seeds se aplicaron de forma idempotente.
- `.\database\scripts\check-tables.ps1`: correcto; listo 13 tablas publicas y los roles `ADMIN`, `SUPER_ADMIN` y `USER`.
- Validacion de idempotencia: correcto; despues de ejecutar seeds nuevamente, quedaron 3 roles, 6 equipos y 3 partidos, sin duplicados.
- `docker compose ps`: correcto; backend healthy, Nginx healthy y PostgreSQL healthy con `0.0.0.0:5433->5432/tcp`.
- `http://localhost:8081`: correcto, respondio `200`.
- `http://localhost:8081/api/v1/health`: correcto, respondio `200` con estado `ok` del backend.
- Puerto local para pgAdmin: correcto; `localhost:5433` acepta conexiones TCP y las 13 tablas quedan disponibles para pgAdmin despues de refrescar el arbol de objetos.
- No se ejecuto `reset-database.ps1` para no eliminar datos locales existentes durante la validacion.

## Checklist final

- [x] La migracion up crea todas las tablas minimas.
- [x] La migracion down revierte las tablas en orden correcto.
- [x] Existen claves primarias.
- [x] Existen claves foraneas principales.
- [x] Existen restricciones UNIQUE obligatorias.
- [x] Existen CHECK constraints para estados principales.
- [x] `score_events` no maneja puntajes por sala.
- [x] `rooms` solo usa `ACTIVE` y `CLOSED`.
- [x] `predictions` mantiene `UNIQUE(user_id, match_id)`.
- [x] `room_members` mantiene `UNIQUE(room_id, user_id)`.
- [x] `audit_logs` permite registrar `ROOM_LEFT`.
- [x] Seeds minimos existen.
- [x] Existen scripts manuales para aplicar migraciones y seeds.
- [x] Existen scripts manuales para validar tablas y resetear la base local.
- [x] Docker Compose no ejecuta migraciones ni seeds automaticamente.
- [x] Los scripts usan variables del contenedor PostgreSQL y no credenciales reales.
- [x] pgAdmin puede inspeccionar la base por `localhost:5433` despues de refrescar.
- [x] Backend sigue healthy despues de aplicar scripts manuales.
- [x] Nginx responde en `http://localhost:8081`.
- [x] `http://localhost:8081/api/v1/health` responde correctamente.
- [x] No se implementa logica de negocio todavia.
- [x] No se cambia el stack tecnologico.
- [x] No se usaron credenciales reales.
