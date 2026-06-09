-- Migration: Create score_events (if not exists) and user_scores
-- Fase 9: Puntuacion y Ranking Global

-- Handle score_events: create if not exist, or add columns if old schema exists
CREATE TABLE IF NOT EXISTS score_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    prediction_id UUID NOT NULL REFERENCES predictions(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN ('EXACT_SCORE', 'WINNER_CORRECT', 'GOAL_DIFFERENCE_CORRECT', 'EARLY_BONUS', 'STREAK_BONUS')),
    points INTEGER NOT NULL CHECK (points > 0),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- If the table already existed with old schema, add the missing columns
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'score_events' AND column_name = 'event_type'
    ) THEN
        ALTER TABLE score_events ADD COLUMN event_type VARCHAR(50);
        ALTER TABLE score_events ADD COLUMN description TEXT;
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'score_events' AND column_name = 'reason'
        ) THEN
            UPDATE score_events SET event_type = reason WHERE event_type IS NULL;
            ALTER TABLE score_events DROP COLUMN reason;
        END IF;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS user_scores (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_points INTEGER NOT NULL DEFAULT 0,
    predictions_count INTEGER NOT NULL DEFAULT 0,
    exact_scores_count INTEGER NOT NULL DEFAULT 0,
    winner_correct_count INTEGER NOT NULL DEFAULT 0,
    goal_difference_correct_count INTEGER NOT NULL DEFAULT 0,
    streak_count INTEGER NOT NULL DEFAULT 0,
    last_scored_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_score_events_user_id ON score_events(user_id);
CREATE INDEX IF NOT EXISTS idx_score_events_match_id ON score_events(match_id);
CREATE INDEX IF NOT EXISTS idx_score_events_prediction_id ON score_events(prediction_id);
CREATE INDEX IF NOT EXISTS idx_user_scores_total_points ON user_scores(total_points DESC);
