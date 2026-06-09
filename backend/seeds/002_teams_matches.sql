-- Seed: Equipos de ejemplo y partidos del Mundial
-- Fase 7: Equipos y Partidos

-- Equipos
INSERT INTO teams (name, code, country) VALUES
    ('Perú', 'PER', 'Perú'),
    ('Argentina', 'ARG', 'Argentina'),
    ('Brasil', 'BRA', 'Brasil'),
    ('Francia', 'FRA', 'Francia'),
    ('Alemania', 'GER', 'Alemania'),
    ('España', 'ESP', 'España'),
    ('Inglaterra', 'ENG', 'Inglaterra'),
    ('Portugal', 'POR', 'Portugal'),
    ('México', 'MEX', 'México'),
    ('Japón', 'JPN', 'Japón')
ON CONFLICT (code) DO NOTHING;

-- Partidos de ejemplo (solo si hay equipos)
WITH team_per AS (SELECT id FROM teams WHERE code = 'PER'),
     team_arg AS (SELECT id FROM teams WHERE code = 'ARG'),
     team_bra AS (SELECT id FROM teams WHERE code = 'BRA'),
     team_fra AS (SELECT id FROM teams WHERE code = 'FRA'),
     team_ger AS (SELECT id FROM teams WHERE code = 'GER'),
     team_esp AS (SELECT id FROM teams WHERE code = 'ESP'),
     team_eng AS (SELECT id FROM teams WHERE code = 'ENG'),
     team_por AS (SELECT id FROM teams WHERE code = 'POR'),
     team_mex AS (SELECT id FROM teams WHERE code = 'MEX'),
     team_jpn AS (SELECT id FROM teams WHERE code = 'JPN')
INSERT INTO matches (home_team_id, away_team_id, starts_at, status)
SELECT t1.id, t2.id, t1.starts_at, 'SCHEDULED'
FROM (VALUES
    ((SELECT id FROM team_arg), (SELECT id FROM team_bra), NOW() + INTERVAL '3 days'),
    ((SELECT id FROM team_esp), (SELECT id FROM team_fra), NOW() + INTERVAL '5 days'),
    ((SELECT id FROM team_ger), (SELECT id FROM team_eng), NOW() + INTERVAL '7 days'),
    ((SELECT id FROM team_por), (SELECT id FROM team_per), NOW() + INTERVAL '10 days'),
    ((SELECT id FROM team_mex), (SELECT id FROM team_jpn), NOW() + INTERVAL '14 days')
) AS t1(id, id2, starts_at)
WHERE EXISTS (SELECT 1 FROM teams WHERE code IN ('PER','ARG','BRA','FRA','GER','ESP','ENG','POR','MEX','JPN'));
