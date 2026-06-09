-- Seed: Equipos de ejemplo y partidos del Mundial
-- Fase 7: Equipos y Partidos

-- Equipos
INSERT INTO teams (name, code, country, flag_url) VALUES
    ('Perú', 'PER', 'Perú', NULL),
    ('Argentina', 'ARG', 'Argentina', NULL),
    ('Brasil', 'BRA', 'Brasil', NULL),
    ('Francia', 'FRA', 'Francia', NULL),
    ('Alemania', 'GER', 'Alemania', NULL),
    ('España', 'ESP', 'España', NULL),
    ('Inglaterra', 'ENG', 'Inglaterra', NULL),
    ('Portugal', 'POR', 'Portugal', NULL),
    ('México', 'MEX', 'México', NULL),
    ('Japón', 'JPN', 'Japón', NULL)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
    country = EXCLUDED.country,
    flag_url = EXCLUDED.flag_url;

-- Partidos de ejemplo
INSERT INTO matches (home_team_id, away_team_id, starts_at, status)
SELECT home.id, away.id, fixture.starts_at, 'SCHEDULED'
FROM (VALUES
    ('ARG', 'BRA', NOW() + INTERVAL '3 days'),
    ('ESP', 'FRA', NOW() + INTERVAL '5 days'),
    ('GER', 'ENG', NOW() + INTERVAL '7 days'),
    ('POR', 'PER', NOW() + INTERVAL '10 days'),
    ('MEX', 'JPN', NOW() + INTERVAL '14 days')
) AS fixture(home_code, away_code, starts_at)
JOIN teams AS home ON home.code = fixture.home_code
JOIN teams AS away ON away.code = fixture.away_code
WHERE NOT EXISTS (
    SELECT 1 FROM matches existing
    WHERE existing.home_team_id = home.id
      AND existing.away_team_id = away.id
      AND existing.starts_at = fixture.starts_at
);
