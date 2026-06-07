INSERT INTO teams (name, country_code, flag_url)
VALUES
    ('Argentina', 'ARG', NULL),
    ('Brasil', 'BRA', NULL),
    ('Francia', 'FRA', NULL),
    ('Espana', 'ESP', NULL),
    ('Mexico', 'MEX', NULL),
    ('Estados Unidos', 'USA', NULL)
ON CONFLICT (country_code) DO UPDATE
SET name = EXCLUDED.name,
    flag_url = EXCLUDED.flag_url;

INSERT INTO matches (home_team_id, away_team_id, match_date, status)
SELECT home_team.id, away_team.id, fixture.match_date, 'SCHEDULED'
FROM (
    VALUES
        ('ARG', 'MEX', TIMESTAMPTZ '2026-06-11 19:00:00+00'),
        ('BRA', 'USA', TIMESTAMPTZ '2026-06-12 22:00:00+00'),
        ('FRA', 'ESP', TIMESTAMPTZ '2026-06-13 20:00:00+00')
) AS fixture(home_code, away_code, match_date)
JOIN teams AS home_team ON home_team.country_code = fixture.home_code
JOIN teams AS away_team ON away_team.country_code = fixture.away_code
WHERE NOT EXISTS (
    SELECT 1
    FROM matches existing_match
    WHERE existing_match.home_team_id = home_team.id
      AND existing_match.away_team_id = away_team.id
      AND existing_match.match_date = fixture.match_date
);
