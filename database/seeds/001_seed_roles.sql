INSERT INTO roles (name, description)
VALUES
    ('SUPER_ADMIN', 'Rol global con permisos administrativos completos.'),
    ('ADMIN', 'Rol global para administracion operativa del MVP.'),
    ('USER', 'Rol global base para usuarios registrados.')
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description;
