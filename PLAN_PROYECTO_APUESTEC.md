# PLAN_PROYECTO_APUESTEC.md

## 1. Resumen Ejecutivo

ApuesTec es una plataforma educativa de predicciones deportivas para partidos del Mundial. El sistema permite a usuarios registrados predecir marcadores, acumular puntos, participar en rankings globales y por sala, crear salas privadas, gestionar invitaciones temporales y competir mediante reglas de gamificacion.

El sistema no maneja dinero real, pagos, cuotas, apuestas monetarias ni integracion con casas de apuestas. Su objetivo es academico y recreativo.

Este documento define el plan tecnico, arquitectonico y funcional del proyecto para evitar decisiones improvisadas, codigo espagueti y mezcla de responsabilidades.

## 2. Alcance del MVP

El MVP incluira:

- Registro e inicio de sesion con email y contrasena.
- Login alternativo con Google OAuth 2.0 / OpenID Connect.
- Emision de JWT propios desde el backend.
- Refresh tokens revocables almacenados como hash.
- Roles globales: SUPER_ADMIN, ADMIN y USER.
- Roles por sala: OWNER, MODERATOR y MEMBER.
- CRUD basico de usuarios segun permisos.
- Gestion manual de equipos y partidos.
- Registro de resultados oficiales por ADMIN o SUPER_ADMIN.
- Creacion de salas privadas.
- Invitaciones temporales con codigo y QR.
- Registro y edicion de predicciones hasta 10 minutos antes del partido.
- Calculo de puntos segun reglas definidas.
- Ranking global.
- Ranking por sala.
- Auditoria de eventos importantes.
- Redis para rankings, cache, rate limiting e invitaciones temporales.
- Docker Compose para levantar frontend, backend, PostgreSQL, Redis, Nginx y k6.
- Nginx como punto de entrada, reverse proxy y balanceador.
- Pruebas de estres con k6.
- Documentacion tecnica y operativa.

## 3. Alcance Fuera del MVP

No se implementara en el MVP:

- Dinero real.
- Pagos.
- Cuotas reales.
- Logica de apuestas monetarias.
- Integracion con APIs externas de partidos.
- Firebase Auth.
- IA o predicciones automaticas.
- Marketplace, recompensas economicas o premios reales.
- Panel avanzado de analitica predictiva.
- Aplicacion movil nativa.

## 4. Stack Tecnologico

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

Base de datos principal:

- PostgreSQL.

Base de datos en memoria:

- Redis.

Contenerizacion:

- Docker.
- Docker Compose.

Balanceo de carga:

- Nginx.

Pruebas de estres:

- k6.

Autenticacion:

- JWT propio generado por backend.
- Google OAuth 2.0 / OpenID Connect solo para validar identidad.

## 5. Tipo de Repositorio

El proyecto usara un monorepo con separacion clara por capas:

```text
ApuesTec/
  backend/
  frontend/
  database/
  nginx/
  tests/
    k6/
  docker-compose.yml
  PLAN_PROYECTO_APUESTEC.md
```

## 6. Arquitectura General en Capas

La arquitectura general sera por capas:

1. Frontend: presentacion y experiencia de usuario.
2. Nginx: reverse proxy, balanceador y punto unico de entrada.
3. Backend: API REST, logica de negocio, seguridad y coordinacion de datos.
4. PostgreSQL: persistencia oficial.
5. Redis: rankings, cache, rate limiting e invitaciones temporales.
6. k6: pruebas externas de carga y estres.

## 7. Diagrama Textual de Arquitectura

```text
Cliente Web
   |
   v
Frontend - Next.js
   |
   v
Nginx - Reverse Proxy / Load Balancer
   |
   |--- Backend Go API instancia 1
   |--- Backend Go API instancia 2
   |--- Backend Go API instancia 3
          |
          |--- PostgreSQL
          |--- Redis
```

## 8. Arquitectura Backend Feature-Based

El backend no debe implementarse como MVC puro. Se usara una arquitectura limpia, ligera y pragmatica basada en funcionalidades.

Patrones obligatorios:

- Feature-Based Architecture.
- Service Layer.
- Repository Pattern.
- DTO Pattern.
- Middleware Pattern.
- Separacion clara entre handler, service y repository.

Estructura recomendada:

```text
backend/
  cmd/
    api/
      main.go
  internal/
    config/
    database/
    redis/
    middleware/
    auth/
      handler.go
      service.go
      repository.go
      dto.go
      model.go
    users/
      handler.go
      service.go
      repository.go
      dto.go
      model.go
    roles/
    rooms/
    matches/
    predictions/
    scores/
    leaderboards/
    audit/
  migrations/
  seeds/
  Dockerfile
  .env.example
```

Reglas backend:

- Los handlers reciben HTTP, validan entrada basica y llaman services.
- Los services contienen reglas de negocio, permisos y coordinacion transaccional.
- Los repositories acceden a PostgreSQL.
- Redis no debe usarse directamente desde handlers.
- PostgreSQL no debe usarse directamente desde handlers.
- La logica de puntuacion vive en `scores`.
- La logica de rankings vive en `leaderboards`.
- La logica de autenticacion vive en `auth`.
- La logica de auditoria vive en `audit`.
- Los errores deben responderse con formato estandar.
- Todas las entradas deben validarse.
- No se deben subir secretos al repositorio.

## 9. Arquitectura Frontend Feature-Based

El frontend usara Next.js App Router y arquitectura por funcionalidades.

Estructura recomendada:

```text
frontend/
  app/
    login/
    register/
    dashboard/
    rooms/
    rooms/[id]/
    matches/
    predictions/
    leaderboard/
    admin/
    layout.tsx
    page.tsx

  features/
    auth/
      components/
      hooks/
      services/
      types/
      schemas/
    rooms/
      components/
      hooks/
      services/
      types/
      schemas/
    matches/
      components/
      hooks/
      services/
      types/
      schemas/
    predictions/
      components/
      hooks/
      services/
      types/
      schemas/
    leaderboards/
      components/
      hooks/
      services/
      types/
    admin/
      components/
      hooks/
      services/
      types/

  shared/
    components/
      ui/
      layout/
      feedback/
    hooks/
    services/
    types/
    lib/
    utils/

  public/
  Dockerfile
  .env.example
```

Reglas frontend:

- Usar TypeScript.
- Usar App Router.
- Separar logica por features.
- No hacer llamadas HTTP directamente dentro de componentes grandes.
- Centralizar llamadas HTTP en services.
- Usar schemas para validacion cuando aplique.
- Proteger rutas por autenticacion y rol.
- Manejar loading, error y empty states.
- Usar shadcn/ui solo cuando sea necesario.
- Usar Recharts solo para graficos y estadisticas.
- Mantener UI moderna, clara y responsive.

## 10. Modulos del Sistema

Modulos principales:

- Auth.
- Users.
- Roles.
- Rooms.
- Room Invites.
- Teams.
- Matches.
- Predictions.
- Scores.
- Leaderboards.
- Audit.
- Admin.
- Security.
- Docker/Nginx.
- k6 tests.

## 11. Roles y Permisos

Roles globales:

- SUPER_ADMIN: control total del sistema.
- ADMIN: gestion de equipos, partidos, resultados, rankings y reportes.
- USER: registro, salas, predicciones y rankings.

Roles por sala:

- OWNER: creador de sala, administra sala, invitaciones, moderadores y cierre.
- MODERATOR: genera invitaciones, ve miembros y puede expulsar MEMBERS.
- MEMBER: participa, predice y consulta rankings.

Reglas:

- Una sala puede tener varios MODERATORS.
- El OWNER puede cerrar la sala.
- No se eliminan fisicamente salas en el MVP; una sala activa usa estado `ACTIVE` y una sala cerrada usa estado `CLOSED`.
- Si en el futuro un administrador elimina fisicamente un registro desde base de datos, eso sera una operacion administrativa distinta, no un estado funcional de la sala.
- Solo usuarios registrados pueden unirse a salas.
- Los permisos globales y de sala deben validarse en services o middleware reutilizable.

## 12. Modelo de Datos Propuesto

Tablas minimas:

- users.
- auth_accounts.
- refresh_tokens.
- roles.
- user_roles.
- rooms.
- room_members.
- room_invites.
- teams.
- matches.
- predictions.
- score_events.
- audit_logs.

Restricciones obligatorias:

```text
UNIQUE(users.email)
UNIQUE(room_members.room_id, room_members.user_id)
UNIQUE(predictions.user_id, predictions.match_id)
```

Estados recomendados:

```text
users.status: ACTIVE, BLOCKED, DISABLED
rooms.status: ACTIVE, CLOSED
room_invites.status: ACTIVE, EXPIRED, REVOKED
matches.status: SCHEDULED, LOCKED, FINISHED, CANCELLED
```

Campos minimos sugeridos:

```text
users: id, full_name, email, avatar_url, password_hash, status, failed_login_attempts, locked_until, created_at, updated_at
auth_accounts: id, user_id, provider, provider_user_id, provider_email, created_at
refresh_tokens: id, user_id, token_hash, user_agent, ip_address, expires_at, revoked_at, created_at
roles: id, name, description
user_roles: user_id, role_id, created_at
rooms: id, name, created_by, status, created_at, updated_at, closed_at
room_members: id, room_id, user_id, room_role, joined_at
room_invites: id, room_id, code_hash, created_by, expires_at, used_count, status, created_at
teams: id, name, country_code, flag_url, created_at
matches: id, home_team_id, away_team_id, match_date, home_score, away_score, status, created_at, updated_at
predictions: id, user_id, match_id, predicted_home_score, predicted_away_score, is_exact_score, is_winner_correct, is_goal_difference_correct, base_points, early_bonus_points, streak_bonus_points, total_points, created_at, updated_at, locked_at
score_events: id, user_id, match_id, prediction_id, points, reason, created_at
audit_logs: id, user_id, action, entity, entity_id, ip_address, user_agent, created_at
```

Reglas para `score_events`:

- Para el MVP, `score_events` registra unicamente eventos reales de puntaje global.
- No debe crear eventos que dupliquen puntos por sala.
- No debe incluir campos para puntajes independientes de sala.
- El calculo de puntos debe poder reconstruirse desde PostgreSQL.

Comportamiento recomendado para `room_members`:

- Para el MVP, cuando un usuario abandona una sala, se puede eliminar su registro activo de `room_members`.
- El evento de salida debe quedar registrado en `audit_logs`.
- Esta decision permite que el usuario pueda volver a unirse despues mediante una invitacion valida sin conflicto con `UNIQUE(room_members.room_id, room_members.user_id)`.
- No se elimina el usuario.
- No se elimina la sala.
- No se eliminan predicciones.
- No se eliminan `score_events`.
- Solo se elimina o desactiva la pertenencia del usuario a la sala.

## 13. Reglas de Negocio

Predicciones:

- Cada usuario registra una sola prediccion por partido.
- La prediccion genera puntos una sola vez sobre el puntaje oficial global del usuario.
- La prediccion no genera puntos separados por sala.
- El usuario solo predice marcador.
- Se permite crear o editar hasta 10 minutos antes del inicio del partido.
- Despues del limite, la prediccion queda bloqueada.

Partidos:

- No se usara API externa en el MVP.
- ADMIN y SUPER_ADMIN crean equipos y partidos.
- ADMIN y SUPER_ADMIN registran resultados finales.
- El calculo de puntajes ocurre despues de registrar resultado y finalizar partido.

Salas:

- Cualquier usuario registrado puede crear sala.
- El creador queda como OWNER.
- OWNER y MODERATOR pueden generar invitaciones.
- Una sala solo puede tener una invitacion activa al mismo tiempo.
- Si OWNER o MODERATOR genera una nueva invitacion para la misma sala, la invitacion activa anterior queda revocada automaticamente.
- Las salas no se eliminan fisicamente en el MVP; solo se cierran con estado `CLOSED`.
- Una sala activa usa estado `ACTIVE`.
- Una sala cerrada usa estado `CLOSED`.
- Si en el futuro un administrador elimina fisicamente un registro desde base de datos, eso sera una operacion administrativa distinta, no un estado funcional de la sala.
- Duraciones permitidas para invitaciones: 1, 3, 5, 10, 15 y 20 minutos.
- No se aceptan valores fuera de esta lista.

- Duracion por defecto: 5 minutos.
- Solo se permite generar una invitacion cada 10 segundos por el mismo usuario.
- El invite no puede ser usado por el mismo usuario que lo creo.
- El invite no puede ser usado por un usuario que ya sea miembro de la sala.
- Solo usuarios registrados pueden usar invitaciones.
- Los codigos de invitacion deben guardarse como hash en PostgreSQL.
- Redis debe usar claves basadas en hash del codigo, no el codigo en texto plano.
- Redis debe aplicar TTL a las invitaciones temporales.
- Un usuario puede abandonar una sala por decision propia.
- Al abandonar una sala, deja de ser miembro activo de esa sala.
- Al abandonar una sala, deja de aparecer en el ranking de esa sala.
- Abandonar una sala no elimina sus predicciones.
- Abandonar una sala no elimina su puntaje global.
- Abandonar una sala no elimina sus `score_events`.
- Abandonar una sala no recalcula puntos.
- Abandonar una sala no elimina la sala.
- Abandonar una sala no afecta su participacion en otras salas.
- Si vuelve a unirse a la sala, debe hacerlo mediante una invitacion valida.
- Al volver a unirse, aparece con su puntaje global acumulado actual.
- El ranking de sala refleja el puntaje global acumulado actual de cada miembro.
- Si un usuario entra a una sala despues de haber acumulado puntos, aparece con su puntaje global actual.
- Si un usuario abandona una sala, deja de aparecer en el ranking de esa sala.
- Si vuelve a unirse a una sala, aparece nuevamente con su puntaje global acumulado actual.
- No se recalculan puntos antiguos por sala.
- No se duplican puntos por pertenecer a varias salas.
- La sala solo muestra una vista filtrada de sus miembros y sus puntos globales.

## 14. Reglas de Puntuacion

La puntuacion es acumulativa:

- Marcador exacto: +5 puntos.
- Ganador correcto o empate correcto: +3 puntos.
- Diferencia de goles correcta: +2 puntos.
- Prediccion anticipada: +1 punto si fue registrada con mas de 24 horas de anticipacion.
- Bonus por racha: +2 puntos por cada 3 partidos consecutivos acertando ganador o empate.

Diferencia de goles:

```text
predicted_diff = predicted_home_score - predicted_away_score
real_diff = home_score - away_score
```

La diferencia es correcta solo si:

```text
predicted_diff == real_diff
```

Reglas de racha:

- 3 aciertos consecutivos: +2.
- 6 aciertos consecutivos: +4.
- 9 aciertos consecutivos: +6.
- Si falla, la racha se reinicia.
- Marcador exacto tambien cuenta como acierto de ganador/empate.

Reglas de puntaje unico y rankings:

- El usuario tiene un unico puntaje acumulado global.
- La prediccion es unica por usuario y partido.
- Los puntos se calculan una sola vez por prediccion.
- Las salas no generan puntos independientes.
- Las salas no tienen una billetera separada.
- El ranking global muestra a todos los usuarios ordenados por su puntaje global acumulado.
- El ranking de sala muestra unicamente a los miembros activos de esa sala, ordenados por el mismo puntaje global acumulado.
- Si un usuario pertenece a varias salas, su puntaje es el mismo en todas.
- Los puntos no se duplican por pertenecer a varias salas.
- Si un usuario entra a una sala despues de haber acumulado puntos, aparece con su puntaje global actual.
- Si un usuario abandona una sala, deja de aparecer en el ranking de esa sala.
- Si vuelve a unirse a una sala, aparece nuevamente con su puntaje global acumulado actual.
- Redis puede mantener leaderboards por sala para optimizar consultas, pero esos leaderboards solo son una vista o cache del puntaje global.
- PostgreSQL sigue siendo la fuente oficial de verdad.

Ejemplo esperado:

```text
Usuario Diego:
- Puntaje global: 120 puntos.
- Pertenece a Sala Amigos.
- Pertenece a Sala Tecsup.

Resultado:
- En ranking global aparece con 120 puntos.
- En Sala Amigos aparece con 120 puntos.
- En Sala Tecsup aparece con 120 puntos.
- No se duplica a 240 puntos.
- No se calcula un puntaje independiente por sala.
```

## 15. Reglas de Seguridad

Autenticacion:

- Access token con duracion de 1 hora.
- Refresh token revocable.
- Refresh token almacenado como hash en PostgreSQL.
- Refresh token no debe almacenarse en `localStorage`.
- En frontend debe preferirse cookie `HttpOnly`, `Secure` y `SameSite` cuando aplique, o almacenamiento en memoria para access tokens de corta duracion.
- Backend siempre genera sus propios JWT.
- Google solo valida identidad.

Contrasenas:

- Usar bcrypt o Argon2id.
- Nunca guardar contrasenas en texto plano.

Sesiones:

- Logout revoca refresh token.
- Cambio de contrasena revoca todas las sesiones.
- 5 intentos fallidos bloquean login por 15 minutos.

Protecciones:

- Rate limiting con Redis.
- Validacion de entradas.
- CORS configurado explicitamente.
- Variables de entorno por capa.
- No exponer secretos al frontend.
- No subir `.env` reales.
- Registrar eventos relevantes en audit_logs.

## 16. Estrategia de Redis

Redis se usara como capa de optimizacion, no como fuente oficial de datos.

Claves recomendadas:

```text
leaderboard:global
leaderboard:room:{room_id}
cache:matches:upcoming
rate_limit:login:{ip}
rate_limit:user:{user_id}
invite_code:{code_hash}
```

Uso:

- Sorted Sets para rankings.
- TTL para invitaciones.
- TTL para rate limiting.
- Cache para proximos partidos.
- PostgreSQL seguira siendo la fuente de verdad.
- Los codigos de invitacion deben guardarse como hash en PostgreSQL.
- Redis debe usar claves basadas en hash, nunca el codigo plano.
- Los rankings deben poder reconstruirse desde PostgreSQL si Redis se reinicia o pierde datos.
- Los leaderboards por sala en Redis son vistas/cache del puntaje global filtrado por miembros activos de la sala.
- Redis no debe representar una billetera ni un puntaje independiente por sala.

## 17. Estrategia de Docker y Nginx

Servicios minimos:

- frontend.
- backend.
- postgres.
- redis.
- nginx.
- k6.

Reglas:

- Nginx sera el unico punto de entrada externo.
- Backend no debe exponerse directamente al host en escenario con balanceo.
- El sistema debe soportar escalamiento horizontal del backend.
- Docker Compose debe permitir `docker compose up --scale backend=3`.
- Nginx debe distribuir trafico entre replicas usando Round Robin.

## 18. Estrategia de Pruebas con k6

Crear carpeta:

```text
tests/k6/
```

Scripts requeridos:

- Registro de usuarios.
- Login.
- Consulta de partidos.
- Creacion de salas.
- Union a salas por codigo.
- Registro concurrente de predicciones.
- Consulta de ranking global.
- Consulta de ranking por sala.
- Comparacion con 1 backend vs 3 backends.

Metricas a documentar:

- http_req_duration.
- p95.
- p99.
- requests por segundo.
- porcentaje de errores.
- usuarios virtuales concurrentes.
- comportamiento antes y despues del escalamiento.

## 19. Estrategia de Auditoria

Eventos minimos:

```text
USER_REGISTERED
USER_LOGIN
USER_LOGIN_GOOGLE
USER_LOGOUT
FAILED_LOGIN_ATTEMPT
ROOM_CREATED
ROOM_INVITE_CREATED
ROOM_JOINED
ROOM_MEMBER_REMOVED
ROOM_MODERATOR_ASSIGNED
ROOM_CLOSED
ROOM_LEFT
PREDICTION_CREATED
PREDICTION_UPDATED
MATCH_CREATED
MATCH_RESULT_UPDATED
SCORE_CALCULATED
REFRESH_TOKEN_REVOKED
PASSWORD_CHANGED
```

Cada evento debe guardar:

- user_id cuando exista.
- action.
- entity.
- entity_id.
- ip_address.
- user_agent.
- created_at.

Evento `ROOM_LEFT`:

- Se registra cuando un usuario abandona voluntariamente una sala.
- Debe guardar `user_id`.
- Debe guardar `action = ROOM_LEFT`.
- Debe guardar `entity = ROOM`.
- Debe guardar `entity_id = room_id`.
- Debe guardar `ip_address`.
- Debe guardar `user_agent`.
- Debe guardar `created_at`.

## 20. Skills Recomendados para el Agente

No se deben instalar skills todavia.

Skills recomendados para futuras fases:

- Skill para Go API / Gin.
- Skill para Next.js / React / TypeScript.
- Skill para PostgreSQL.
- Skill para Docker / Nginx.
- Skill para seguridad / JWT / OAuth.
- Skill para k6.
- Skill para revision de arquitectura limpia.
- Skill para documentacion tecnica.

## 21. Reglas Complementarias Antes de Codificar

Antes de iniciar cualquier implementacion, todos los agentes o desarrolladores deben cumplir estas reglas:

- Leer `AGENTS.md` si existe antes de modificar el repositorio.
- Leer y respetar `PLAN_PROYECTO_APUESTEC.md` antes de crear, modificar o eliminar codigo.
- Mantener todo lo definido sobre stack, arquitectura, seguridad y reglas de negocio.
- Las salas no se eliminan fisicamente en el MVP; solo se cierran con estado `CLOSED`.
- Las invitaciones solo aceptan duraciones de 1, 3, 5, 10, 15 y 20 minutos.
- Los codigos de invitacion deben guardarse como hash en PostgreSQL.
- Redis debe usar claves de invitacion basadas en hash, no codigos en texto plano.
- Redis no es fuente de verdad.
- Los rankings deben poder reconstruirse desde PostgreSQL.
- `score_events` debe registrar solo eventos reales de puntaje global en el MVP.
- Las salas no generan puntajes independientes ni duplican puntos.
- El ranking de sala debe reflejar el puntaje global acumulado actual de sus miembros activos.
- El refresh token no debe almacenarse en `localStorage`.
- Todo modulo nuevo debe respetar el contrato minimo de endpoints REST definido en este documento.
- Cada fase debe cumplir sus criterios de aceptacion antes de avanzar a la siguiente.

## 22. Contrato Minimo de Endpoints REST por Modulo

Convenciones generales:

- Prefijo recomendado para API: `/api/v1`.
- Todas las respuestas deben usar formato estandar.
- Endpoints privados requieren `Authorization: Bearer <access_token>` salvo que se use cookie segura para sesion.
- Los errores deben incluir codigo, mensaje y detalles validables cuando aplique.
- Los endpoints administrativos requieren validacion de rol global.
- Los endpoints de sala requieren validacion de rol dentro de la sala cuando corresponda.

Auth:

- `POST /api/v1/auth/register`: registra usuario local.
- `POST /api/v1/auth/login`: inicia sesion local.
- `POST /api/v1/auth/google`: valida identidad Google y emite JWT propio.
- `POST /api/v1/auth/refresh`: renueva access token usando refresh token valido.
- `POST /api/v1/auth/logout`: revoca refresh token activo.
- `POST /api/v1/auth/change-password`: cambia contrasena y revoca sesiones previas.

Users:

- `GET /api/v1/users/me`: obtiene perfil del usuario autenticado.
- `PATCH /api/v1/users/me`: actualiza perfil propio.
- `GET /api/v1/admin/users`: lista usuarios para ADMIN o SUPER_ADMIN.
- `PATCH /api/v1/admin/users/{user_id}/status`: cambia estado de usuario.
- `PATCH /api/v1/admin/users/{user_id}/roles`: actualiza roles globales.

Rooms:

- `POST /api/v1/rooms`: crea sala y asigna OWNER.
- `GET /api/v1/rooms`: lista salas del usuario autenticado.
- `GET /api/v1/rooms/{room_id}`: obtiene detalle de sala.
- `PATCH /api/v1/rooms/{room_id}`: actualiza datos permitidos de sala.
- `POST /api/v1/rooms/{room_id}/close`: cierra sala con estado `CLOSED`.
- `GET /api/v1/rooms/{room_id}/members`: lista miembros de sala.
- `DELETE /api/v1/rooms/{room_id}/members/me`: permite que el usuario autenticado abandone una sala.
- `DELETE /api/v1/rooms/{room_id}/members/{user_id}`: expulsa miembro segun permisos.
- `PATCH /api/v1/rooms/{room_id}/members/{user_id}/role`: cambia rol de sala segun permisos.

Descripcion funcional de abandonar sala:

- El usuario autenticado puede abandonar una sala donde actualmente es miembro.
- Al abandonar la sala, deja de aparecer en el ranking de esa sala.
- Su puntaje global no se modifica.
- Sus predicciones no se eliminan.
- Sus `score_events` no se eliminan.
- No se recalculan puntos.
- No se elimina la sala.
- No se afecta su participacion en otras salas.
- Si luego vuelve a unirse mediante una invitacion valida, aparece nuevamente con su puntaje global acumulado actual.

Room Invites:

- `POST /api/v1/rooms/{room_id}/invites`: genera invitacion temporal.
- `GET /api/v1/rooms/{room_id}/invites/active`: consulta invitacion activa permitida.
- `POST /api/v1/rooms/join`: une usuario autenticado usando codigo valido.
- `POST /api/v1/rooms/{room_id}/invites/revoke`: revoca invitacion activa.

Teams:

- `POST /api/v1/admin/teams`: crea equipo.
- `GET /api/v1/teams`: lista equipos.
- `GET /api/v1/teams/{team_id}`: obtiene equipo.
- `PATCH /api/v1/admin/teams/{team_id}`: actualiza equipo.

Matches:

- `POST /api/v1/admin/matches`: crea partido.
- `GET /api/v1/matches`: lista partidos.
- `GET /api/v1/matches/{match_id}`: obtiene partido.
- `PATCH /api/v1/admin/matches/{match_id}`: actualiza partido antes de iniciar.
- `POST /api/v1/admin/matches/{match_id}/result`: registra resultado oficial.
- `POST /api/v1/admin/matches/{match_id}/finish`: marca partido como finalizado.
- `POST /api/v1/admin/matches/{match_id}/score`: activa calculo de puntajes.

Predictions:

- `POST /api/v1/predictions`: crea prediccion.
- `GET /api/v1/predictions/me`: lista predicciones del usuario autenticado.
- `GET /api/v1/matches/{match_id}/prediction`: obtiene prediccion propia para un partido.
- `PATCH /api/v1/predictions/{prediction_id}`: edita prediccion antes del bloqueo.

Scores:

- `GET /api/v1/scores/me`: obtiene resumen de puntaje propio.
- `GET /api/v1/admin/matches/{match_id}/scores`: lista puntajes calculados de un partido.
- `POST /api/v1/admin/scores/rebuild-leaderboards`: reconstruye rankings desde PostgreSQL.

Leaderboards:

- `GET /api/v1/leaderboards/global`: consulta ranking global.
- `GET /api/v1/rooms/{room_id}/leaderboard`: consulta ranking de sala.

Audit:

- `GET /api/v1/admin/audit-logs`: consulta eventos de auditoria.
- `GET /api/v1/admin/audit-logs/{audit_log_id}`: obtiene detalle de evento.

Health:

- `GET /api/v1/health`: verifica estado basico de API.
- `GET /api/v1/health/dependencies`: verifica PostgreSQL y Redis.

## 23. Criterios de Aceptacion por Fase

Fase 1: Planificacion y documentacion base

- `PLAN_PROYECTO_APUESTEC.md` esta completo y actualizado.
- `AGENTS.md` fue leido y respetado.
- El stack tecnologico no fue cambiado.
- Las reglas complementarias antes de codificar estan documentadas.

Fase 2: Infraestructura base

- `docker-compose.yml` levanta los servicios minimos.
- Nginx es el unico punto de entrada externo.
- PostgreSQL y Redis tienen health checks o validacion equivalente.
- El backend puede escalarse con `docker compose up --scale backend=3`.

Fase 3: Backend base

- Existe health check funcional.
- La configuracion se carga desde variables de entorno.
- Hay conexion validada a PostgreSQL y Redis.
- Existe manejo estandar de errores.

Fase 4: Base de datos

- Las migraciones crean todas las tablas minimas.
- Las restricciones unicas obligatorias existen.
- `score_events` registra unicamente eventos reales de puntaje global con los campos definidos para el MVP.
- Los seeds minimos permiten iniciar pruebas locales.

Fase 5: Autenticacion y usuarios

- Registro local funcional.
- Login local funcional.
- Login con Google emite JWT propio del backend.
- Refresh token se guarda hasheado y no se almacena en `localStorage`.
- Logout revoca refresh token.
- Bloqueo temporal por intentos fallidos funciona.

Fase 6: Roles y permisos

- Roles globales se validan en endpoints administrativos.
- Roles por sala se validan en operaciones de sala.
- No existe logica de autorizacion duplicada innecesariamente.

Fase 7: Salas e invitaciones

- Crear sala asigna OWNER.
- Cerrar sala usa estado `CLOSED` y no borrado fisico.
- Las salas solo usan estados `ACTIVE` y `CLOSED`.
- No existe estado `DELETED` en el MVP.
- Solo existe una invitacion activa por sala.
- Una nueva invitacion revoca la anterior.
- Invitaciones aceptan solo 1, 3, 5, 10, 15 y 20 minutos.
- La duracion por defecto es 5 minutos.
- Codigos se guardan hasheados y Redis usa clave basada en hash.
- Redis aplica TTL a invitaciones temporales.
- Un usuario ya miembro no puede reutilizar invitacion.
- El creador de la invitacion no puede usar su propio codigo.
- El usuario autenticado puede abandonar una sala usando `DELETE /api/v1/rooms/{room_id}/members/me`.
- Al abandonar una sala, deja de aparecer en el ranking de esa sala.
- Al abandonar una sala, su puntaje global no cambia.
- Al abandonar una sala, sus predicciones no se eliminan.
- Al abandonar una sala, sus `score_events` no se eliminan.
- Abandonar una sala no afecta su participacion en otras salas.
- El evento `ROOM_LEFT` queda registrado en `audit_logs`.
- Si el usuario vuelve a unirse mediante invitacion valida, aparece nuevamente con su puntaje global acumulado actual.

Fase 8: Partidos y equipos

- ADMIN o SUPER_ADMIN puede gestionar equipos.
- ADMIN o SUPER_ADMIN puede gestionar partidos.
- No se puede editar indebidamente un partido bloqueado o finalizado.
- Se puede registrar resultado oficial.

Fase 9: Predicciones

- Cada usuario solo tiene una prediccion por partido.
- La edicion se bloquea 10 minutos antes del inicio.
- Los inputs de marcador son validados.
- El sistema impide predicciones duplicadas.

Fase 10: Puntuacion y rankings

- Las reglas de puntos se calculan correctamente.
- El sistema maneja un unico puntaje global por usuario.
- Las salas no generan puntajes independientes.
- La prediccion es unica por usuario y partido.
- Los puntos se calculan una sola vez.
- Se generan eventos `score_events` solo para puntaje global oficial.
- No se crean eventos que dupliquen puntos por sala.
- El ranking global muestra el puntaje global acumulado.
- El ranking de sala muestra miembros de la sala ordenados por su puntaje global acumulado.
- Los puntos no se duplican por pertenecer a varias salas.
- Si un usuario entra a una sala despues de tener puntos, aparece con su puntaje global actual.
- Si un usuario abandona una sala, deja de aparecer en el ranking de esa sala.
- Redis actualiza rankings con Sorted Sets.
- Redis puede cachear leaderboards por sala, pero PostgreSQL sigue siendo la fuente oficial.
- Los rankings pueden reconstruirse desde PostgreSQL.

Fase 11: Frontend

- Las rutas principales existen y son responsive.
- Las rutas protegidas validan autenticacion y rol.
- No se almacena refresh token en `localStorage`.
- Las pantallas manejan loading, error y empty states.

Fase 12: k6 y documentacion

- Existen scripts k6 para flujos principales.
- Se documentan metricas p95, p99, errores y throughput.
- Se comparan resultados con 1 backend y 3 backends.
- La documentacion explica ejecucion local y despliegue Linux.

## 24. Plan de Desarrollo por Fases

Fase 1: Planificacion y documentacion base

- Completar `PLAN_PROYECTO_APUESTEC.md`.
- Definir estructura de carpetas.
- Definir convenciones.
- Revisar dependencias existentes.
- Completar `.env.example` por capa.

Fase 2: Infraestructura base

- Crear Dockerfile de backend.
- Crear Dockerfile de frontend.
- Completar docker-compose.yml.
- Completar nginx.conf.
- Configurar PostgreSQL y Redis.
- Validar ejecucion local.

Fase 3: Backend base

- Crear estructura `cmd/api` e `internal`.
- Configurar Gin.
- Configurar carga de variables de entorno.
- Configurar conexion PostgreSQL.
- Configurar conexion Redis.
- Crear middleware de errores, logging, CORS y rate limit.
- Crear respuesta estandar de API.

Fase 4: Base de datos

- Crear migraciones iniciales.
- Crear seeds de roles.
- Crear seeds de equipos y partidos.
- Validar restricciones unicas.
- Validar relaciones.

Fase 5: Autenticacion y usuarios

- Registro local.
- Login local.
- Login con Google.
- JWT propio.
- Refresh tokens revocables.
- Logout.
- Bloqueo temporal por intentos fallidos.
- Auditoria de eventos de auth.

Fase 6: Roles y permisos

- Middleware de autenticacion.
- Middleware de rol global.
- Validacion de permisos por sala.
- Services reutilizables de autorizacion.

Fase 7: Salas e invitaciones

- Crear sala.
- Listar salas del usuario.
- Gestionar miembros.
- Crear invitacion temporal.
- Generar QR en frontend.
- Unirse por codigo.
- Abandonar sala por decision propia.
- Invalidar invitaciones previas.

Fase 8: Partidos y equipos

- CRUD administrativo de equipos.
- CRUD administrativo de partidos.
- Registro de resultados.
- Bloqueo de partidos.
- Finalizacion de partidos.

Fase 9: Predicciones

- Crear prediccion.
- Editar prediccion.
- Bloquear edicion 10 minutos antes.
- Validar UNIQUE(user_id, match_id).
- Consultar predicciones del usuario.

Fase 10: Puntuacion y rankings

- Implementar calculo de puntos.
- Crear score_events.
- Actualizar rankings en Redis.
- Consultar ranking global.
- Consultar ranking por sala.
- Sincronizar o reconstruir rankings desde PostgreSQL si Redis se reinicia.

Fase 11: Frontend

- Layout base.
- Login y registro.
- Dashboard.
- Salas.
- Invitaciones.
- Partidos.
- Predicciones.
- Ranking global.
- Ranking por sala.
- Panel admin.
- Estados de carga, error y vacio.

Fase 12: k6 y documentacion

- Crear scripts k6.
- Ejecutar pruebas con 1 backend.
- Ejecutar pruebas con 3 backends.
- Documentar resultados.
- Documentar uso de Docker Compose.
- Documentar despliegue basico en Linux.

## 25. Checklist de Entregables

- Propuesta de arquitectura del sistema.
- Aplicacion funcional.
- Dockerfile backend.
- Dockerfile frontend.
- docker-compose.yml.
- nginx.conf.
- Escalamiento horizontal del backend.
- Balanceo con Nginx.
- PostgreSQL configurado.
- Redis configurado.
- Sistema de autenticacion.
- Google OAuth como alternativa.
- Roles y permisos.
- Salas privadas.
- Invitaciones temporales.
- Sistema de predicciones.
- Sistema de puntuacion.
- Rankings globales y por sala.
- Auditoria.
- Pruebas k6.
- Documentacion de pruebas.
- Documentacion de ejecucion local.
- Documentacion de despliegue Linux.
- Variables de entorno por capa.

## 26. Riesgos Tecnicos y Mitigaciones

Riesgo: mezcla de logica en handlers.

Mitigacion: services obligatorios para reglas de negocio y repositories para persistencia.

Riesgo: inconsistencias entre PostgreSQL y Redis.

Mitigacion: PostgreSQL sera fuente de verdad y Redis podra reconstruirse.

Riesgo: refresh tokens comprometidos.

Mitigacion: guardar solo hash, permitir revocacion y registrar auditoria.

Riesgo: mala gestion de permisos por sala.

Mitigacion: centralizar autorizacion en services o helpers reutilizables.

Riesgo: problemas de escalamiento por estado en backend.

Mitigacion: backend stateless, sesiones via JWT y datos compartidos en PostgreSQL/Redis.

Riesgo: invitaciones reutilizadas o vencidas.

Mitigacion: TTL en Redis, hash en PostgreSQL y revocacion de codigos anteriores.

Riesgo: pruebas k6 poco representativas.

Mitigacion: probar flujos reales y comparar 1 backend vs 3 backends.

## 27. Buenas Practicas de Codigo

Backend:

- Handlers pequenos.
- Services con logica de negocio.
- Repositories sin logica de negocio.
- DTOs separados de modelos internos.
- Errores estandarizados.
- Validacion de inputs.
- Transacciones cuando haya multiples escrituras relacionadas.
- Logs sin secretos.
- Tests para logica critica de puntuacion y permisos.

Frontend:

- Componentes pequenos.
- Services por feature.
- Types por feature.
- Hooks para logica reutilizable.
- Validacion de formularios.
- Rutas protegidas.
- UI responsive.
- Estados de carga, error y vacio.

General:

- No subir secretos.
- No instalar dependencias innecesarias.
- No duplicar logica.
- No mezclar responsabilidades.
- Mantener documentacion actualizada.

## 28. Convenciones de Nombres

Backend:

- Carpetas en minuscula plural para features: `users`, `rooms`, `matches`.
- Archivos estandar: `handler.go`, `service.go`, `repository.go`, `dto.go`, `model.go`.
- Variables de entorno en mayusculas con snake case.
- Endpoints REST en plural y con convencion consistente.

Frontend:

- Features en minuscula plural.
- Componentes React en PascalCase.
- Hooks con prefijo `use`.
- Tipos en PascalCase.
- Services con nombres descriptivos por dominio.

Base de datos:

- Tablas en snake_case plural.
- Columnas en snake_case.
- Indices con nombres descriptivos.
- Estados en uppercase.

## 29. Estructura de Carpetas Recomendada

```text
ApuesTec/
  backend/
    cmd/
      api/
    internal/
      auth/
      users/
      roles/
      rooms/
      matches/
      predictions/
      scores/
      leaderboards/
      audit/
      config/
      database/
      redis/
      middleware/
    migrations/
    seeds/
    Dockerfile
    .env.example

  frontend/
    app/
    features/
    shared/
    public/
    Dockerfile
    .env.example

  database/
    .env.example

  nginx/
    nginx.conf

  tests/
    k6/

  docker-compose.yml
  README.md
  PLAN_PROYECTO_APUESTEC.md
```

## 30. Orden Recomendado para Empezar a Codificar

1. Completar documentacion del plan.
2. Completar estructura de carpetas base.
3. Configurar Docker Compose, PostgreSQL, Redis y Nginx.
4. Crear backend minimo con health check.
5. Crear frontend minimo conectado al backend.
6. Crear migraciones y seeds.
7. Implementar autenticacion.
8. Implementar roles y permisos.
9. Implementar salas e invitaciones.
10. Implementar partidos.
11. Implementar predicciones.
12. Implementar puntuacion.
13. Implementar rankings.
14. Implementar frontend completo por features.
15. Implementar k6.
16. Documentar pruebas y despliegue.

## 31. Documentacion Obligatoria

La documentacion del proyecto debera cubrir:

- Arquitectura del sistema.
- Stack tecnologico.
- Modelo de base de datos.
- Roles y permisos.
- Reglas de puntuacion.
- Seguridad implementada.
- Uso de Docker Compose.
- Escalamiento con Nginx.
- Pruebas con k6.
- Como ejecutar el proyecto localmente.
- Como preparar el despliegue en servidor Linux.
- Variables de entorno por capa.
- Arquitectura backend por features.
- Arquitectura frontend por features.

## 32. Restricciones Obligatorias

- No usar Firebase Auth.
- No usar API externa de partidos en el MVP.
- No usar dinero real.
- No implementar pagos.
- No implementar cuotas reales.
- No crear endpoints sin separacion de responsabilidades.
- No acceder a PostgreSQL desde handlers.
- No acceder a Redis desde handlers.
- No omitir seguridad.
- No omitir Docker.
- No omitir Nginx.
- No omitir Redis.
- No omitir k6.
- No crear skills manuales dentro del repositorio.
- No instalar dependencias innecesarias.
- No desarrollar pantallas o endpoints antes de completar el plan.
