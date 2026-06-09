# Fase 7: Base de Equipos y Partidos

## Objetivo

Implementar la base de equipos y partidos para ApuesTec, permitiendo:
- Registrar, listar, editar y eliminar equipos
- Registrar partidos manualmente
- Listar partidos próximos y finalizados
- Ver detalle básico de un partido
- Mostrar partidos desde el frontend
- Agregar sección visual de próximos partidos en Home pública
- Agregar vista protegida administrativa para cargar partidos manualmente

## Archivos creados/modificados

### Backend (`backend/internal/`)

| Archivo | Descripción |
|---------|-------------|
| `teams/model.go` | Estructura `Team` con campos: id, name, code, country, flagUrl, createdAt, updatedAt |
| `teams/dto.go` | `CreateTeamRequest`, `UpdateTeamRequest`, `TeamResponse`, errores de dominio |
| `teams/repository.go` | `Repository` interface + implementación con Create, FindAll, FindByID, Update, Delete |
| `teams/service.go` | `Service` interface + implementación con Create, List, GetByID, Update, Delete |
| `teams/handler.go` | Handlers HTTP para endpoints públicos y admin de equipos |
| `matches/model.go` | Estructura `Match` con campos y constantes de estado |
| `matches/dto.go` | `CreateMatchRequest`, `UpdateMatchRequest`, `UpdateStatusRequest`, `UpdateResultRequest`, `MatchResponse`, `TeamInfo`, errores de dominio |
| `matches/repository.go` | `Repository` interface + implementación con Create, FindAll, FindUpcoming, FindFinished, FindByID, Update, UpdateStatus, UpdateResult |
| `matches/service.go` | `Service` interface + implementación con Create, List, ListUpcoming, ListFinished, GetByID, Update, UpdateStatus, UpdateResult. Enriquecimiento automático con datos de equipos. |
| `matches/handler.go` | Handlers HTTP para endpoints públicos y admin de partidos |
| `matches/team_repo.go` | Adaptador `TeamInfoRepository` para consultar datos de equipos desde el módulo de matches |
| `routes/router.go` | Registro de todas las rutas nuevas con protección de roles |

### Frontend (`frontend/`)

| Archivo | Descripción |
|---------|-------------|
| `features/teams/types/team.types.ts` | Tipos TypeScript para equipos |
| `features/teams/services/team.service.ts` | Servicio HTTP para equipos |
| `features/matches/types/match.types.ts` | Tipos TypeScript para partidos |
| `features/matches/services/match.service.ts` | Servicio HTTP para partidos |
| `features/matches/components/MatchStatusBadge.tsx` | Badge visual para estado del partido |
| `features/matches/components/MatchCard.tsx` | Card de partido para listados |
| `features/matches/components/MatchList.tsx` | Lista grid de MatchCards con empty state |
| `features/admin/components/TeamForm.tsx` | Formulario crear/editar equipo |
| `features/admin/components/MatchForm.tsx` | Formulario crear partido |
| `features/admin/components/MatchResultForm.tsx` | Formulario registrar resultado |
| `app/matches/page.tsx` | Vista pública de partidos con tabs (próximos/finalizados) |
| `app/matches/[id]/page.tsx` | Detalle básico de partido |
| `app/admin/page.tsx` | Panel admin protegido por rol |
| `app/admin/teams/page.tsx` | Gestión admin de equipos |
| `app/admin/matches/page.tsx` | Gestión admin de partidos |

### Archivos modificados

| Archivo | Cambio |
|---------|--------|
| `frontend/app/page.tsx` | PreviewSection ahora obtiene partidos reales del backend |
| `frontend/app/dashboard/page.tsx` | Muestra próximos partidos, cards actualizadas, link a admin |
| `frontend/shared/components/AppLayout.tsx` | Nav incluye enlace a /matches |
| `frontend/shared/components/UserNav.tsx` | Dropdown incluye Partidos y Admin (si aplica) |
| `database/migrations/002_align_teams_matches_phase7.sql` | Migración para alinear esquema existente con Fase 7 |
| `database/seeds/002_seed_sample_teams_matches.sql` | Seed actualizado con 10 equipos y 5 partidos |

## Endpoints implementados

### Públicos (sin autenticación)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/v1/teams` | Listar todos los equipos |
| GET | `/api/v1/teams/:id` | Obtener equipo por ID |
| GET | `/api/v1/matches/upcoming` | Listar partidos próximos (SCHEDULED, LOCKED) |
| GET | `/api/v1/matches/finished` | Listar partidos finalizados (FINISHED, CANCELLED) |
| GET | `/api/v1/matches` | Listar todos los partidos |
| GET | `/api/v1/matches/:id` | Obtener detalle de partido |

### Admin (requiere rol ADMIN o SUPER_ADMIN)
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/v1/admin/teams` | Crear equipo |
| PUT | `/api/v1/admin/teams/:id` | Actualizar equipo |
| DELETE | `/api/v1/admin/teams/:id` | Eliminar equipo |
| POST | `/api/v1/admin/matches` | Crear partido |
| PUT | `/api/v1/admin/matches/:id` | Actualizar partido |
| PATCH | `/api/v1/admin/matches/:id/status` | Cambiar estado |
| PATCH | `/api/v1/admin/matches/:id/result` | Registrar resultado |

## Tablas

### teams
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID (PK) | Identificador único |
| name | VARCHAR(150) | Nombre del equipo |
| code | VARCHAR(10) UNIQUE | Código corto (ej: ARG) |
| country | VARCHAR(100) | País |
| flag_url | TEXT (nullable) | URL de bandera |
| created_at | TIMESTAMPTZ | Fecha de creación |
| updated_at | TIMESTAMPTZ | Fecha de actualización |

### matches
| Columna | Tipo | Descripción |
|---------|------|-------------|
| id | UUID (PK) | Identificador único |
| home_team_id | UUID (FK) | Equipo local |
| away_team_id | UUID (FK) | Equipo visitante |
| starts_at | TIMESTAMPTZ | Fecha/hora del partido |
| status | VARCHAR(20) | SCHEDULED, LOCKED, FINISHED, CANCELLED |
| home_score | INTEGER (nullable) | Goles local |
| away_score | INTEGER (nullable) | Goles visitante |
| locked_at | TIMESTAMPTZ (nullable) | Cuándo se bloqueó |
| created_at | TIMESTAMPTZ | Fecha de creación |
| updated_at | TIMESTAMPTZ | Fecha de actualización |

### Reglas de negocio
1. Un partido debe tener equipo local y visitante diferentes
2. Un partido nuevo inicia como SCHEDULED
3. SCHEDULED → LOCKED (bloqueo manual desde admin)
4. LOCKED → FINISHED (al registrar resultado)
5. SCHEDULED → CANCELLED
6. Los estados válidos son: SCHEDULED, LOCKED, FINISHED, CANCELLED
7. Endpoints admin requieren rol ADMIN o SUPER_ADMIN
8. Endpoints públicos son accesibles sin login

## Pruebas realizadas

### Backend
- `go test ./...` - Tests existentes de auth pasan
- `go build ./...` - Compilación exitosa

### Frontend
- `npm run lint` - ESLint sin errores
- `npm run build` - Build exitoso, todas las rutas generadas correctamente

### Docker
- `docker compose config --quiet` - Configuración válida

## Cómo probar desde frontend

1. Entrar a http://localhost:8081
2. Ver Home pública con sección "Próximos partidos"
3. Ir a `/matches` para ver lista de partidos (sin login)
4. Entrar a detalle de partido `/matches/:id` (sin login)
5. Iniciar sesión con usuario registrado
6. Ir a dashboard y ver "Próximos partidos" y cards actualizadas
7. Si el usuario tiene rol ADMIN o SUPER_ADMIN:
   - Ir a `/admin`
   - Crear equipo desde `/admin/teams`
   - Crear partido desde `/admin/matches`
   - Cambiar estado de partido (Bloquear, Cancelar)
   - Registrar resultado de partido bloqueado
8. Si el usuario es USER, `/admin` muestra "Acceso denegado"

### Promover usuario a ADMIN para pruebas
Ejecutar en la base de datos:
```sql
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email = 'usuario@ejemplo.com' AND r.name = 'ADMIN';
```

## Qué queda pendiente para Fase 8

- Predicciones de usuarios sobre partidos
- Cálculo de puntuaciones basado en aciertos
- Rankings reales (global y por sala)
- Bloqueo automático de partidos cuando inician
- Bloqueo automático de predicciones antes del inicio del partido
- Historial de predicciones del usuario
- Salas privadas con invitaciones
