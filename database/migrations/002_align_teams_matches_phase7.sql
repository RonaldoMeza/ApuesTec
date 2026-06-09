-- Migration: Align teams and matches tables with Phase 7 spec
-- Adds code, country to teams; starts_at, locked_at to matches
-- Drops legacy columns country_code and match_date

-- Teams: add code, country, updated_at
ALTER TABLE teams ADD COLUMN IF NOT EXISTS code VARCHAR(10);
ALTER TABLE teams ADD COLUMN IF NOT EXISTS country VARCHAR(100);
ALTER TABLE teams ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Populate code from country_code if code is NULL
UPDATE teams SET code = country_code WHERE code IS NULL;

-- Populate country from name or code if country is NULL
UPDATE teams SET country = name WHERE country IS NULL;

-- Make code NOT NULL and UNIQUE after population
ALTER TABLE teams ALTER COLUMN code SET NOT NULL;
ALTER TABLE teams ADD CONSTRAINT teams_code_unique UNIQUE (code);

-- Drop old country_code column and its constraint
ALTER TABLE teams DROP CONSTRAINT IF EXISTS teams_country_code_format_check;
ALTER TABLE teams DROP CONSTRAINT IF EXISTS teams_country_code_key;
ALTER TABLE teams DROP COLUMN IF EXISTS country_code;

-- Matches: add starts_at, locked_at
ALTER TABLE matches ADD COLUMN IF NOT EXISTS starts_at TIMESTAMPTZ;
ALTER TABLE matches ADD COLUMN IF NOT EXISTS locked_at TIMESTAMPTZ;

-- Populate starts_at from match_date if starts_at is NULL
UPDATE matches SET starts_at = match_date WHERE starts_at IS NULL;

-- Make starts_at NOT NULL after population
ALTER TABLE matches ALTER COLUMN starts_at SET NOT NULL;

-- Drop old match_date column
ALTER TABLE matches DROP COLUMN IF EXISTS match_date;
