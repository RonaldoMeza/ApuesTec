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
- No se eliminan fisicamente salas; se usa estado CLOSED o DELETED.
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
rooms.status: ACTIVE, CLOSED, DELETED
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

## 13. Reglas de Negocio

Predicciones:

- Cada usuario registra una sola prediccion por partido.
- La prediccion aplica al ranking global.
- La prediccion aplica a todos los rankings de salas donde el usuario sea miembro.
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
- Una nueva invitacion invalida la anterior.
- Las invitaciones tienen duracion configurable entre 1 y 20 minutos.
- Duracion por defecto: 5 minutos.

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

## 15. Reglas de Seguridad

Autenticacion:

- Access token con duracion de 1 hora.
- Refresh token revocable.
- Refresh token almacenado como hash en PostgreSQL.
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
invite_code:{code}
```

Uso:

- Sorted Sets para rankings.
- TTL para invitaciones.
- TTL para rate limiting.
- Cache para proximos partidos.
- PostgreSQL seguira siendo la fuente de verdad.

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

## 21. Plan de Desarrollo por Fases

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

## 22. Checklist de Entregables

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

## 23. Riesgos Tecnicos y Mitigaciones

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

## 24. Buenas Practicas de Codigo

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

## 25. Convenciones de Nombres

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

## 26. Estructura de Carpetas Recomendada

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

## 27. Orden Recomendado para Empezar a Codificar

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

## 28. Documentacion Obligatoria

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

## 29. Restricciones Obligatorias

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
