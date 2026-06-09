# Fase 8: Predicciones de Usuarios

## Objetivo

Implementar predicciones de usuarios sobre partidos. Un usuario autenticado puede registrar, editar y ver sus predicciones de marcador para cada partido.

## Archivos creados/modificados

### Backend (creados)
- `backend/internal/predictions/model.go` — Modelo `Prediction`
- `backend/internal/predictions/dto.go` — DTOs de request/response + errores
- `backend/internal/predictions/repository.go` — Acceso a base de datos
- `backend/internal/predictions/service.go` — Lógica de negocio
- `backend/internal/predictions/handler.go` — Handlers HTTP
- `backend/migrations/003_predictions.sql` — Migración de tabla predictions
- `database/migrations/002_add_predictions.sql` — Migración para bases existentes

### Backend (modificados)
- `backend/internal/audit/repository.go` — Constantes `PREDICTION_CREATED`, `PREDICTION_UPDATED`
- `backend/internal/routes/router.go` — Registro de rutas y dependencias

### Frontend (creados)
- `frontend/features/predictions/types/prediction.types.ts` — Tipos TypeScript
- `frontend/features/predictions/services/prediction.service.ts` — Servicio HTTP
- `frontend/features/predictions/components/PredictionForm.tsx` — Formulario de predicción
- `frontend/features/predictions/components/PredictionCard.tsx` — Card de predicción
- `frontend/features/predictions/components/PredictionStatus.tsx` — Badge de estado
- `frontend/app/predictions/page.tsx` — Página de historial de predicciones

### Frontend (modificados)
- `frontend/app/matches/[id]/page.tsx` — Integración de formulario de predicción
- `frontend/app/dashboard/page.tsx` — Sección "Mis predicciones recientes"
- `frontend/shared/components/AppLayout.tsx` — Enlace "Mis predicciones" en navbar
- `frontend/shared/components/UserNav.tsx` — Ítem "Mis predicciones" en menú

## Tabla predictions

```sql
CREATE TABLE IF NOT EXISTS predictions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    predicted_home_score INTEGER NOT NULL CHECK (predicted_home_score >= 0),
    predicted_away_score INTEGER NOT NULL CHECK (predicted_away_score >= 0),
    is_exact_score BOOLEAN NOT NULL DEFAULT FALSE,
    is_winner_correct BOOLEAN NOT NULL DEFAULT FALSE,
    is_goal_difference_correct BOOLEAN NOT NULL DEFAULT FALSE,
    base_points INTEGER NOT NULL DEFAULT 0,
    early_bonus_points INTEGER NOT NULL DEFAULT 0,
    streak_bonus_points INTEGER NOT NULL DEFAULT 0,
    total_points INTEGER NOT NULL DEFAULT 0,
    locked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT predictions_user_match_unique UNIQUE (user_id, match_id),
    CONSTRAINT predictions_points_check CHECK (
        base_points >= 0 AND early_bonus_points >= 0
        AND streak_bonus_points >= 0 AND total_points >= 0
    )
);
```

Restricciones:
- `UNIQUE(user_id, match_id)` — Una predicción por usuario y partido
- `CHECK (predicted_home_score >= 0)` — Sin marcadores negativos
- `CHECK (predicted_away_score >= 0)` — Sin marcadores negativos
- `locked_at` not NULL indica predicción bloqueada

## Endpoints implementados

Protegidos con JWT:

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/v1/predictions/me` | Lista todas las predicciones del usuario autenticado |
| GET | `/api/v1/matches/:id/prediction` | Obtiene la predicción del usuario para un partido |
| POST | `/api/v1/matches/:id/prediction` | Crea o actualiza predicción |
| DELETE | `/api/v1/matches/:id/prediction` | Elimina predicción si es editable |

### DTO de request (POST)
```json
{
  "homeScorePredicted": 2,
  "awayScorePredicted": 1
}
```

### DTO de response
```json
{
  "id": "uuid",
  "matchId": "uuid",
  "homeScorePredicted": 2,
  "awayScorePredicted": 1,
  "points": 0,
  "isLocked": false,
  "canEdit": true,
  "createdAt": "2026-06-08T22:00:00Z",
  "updatedAt": "2026-06-08T22:00:00Z"
}
```

## Reglas de negocio

1. Solo usuarios autenticados pueden crear predicciones.
2. Un usuario solo puede registrar una predicción por partido (UNIQUE).
3. Si ya existe predicción, se actualiza si está permitido editar.
4. No permitir predicciones en partidos `FINISHED` o `CANCELLED`.
5. No permitir predicciones en partidos `LOCKED`.
6. No permitir predicciones si faltan 10 minutos o menos para `starts_at`.
7. No permitir marcadores negativos.
8. No permitir edición si la predicción ya está bloqueada (`locked_at IS NOT NULL`).
9. El `userID` se obtiene del JWT, no del frontend.
10. `total_points` se deja en 0 por ahora (Fase 9).
11. Se registra auditoría: `PREDICTION_CREATED`, `PREDICTION_UPDATED`.

## Pruebas realizadas

- `go build ./...` — Sin errores
- `go test ./...` — Tests pasan
- `npm run lint` — Sin errores
- `npm run build` — Build exitoso
- `docker compose config --quiet` — Configuración válida

## Cómo probar desde frontend

1. Entrar a `/matches` sin login.
2. Abrir detalle de partido sin login y ver CTA "Inicia sesión para registrar tu predicción".
3. Iniciar sesión.
4. Abrir detalle de partido `SCHEDULED`.
5. Registrar predicción (marcador local y visitante).
6. Ver mensaje de éxito.
7. Recargar la página y confirmar que la predicción sigue visible.
8. Editar predicción.
9. Ir a `/dashboard` y ver predicción reciente en la sección "Mis predicciones recientes".
10. Ir a `/predictions` y ver historial completo.
11. Intentar predecir partido `FINISHED`/`CANCELLED`/`LOCKED` y ver bloqueo visual.
12. Verificar que no se puede crear más de una predicción duplicada.
13. Verificar que `USER` no pueda acceder a rutas admin.
14. Verificar que `ADMIN`/`SUPER_ADMIN` siga pudiendo gestionar partidos/equipos.

## Qué queda pendiente para Fase 9

- Cálculo real de puntos (`points`) cuando un partido finaliza.
- Ranking global completo.
- Rankins por salas.
- Implementación del módulo de salas (`rooms`).
