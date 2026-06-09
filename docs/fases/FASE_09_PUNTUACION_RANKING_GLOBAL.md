# Fase 9: Puntuacion y Ranking Global

## Objetivo

Implementar el calculo real de puntos cuando un partido finaliza, mostrar puntuaciones por prediccion, implementar ranking global, estadisticas de usuario y mantener la arquitectura preparada para rankings por sala.

## Archivos Creados

### Backend

| Archivo | Proposito |
|---------|-----------|
| `backend/migrations/004_scoring.sql` | Migracion con tablas `score_events` y `user_scores` |
| `backend/internal/database/migrate.go` | Runner automatico de migraciones SQL al iniciar el backend |
| `backend/internal/scoring/model.go` | Modelos ScoreEvent, UserScore, ScoredPrediction |
| `backend/internal/scoring/dto.go` | DTOs para ScoreEventResponse |
| `backend/internal/scoring/repository.go` | Repositorio para calculo de puntuacion |
| `backend/internal/scoring/service.go` | Logica de calculo de puntos (exacto, ganador, diferencia, anticipado, racha) |
| `backend/internal/leaderboard/model.go` | Modelo LeaderboardEntry |
| `backend/internal/leaderboard/dto.go` | DTOs para respuesta del ranking global |
| `backend/internal/leaderboard/repository.go` | Consultas de ranking global y posicion de usuario |
| `backend/internal/leaderboard/service.go` | Servicio de ranking global con limite configurable |
| `backend/internal/leaderboard/handler.go` | Handler GET /api/v1/leaderboard/global |
| `backend/internal/userstats/dto.go` | DTOs para estadisticas de usuario |
| `backend/internal/userstats/repository.go` | Consultas de puntaje, rachas y eventos |
| `backend/internal/userstats/service.go` | Logica de estadisticas de usuario |
| `backend/internal/userstats/handler.go` | Handlers GET /api/v1/users/me/stats y /api/v1/users/me/score-events |

### Archivos Modificados

| Archivo | Cambio |
|---------|--------|
| `backend/cmd/api/main.go` | Agregada llamada a `database.RunMigrations` al iniciar el servidor |
| `backend/internal/matches/service.go` | Agregada integracion con scoring via `ScoringFunc` |
| `backend/internal/matches/handler.go` | Agregado handler `RecalculateScore` |
| `backend/internal/predictions/dto.go` | Agregados campos de puntuacion detallada al response |
| `backend/internal/audit/repository.go` | Agregadas constantes `MATCH_RESULT_UPDATED` y `SCORE_CALCULATED` |
| `backend/internal/routes/router.go` | Registradas nuevas rutas de leaderboard, userstats y scoring |

### Frontend

| Archivo | Proposito |
|---------|-----------|
| `frontend/features/leaderboard/types/leaderboard.types.ts` | Tipos LeaderboardEntry, UserStats, ScoreEvent |
| `frontend/features/leaderboard/services/leaderboard.service.ts` | Servicio de llamadas API |
| `frontend/features/leaderboard/components/LeaderboardTable.tsx` | Tabla de ranking global |
| `frontend/features/leaderboard/components/LeaderboardCard.tsx` | Card individual de ranking |
| `frontend/features/leaderboard/components/UserStatsCard.tsx` | Card de estadisticas de usuario |
| `frontend/features/leaderboard/components/ScoreEventList.tsx` | Lista de eventos de puntuacion |
| `frontend/app/leaderboard/page.tsx` | Pagina publica /leaderboard |
| `frontend/app/stats/page.tsx` | Pagina protegida /stats |

### Paginas Modificadas

| Pagina | Cambio |
|--------|--------|
| `frontend/app/dashboard/page.tsx` | Puntos reales, ranking posicion, preview ranking, score events |
| `frontend/app/admin/matches/page.tsx` | Boton "Recalcular puntuacion" y mensaje de exito |
| `frontend/app/page.tsx` | RankingPreview con datos reales |
| `frontend/features/predictions/components/PredictionCard.tsx` | Badges de exacto, ganador, diferencia, anticipado, racha |
| `frontend/features/predictions/components/PredictionForm.tsx` | Puntuacion detallada en partidos finalizados |
| `frontend/features/predictions/types/prediction.types.ts` | Nuevos campos de puntuacion |
| `frontend/features/matches/services/match.service.ts` | Metodo `recalculateScore` |
| `frontend/shared/components/AppLayout.tsx` | Link "Ranking" en navegacion publica |
| `frontend/shared/components/UserNav.tsx` | Link "Mis estadisticas" en menu de usuario |

## Tablas Agregadas

### score_events

Almacena cada evento de puntuacion generado por prediccion.

| Columna | Tipo | Descripcion |
|---------|------|-------------|
| id | UUID PK | Identificador unico |
| user_id | UUID FK | Usuario que recibe los puntos |
| match_id | UUID FK | Partido asociado |
| prediction_id | UUID FK | Prediccion que genero el evento |
| event_type | VARCHAR(50) | EXACT_SCORE, WINNER_CORRECT, GOAL_DIFFERENCE_CORRECT, EARLY_BONUS, STREAK_BONUS |
| points | INTEGER | Puntos del evento |
| description | TEXT | Descripcion del evento |
| created_at | TIMESTAMPTZ | Fecha de creacion |

### user_scores

Materializa el puntaje total acumulado por usuario.

| Columna | Tipo | Descripcion |
|---------|------|-------------|
| user_id | UUID PK/FK | Usuario |
| total_points | INTEGER | Puntos totales acumulados |
| predictions_count | INTEGER | Predicciones con puntos |
| exact_scores_count | INTEGER | Marcadores exactos acertados |
| winner_correct_count | INTEGER | Ganadores/empates correctos |
| goal_difference_correct_count | INTEGER | Diferencias correctas |
| streak_count | INTEGER | Contador de racha |
| last_scored_at | TIMESTAMPTZ | Ultima vez que recibio puntos |
| updated_at | TIMESTAMPTZ | Ultima actualizacion |

## Endpoints Implementados

### Publicos

- `GET /api/v1/leaderboard/global?limit=50` - Ranking global ordenado por puntos descendente

### Protegidos (requiere autenticacion)

- `GET /api/v1/users/me/stats` - Estadisticas del usuario autenticado
- `GET /api/v1/users/me/score-events` - Historial de eventos de puntuacion

### Admin (requiere ADMIN o SUPER_ADMIN)

- `POST /api/v1/admin/matches/:id/recalculate-score` - Recalcula puntuacion de un partido

## Reglas de Puntuacion Implementadas

1. **Marcador exacto**: +5 puntos (prediccion identica al resultado real)
2. **Ganador o empate correcto**: +3 puntos (quien gana o si hay empate)
3. **Diferencia de goles correcta**: +2 puntos (diferencia con signo exacta)
4. **Prediccion anticipada**: +1 punto (si se predijo mas de 24h antes del inicio del partido)
5. **Racha**: +2 puntos por cada 3 partidos consecutivos acertando ganador/empate
6. Los puntos son acumulativos
7. Partidos CANCELLED no otorgan puntos

## Migracion Automatica

- El backend ejecuta migraciones SQL automaticamente al iniciar mediante `internal/database/migrate.go`
- Lee archivos `.sql` del directorio `migrations/` ordenados por prefijo numerico
- Rastrea migraciones aplicadas en la tabla `schema_migrations`
- La migracion `004_scoring.sql` usa `CREATE TABLE IF NOT EXISTS` y `ALTER TABLE ... ADD COLUMN IF NOT EXISTS` para ser compatible con esquemas existentes (la columna `reason` preexistente se migra a `event_type`/`description`)
- Esto evita errores 500 por tablas faltantes al desplegar en entornos nuevos

## Idempotencia del Calculo

- Al registrar resultado, se calculan puntos una sola vez (verificando que no existan score_events previos para esa prediccion)
- Al recalcular, se eliminan score_events del partido, se resetean las predicciones a 0, se reconstruyen los puntajes de usuario desde cero, y se vuelve a calcular
- Esto asegura que el recalculamiento no duplique puntos

## PostgreSQL como Fuente de Verdad

- `user_scores` se mantiene actualizado via upsert durante el calculo
- `RebuildAllUserScores` permite reconstruir todos los puntajes desde las predicciones y score_events
- El ranking global se consulta directamente desde PostgreSQL con `ROW_NUMBER()`
- El calculo de rachas (current/best streak) usa gaps-and-islands con window functions (`ROW_NUMBER`, subconsultas correlacionadas)
- Redis no se utilizo en esta fase; el sistema funciona completamente sin Redis

## Pruebas Realizadas

- `go build ./...` - Backend compila sin errores
- `go test ./...` - Tests existentes pasan
- `npm run lint` - Frontend sin errores de linter
- `npm run build` - Frontend compila sin errores

## Como Probar desde Frontend

1. Entrar a `/leaderboard` sin login - Ver ranking global o empty state
2. Iniciar sesion como USER
3. Ir a `/matches` y crear predicciones para partidos SCHEDULED
4. Iniciar sesion como ADMIN
5. Ir a `/admin/matches`
6. Registrar resultado de un partido LOCKED
7. Confirmar que el partido pasa a FINISHED y que se recalculan puntos
8. Iniciar sesion como USER
9. Ir a `/predictions` - Ver puntos reales y badges
10. Ir a `/stats` - Ver estadisticas detalladas
11. Ir a `/dashboard` - Ver puntos, ranking, preview del ranking global, score events
12. Ir a `/leaderboard` - Ver posicion del usuario en el ranking
13. Como ADMIN, recalcular puntuacion de un partido FINISHED
14. Confirmar que los puntos no se duplican
15. Corregir resultado de un partido - Confirmar que los puntos se actualizan

## Pendiente para Fase 10

- Implementacion de salas (rooms)
- Invitaciones temporales con codigo y QR
- Ranking por sala (vista filtrada de miembros ordenados por puntaje global)
- Redis leaderboards:global como cache del ranking global
- Rate limiting con Redis
- Pruebas de estres con k6
