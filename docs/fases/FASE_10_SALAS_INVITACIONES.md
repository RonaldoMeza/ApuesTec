# Fase 10: Salas Privadas e Invitaciones Temporales

## Objetivo

Implementar salas privadas con membresia, roles por sala (OWNER, MODERATOR, MEMBER), invitaciones temporales mediante codigo, ranking de sala basado en puntaje global acumulado, y toda la interfaz de usuario asociada.

## Archivos Creados

### Backend

| Archivo | Proposito |
|---------|-----------|
| `backend/migrations/005_rooms.sql` | Migracion con tablas `rooms`, `room_members` y `room_invites` |
| `backend/internal/rooms/model.go` | Modelos Room, RoomMember, ServiceError y constantes de rol/estado |
| `backend/internal/rooms/dto.go` | DTOs para creacion, actualizacion, respuesta y leaderboard de sala |
| `backend/internal/rooms/repository.go` | Repositorio con CRUD de salas, membresia y ranking |
| `backend/internal/rooms/service.go` | Logica de negocio con validacion de permisos por rol |
| `backend/internal/rooms/handler.go` | Handlers REST para CRUD de salas y leaderboard |
| `backend/internal/roommembers/model.go` | Constantes de rol y ServiceError |
| `backend/internal/roommembers/dto.go` | DTOs para miembros y cambio de rol |
| `backend/internal/roommembers/repository.go` | Repositorio de membresia (listar, cambiar rol, expulsar) |
| `backend/internal/roommembers/service.go` | Logica de membresia con reglas de permiso por rol |
| `backend/internal/roommembers/handler.go` | Handlers REST para miembros |
| `backend/internal/roominvites/model.go` | Modelo RoomInvite, duraciones permitidas y ServiceError |
| `backend/internal/roominvites/dto.go` | DTOs para creacion, preview y respuesta de invitacion |
| `backend/internal/roominvites/repository.go` | Repositorio de invitaciones con generacion de codigo criptografico |
| `backend/internal/roominvites/service.go` | Logica de invitaciones con validacion, revocacion y auto-join |
| `backend/internal/roominvites/handler.go` | Handlers REST para invitaciones y union |
| `backend/internal/scoring/handler.go` | Handler POST /admin/rebuild-scores para reconstruir puntajes |

### Archivos Modificados

| Archivo | Cambio |
|---------|--------|
| `backend/cmd/api/main.go` | Agregada llamada a `RebuildUserScores` al iniciar el backend |
| `backend/internal/audit/repository.go` | Agregadas 10 constantes de auditoria para salas (ROOM_CREATED, ROOM_MEMBER_JOINED, etc.) |
| `backend/internal/database/migrate.go` | Agregada funcion `RebuildUserScores` con consulta usando subconsultas para evitar producto cartesiano |
| `backend/internal/scoring/repository.go` | Fix critico: `RebuildAllUserScores` ahora usa subconsultas en lugar de JOIN directo para evitar inflar SUM/COUNT; `UpsertUserScore` ahora actualiza `streak_count` en ON CONFLICT |
| `backend/internal/leaderboard/repository.go` | Fix: `COUNT(*)` ahora cuenta desde `users` en lugar de `user_scores` para consistencia |
| `backend/internal/routes/router.go` | Agregadas todas las rutas de salas, miembros, invitaciones y leaderboard de sala |

### Frontend

| Archivo | Proposito |
|---------|-----------|
| `frontend/features/rooms/types/room.types.ts` | Tipos Room, RoomMember, RoomLeaderboardEntry, RoomInvite |
| `frontend/features/rooms/services/room.service.ts` | Servicio de llamadas API para salas |
| `frontend/features/rooms/components/RoomCard.tsx` | Card de sala con estado, rol, miembros |
| `frontend/features/rooms/components/RoomForm.tsx` | Formulario crear/editar sala |
| `frontend/features/rooms/components/RoomMemberList.tsx` | Lista de miembros con controles segun rol |
| `frontend/features/rooms/components/RoomInviteBox.tsx` | Generador de invitacion con selector de duracion y codigo |
| `frontend/features/rooms/components/RoomLeaderboard.tsx` | Tabla de ranking de sala con puntaje global |
| `frontend/features/rooms/components/RoomRoleBadge.tsx` | Badge de rol coloreado |
| `frontend/features/invites/types/invite.types.ts` | Tipos InvitePreview y JoinResponse |
| `frontend/features/invites/services/invite.service.ts` | Servicio de invitaciones |
| `frontend/features/invites/components/InvitePreviewCard.tsx` | Card de preview de invitacion |
| `frontend/features/invites/components/QRCodeBox.tsx` | Visualizacion de QR payload con copia |
| `frontend/app/rooms/page.tsx` | Pagina protegida /rooms - lista de salas del usuario |
| `frontend/app/rooms/create/page.tsx` | Pagina protegida /rooms/create - formulario de creacion |
| `frontend/app/rooms/[id]/page.tsx` | Pagina protegida /rooms/[id] - detalle con tabs (Ranking, Miembros, Invitacion, Configuracion) |
| `frontend/app/join/[code]/page.tsx` | Pagina /join/[code] - preview y union a sala por codigo |

### Paginas Modificadas

| Pagina | Cambio |
|--------|--------|
| `frontend/app/dashboard/page.tsx` | Agregada seccion "Mis Salas" con hasta 3 salas y CTA |
| `frontend/app/page.tsx` | Agregada seccion "Salas Privadas" explicativa con CTA |
| `frontend/shared/components/AppLayout.tsx` | Agregado enlace "Salas" en navegacion para usuarios autenticados |

## Tablas Agregadas

### rooms

| Columna | Tipo | Descripcion |
|---------|------|-------------|
| id | UUID PK | Identificador unico |
| name | VARCHAR(120) NOT NULL | Nombre de la sala |
| description | TEXT NULL | Descripcion opcional |
| owner_id | UUID FK | Usuario propietario |
| status | VARCHAR(20) DEFAULT 'ACTIVE' | ACTIVE o CLOSED |
| created_at | TIMESTAMPTZ | Fecha de creacion |
| updated_at | TIMESTAMPTZ | Ultima actualizacion |
| closed_at | TIMESTAMPTZ NULL | Fecha de cierre |

### room_members

| Columna | Tipo | Descripcion |
|---------|------|-------------|
| id | UUID PK | Identificador unico |
| room_id | UUID FK | Sala |
| user_id | UUID FK | Usuario |
| role | VARCHAR(20) | OWNER, MODERATOR o MEMBER |
| joined_at | TIMESTAMPTZ | Fecha de ingreso |
| created_at | TIMESTAMPTZ | Fecha de creacion |
| updated_at | TIMESTAMPTZ | Ultima actualizacion |
| UNIQUE(room_id, user_id) | | No duplicados |

### room_invites

| Columna | Tipo | Descripcion |
|---------|------|-------------|
| id | UUID PK | Identificador unico |
| room_id | UUID FK | Sala asociada |
| code | VARCHAR(32) UNIQUE | Codigo criptografico de 32 caracteres hex |
| qr_payload | TEXT NULL | URL para QR |
| created_by | UUID FK | Creador de la invitacion |
| expires_at | TIMESTAMPTZ | Fecha de expiracion |
| used_at | TIMESTAMPTZ NULL | Fecha de uso |
| revoked_at | TIMESTAMPTZ NULL | Fecha de revocacion |
| created_at | TIMESTAMPTZ | Fecha de creacion |

## Endpoints Implementados

### Salas (requiere autenticacion)

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | /api/v1/rooms | Lista salas del usuario autenticado |
| POST | /api/v1/rooms | Crea una sala (creador queda OWNER) |
| GET | /api/v1/rooms/:id | Detalle de sala (requiere ser miembro) |
| PUT | /api/v1/rooms/:id | Editar nombre/descripcion (solo OWNER) |
| PATCH | /api/v1/rooms/:id/close | Cerrar sala (solo OWNER) |
| GET | /api/v1/rooms/:id/members | Lista miembros |
| PATCH | /api/v1/rooms/:id/members/:userId/role | Cambiar rol (solo OWNER) |
| DELETE | /api/v1/rooms/:id/members/:userId | Expulsar miembro (OWNER/MODERATOR segun reglas) |
| POST | /api/v1/rooms/:id/leave | Abandonar sala (no OWNER) |
| GET | /api/v1/rooms/:id/leaderboard | Ranking de sala basado en puntaje global |
| POST | /api/v1/rooms/:id/invites | Generar invitacion temporal |

### Invitaciones

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | /api/v1/invites/:code | Vista previa publica de invitacion |
| POST | /api/v1/invites/:code/join | Unirse a sala con codigo (autenticado) |

### Admin

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| POST | /api/v1/admin/rebuild-scores | Reconstruir todos los puntajes de usuario |

## Roles de Sala

| Rol | Permisos |
|-----|----------|
| OWNER | Editar sala, cerrar sala, generar invitaciones, cambiar roles (excepto OWNER), expulsar MODERATOR y MEMBER |
| MODERATOR | Generar invitaciones, ver miembros, expulsar solo MEMBER |
| MEMBER | Ver sala, ver ranking, ver miembros, abandonar sala |

## Reglas de Permisos

1. Solo OWNER puede editar nombre/descripcion de la sala
2. Solo OWNER puede cerrar la sala
3. Solo OWNER puede cambiar roles (MEMBER <-> MODERATOR), nunca al OWNER
4. OWNER puede expulsar MODERATOR o MEMBER, pero no a si mismo
5. MODERATOR puede expulsar solo MEMBER
6. MEMBER no puede expulsar ni administrar
7. MEMBER y MODERATOR pueden abandonar la sala
8. OWNER no puede abandonar (debe cerrar la sala)
9. Salas CLOSED no permiten invitaciones ni nuevos miembros

## Flujo de Invitaciones

1. OWNER o MODERATOR genera invitacion con duracion (1, 3, 5, 10, 15 o 20 minutos)
2. Al generar, se revoca automaticamente cualquier invitacion activa anterior de la misma sala
3. Se genera codigo criptografico de 32 caracteres hexadecimales (16 bytes via crypto/rand)
4. La invitacion tiene expiracion basada en duracion elegida
5. Cualquier usuario autenticado puede usar el codigo para unirse
6. Al unirse, la invitacion se marca como usada y el usuario obtiene rol MEMBER
7. No se puede usar invitacion expirada, revocada o ya usada
8. Si el usuario ya pertenece a la sala, se rechaza con mensaje claro
9. Si la sala esta CLOSED, no se permite unirse

## Ranking de Sala

- Se consulta via `GET /api/v1/rooms/:id/leaderboard`
- Devuelve miembros de la sala ordenados por su puntaje global acumulado (`user_scores.total_points`)
- No se crean puntajes independientes por sala
- Incluye: rank, userId, fullName, totalPoints, predictionsCount, exactScoresCount, roomRole
- Se muestra mensaje: "El ranking de sala usa el puntaje global acumulado"

## Fix del Ranking Global (bug critico corregido)

El `RebuildAllUserScores` original usaba `LEFT JOIN` directo entre `predictions` y `score_events`, generando un producto cartesiano que inflaba los valores de `SUM` y `COUNT`. Se reemplazo por subconsultas independientes que agregan cada tabla por separado antes del JOIN.

Adicionalmente:
- `UpsertUserScore` ahora actualiza `streak_count` en ON CONFLICT
- `leaderboard COUNT` ahora cuenta desde `users` (no `user_scores`) para consistencia
- Se agrego `POST /api/v1/admin/rebuild-scores` para reconstruccion manual
- Se llama `RebuildUserScores` automaticamente al iniciar el backend

## Pruebas Realizadas

- `go build ./...` - Backend compila sin errores
- `go test ./...` - Tests existentes pasan
- `npm run lint` - Frontend sin errores de linter
- `npm run build` - Frontend compila sin errores
- `docker compose config --quiet` - Configuracion valida

## Como Probar desde Frontend

1. Iniciar sesion como USER1
2. Ir a /rooms - Ver empty state
3. Crear sala - Ver redireccion a /rooms/[id]
4. Ver que aparece como OWNER
5. Ver ranking de sala (vacio o con USER1)
6. Ir a pestaña "Invitacion" y generar invitacion de 5 minutos
7. Copiar codigo o enlace
8. Cerrar sesion e iniciar como USER2
9. Abrir /join/[codigo]
10. Unirse a la sala
11. Ver que aparece como MEMBER
12. Volver con USER1 y promover USER2 a MODERATOR
13. Verificar que MODERATOR puede generar invitacion
14. Verificar que MODERATOR puede expulsar solo MEMBERS
15. Verificar que MEMBER no puede administrar
16. Cerrar sala como OWNER
17. Verificar que sala CLOSED no permite nuevas invitaciones
18. Verificar que invitacion expirada no permite unirse
19. Verificar que ranking de sala muestra puntaje global
20. Verificar que /leaderboard global sigue funcionando correctamente con puntos reales

## Mejoras: Salas Publicas, Busqueda por Red y Contrasena

### Cambios en Base de Datos (006_rooms_visibility.sql)

| Columna | Tabla | Tipo | Descripcion |
|---------|-------|------|-------------|
| `visibility` | rooms | `VARCHAR(7) DEFAULT 'PRIVATE'` | `PUBLIC` o `PRIVATE` |
| `password_hash` | rooms | `VARCHAR(255)` | Hash bcrypt de la contrasena (solo publicas) |
| `network_prefix` | rooms | `VARCHAR(15)` | Prefijo de red `/24` del creador |

Indice: `idx_rooms_public_network` sobre `(visibility, network_prefix)` filtrando `WHERE visibility = 'PUBLIC'`

### Nuevos Endpoints

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| `GET` | `/api/v1/rooms/public?q=busqueda` | Busca salas publicas en la misma red |
| `POST` | `/api/v1/rooms/{id}/join` | Unirse a sala publica (requiere password si tiene) |

### Funcionalidades Nuevas

- **Salas publicas/privadas**: Al crear sala, switch para elegir visibilidad. Por defecto privada.
- **Contrasena de sala**: Las salas publicas requieren contrasena (bcrypt). Quien la conozca puede unirse.
- **Busqueda por red**: Al buscar salas publicas, solo aparecen las que comparten el mismo prefijo de red `/24` (e.g., `192.168.1`).
- **Unirse por codigo**: Input siempre visible en `/rooms`, fuera de las pestañas, accesible sin necesidad de tener salas creadas.
- **Busqueda dinamica**: Input con debounce (400ms) dentro de la pestana "Descubrir" que actualiza resultados automaticamente al escribir.
- **Pestana "Descubrir"**: En `/rooms`, pestana para explorar salas publicas en la misma red con buscador.
- **Modal de contrasena**: Al unirse a una sala con contrasena, se abre un modal para ingresarla.
- **Toasts (sonner)**: Notificaciones para crear, unirse, abandonar, actualizar, cerrar sala, cambiar rol, expulsar y generar invitacion.
- **Cursor pointer**: Todos los botones e inputs interactivos tienen `cursor-pointer` y `disabled:cursor-not-allowed`.
- **Sala cerrada**: Al cerrar una sala, todos los campos de configuracion se bloquean (inputs, switch, botones), la pestaña "Invitacion" se oculta y los botones de administrar miembros se deshabilitan.

### Como Probar las Nuevas Funcionalidades

1. USER1 crea sala publica con contrasena (e.g., "mundial123")
2. USER1 ve que la sala aparece con badge "Publica" en el detalle
3. USER2 abre `/rooms` y ve el input "Unirse por codigo" siempre visible
4. USER2 va a pestana "Descubrir" y escribe en el buscador
5. Los resultados se actualizan automaticamente mientras escribe
6. USER2 hace clic en "Unirse" y si tiene contrasena, ingresa password en modal
7. USER2 es redirigido a la sala como MEMBER con toast de exito
8. USER1 cierra la sala y todos los campos de configuracion se bloquean

## Pendiente para Fase 11

- Cache de ranking global y de sala con Redis
- Rate limiting con Redis
- Pruebas de estres con k6
- QR code visual real (actualmente solo se muestra payload)
- Notificaciones en tiempo real (WebSocket o Server-Sent Events)
- Moderacion de contenido en salas
- Historial de actividad de sala
