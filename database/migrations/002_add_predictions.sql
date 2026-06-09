-- Migration: Create predictions table (for databases that already had teams/matches)
-- Fase 8: Predicciones de Usuarios

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
        base_points >= 0
        AND early_bonus_points >= 0
        AND streak_bonus_points >= 0
        AND total_points >= 0
    )
);

CREATE INDEX IF NOT EXISTS predictions_match_id_idx ON predictions(match_id);
