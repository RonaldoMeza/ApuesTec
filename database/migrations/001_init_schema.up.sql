BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT roles_name_check CHECK (name IN ('SUPER_ADMIN', 'ADMIN', 'USER'))
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    password_hash TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_status_check CHECK (status IN ('ACTIVE', 'BLOCKED', 'DISABLED')),
    CONSTRAINT users_failed_login_attempts_check CHECK (failed_login_attempts >= 0)
);

CREATE TABLE auth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    provider_email VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT auth_accounts_provider_user_unique UNIQUE (provider, provider_user_id),
    CONSTRAINT auth_accounts_user_provider_unique UNIQUE (user_id, provider)
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address INET,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT refresh_tokens_expiration_check CHECK (expires_at > created_at)
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(120) NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ,
    CONSTRAINT rooms_status_check CHECK (status IN ('ACTIVE', 'CLOSED')),
    CONSTRAINT rooms_closed_at_check CHECK (
        (status = 'CLOSED' AND closed_at IS NOT NULL)
        OR (status = 'ACTIVE' AND closed_at IS NULL)
    )
);

CREATE TABLE room_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    room_role VARCHAR(20) NOT NULL DEFAULT 'MEMBER',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT room_members_room_user_unique UNIQUE (room_id, user_id),
    CONSTRAINT room_members_role_check CHECK (room_role IN ('OWNER', 'MODERATOR', 'MEMBER'))
);

CREATE TABLE room_invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    code_hash TEXT NOT NULL UNIQUE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    duration_minutes INTEGER NOT NULL DEFAULT 5,
    expires_at TIMESTAMPTZ NOT NULL,
    used_count INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT room_invites_status_check CHECK (status IN ('ACTIVE', 'EXPIRED', 'REVOKED')),
    CONSTRAINT room_invites_duration_check CHECK (duration_minutes IN (1, 3, 5, 10, 15, 20)),
    CONSTRAINT room_invites_used_count_check CHECK (used_count >= 0),
    CONSTRAINT room_invites_expiration_check CHECK (expires_at > created_at)
);

CREATE UNIQUE INDEX room_invites_one_active_per_room_idx
    ON room_invites (room_id)
    WHERE status = 'ACTIVE';

CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(120) NOT NULL UNIQUE,
    country_code CHAR(3) NOT NULL UNIQUE,
    flag_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT teams_country_code_format_check CHECK (country_code = UPPER(country_code))
);

CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    home_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    away_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    match_date TIMESTAMPTZ NOT NULL,
    home_score INTEGER,
    away_score INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'SCHEDULED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT matches_status_check CHECK (status IN ('SCHEDULED', 'LOCKED', 'FINISHED', 'CANCELLED')),
    CONSTRAINT matches_different_teams_check CHECK (home_team_id <> away_team_id),
    CONSTRAINT matches_home_score_check CHECK (home_score IS NULL OR home_score >= 0),
    CONSTRAINT matches_away_score_check CHECK (away_score IS NULL OR away_score >= 0),
    CONSTRAINT matches_finished_scores_check CHECK (
        (status = 'FINISHED' AND home_score IS NOT NULL AND away_score IS NOT NULL)
        OR status <> 'FINISHED'
    )
);

CREATE TABLE predictions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    predicted_home_score INTEGER NOT NULL,
    predicted_away_score INTEGER NOT NULL,
    is_exact_score BOOLEAN NOT NULL DEFAULT FALSE,
    is_winner_correct BOOLEAN NOT NULL DEFAULT FALSE,
    is_goal_difference_correct BOOLEAN NOT NULL DEFAULT FALSE,
    base_points INTEGER NOT NULL DEFAULT 0,
    early_bonus_points INTEGER NOT NULL DEFAULT 0,
    streak_bonus_points INTEGER NOT NULL DEFAULT 0,
    total_points INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    locked_at TIMESTAMPTZ,
    CONSTRAINT predictions_user_match_unique UNIQUE (user_id, match_id),
    CONSTRAINT predictions_predicted_home_score_check CHECK (predicted_home_score >= 0),
    CONSTRAINT predictions_predicted_away_score_check CHECK (predicted_away_score >= 0),
    CONSTRAINT predictions_points_check CHECK (
        base_points >= 0
        AND early_bonus_points >= 0
        AND streak_bonus_points >= 0
        AND total_points >= 0
    )
);

CREATE TABLE score_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE RESTRICT,
    prediction_id UUID REFERENCES predictions(id) ON DELETE SET NULL,
    points INTEGER NOT NULL,
    reason VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT score_events_points_check CHECK (points <> 0)
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity VARCHAR(100) NOT NULL,
    entity_id UUID,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX auth_accounts_user_id_idx ON auth_accounts(user_id);
CREATE INDEX refresh_tokens_user_id_idx ON refresh_tokens(user_id);
CREATE INDEX user_roles_role_id_idx ON user_roles(role_id);
CREATE INDEX rooms_created_by_idx ON rooms(created_by);
CREATE INDEX room_members_user_id_idx ON room_members(user_id);
CREATE INDEX room_invites_room_id_idx ON room_invites(room_id);
CREATE INDEX matches_home_team_id_idx ON matches(home_team_id);
CREATE INDEX matches_away_team_id_idx ON matches(away_team_id);
CREATE INDEX matches_match_date_idx ON matches(match_date);
CREATE INDEX predictions_match_id_idx ON predictions(match_id);
CREATE INDEX score_events_user_id_idx ON score_events(user_id);
CREATE INDEX score_events_match_id_idx ON score_events(match_id);
CREATE INDEX audit_logs_user_id_idx ON audit_logs(user_id);
CREATE INDEX audit_logs_action_idx ON audit_logs(action);

COMMIT;
