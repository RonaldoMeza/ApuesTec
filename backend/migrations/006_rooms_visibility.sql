ALTER TABLE rooms ADD COLUMN IF NOT EXISTS visibility VARCHAR(7) NOT NULL DEFAULT 'PRIVATE';
ALTER TABLE rooms ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
ALTER TABLE rooms ADD COLUMN IF NOT EXISTS network_prefix VARCHAR(15);

CREATE INDEX IF NOT EXISTS idx_rooms_public_network ON rooms(visibility, network_prefix) WHERE visibility = 'PUBLIC';
