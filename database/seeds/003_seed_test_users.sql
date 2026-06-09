-- Seed: Usuarios de prueba para desarrollo
-- Crea un usuario por cada rol (USER, ADMIN, SUPER_ADMIN)
-- Contraseña para todos: test1234

-- Inserción de usuarios
INSERT INTO users (full_name, email, password_hash, status)
SELECT full_name, email, crypt, 'ACTIVE'
FROM (VALUES
    ('Usuario Regular', 'user@apuestec.dev', crypt('test1234', gen_salt('bf', 8))),
    ('Admin ApuesTec', 'admin@apuestec.dev', crypt('test1234', gen_salt('bf', 8))),
    ('Super Admin', 'superadmin@apuestec.dev', crypt('test1234', gen_salt('bf', 8)))
) AS u(full_name, email, crypt)
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE users.email = u.email
);

-- Asignación de roles
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email = 'user@apuestec.dev' AND r.name = 'USER'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email = 'admin@apuestec.dev' AND r.name = 'ADMIN'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email = 'superadmin@apuestec.dev' AND r.name = 'SUPER_ADMIN'
ON CONFLICT DO NOTHING;
