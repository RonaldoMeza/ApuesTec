# Checklist de Pruebas - Fase 5: Autenticacion y Usuarios

> Marca cada item como `[x]` cuando lo hayas verificado.

---

## 1. Compilacion e infraestructura

- [ ] `go test ./...` pasa sin errores (desde `backend/`)
- [ ] `go build ./...` compila correctamente
- [ ] `docker compose config --quiet` pasa sin errores (desde raiz del proyecto)
- [ ] `docker compose up --build -d` levanta todos los servicios sin errores
- [ ] `GET /api/v1/health` responde `200` con `{"status":"ok"}`
- [ ] `GET /api/v1/health/dependencies` responde `200` con postgres y redis ok

---

## 2. Registro de usuario

- [ ] `POST /api/v1/auth/register` con datos validos devuelve `201` y `accessToken`, `refreshToken`, `user`
- [ ] El `user` contiene: `id`, `fullName`, `email`, `roles: ["USER"]`
- [ ] `POST /api/v1/auth/register` con email repetido devuelve `409 CONFLICT` con `EMAIL_EXISTS`
- [ ] `POST /api/v1/auth/register` con password < 8 caracteres devuelve `400 VALIDATION_ERROR`
- [ ] `POST /api/v1/auth/register` con email invalido devuelve `400 VALIDATION_ERROR`
- [ ] En `audit_logs` existe un registro `USER_REGISTERED` para el nuevo usuario

---

## 3. Login

- [ ] `POST /api/v1/auth/login` con credenciales correctas devuelve `200` con `accessToken`, `refreshToken`, `user`
- [ ] El `user` contiene los datos correctos del usuario registrado
- [ ] `POST /api/v1/auth/login` con email incorrecto devuelve `401 INVALID_CREDENTIALS`
- [ ] `POST /api/v1/auth/login` con password incorrecto devuelve `401 INVALID_CREDENTIALS`
- [ ] En `audit_logs` existe un registro `USER_LOGGED_IN` para login exitoso
- [ ] En `audit_logs` existe un registro `LOGIN_FAILED` para login fallido

---

## 4. Bloqueo por intentos fallidos

- [ ] Enviar 5 logins con password incorrecto para el mismo email
- [ ] El 5to intento devuelve `429 USER_LOCKED` con mensaje de bloqueo
- [ ] En `audit_logs` existe un registro `USER_LOCKED`
- [ ] Intentar login inmediatamente despues del bloqueo devuelve `429 USER_LOCKED`
- [ ] (Opcional) Esperar 15 minutos (o reducir `LOGIN_LOCK_MINUTES=1` en `.env`) y verificar que el login vuelve a funcionar

---

## 5. Refresh token

- [ ] `POST /api/v1/auth/refresh` con refresh token valido devuelve `200` con NUEVO `accessToken` y NUEVO `refreshToken`
- [ ] El refresh token anterior queda revocado (no se puede reusar)
- [ ] `POST /api/v1/auth/refresh` con token invalido devuelve `401 INVALID_REFRESH_TOKEN`
- [ ] `POST /api/v1/auth/refresh` con token ya usado (replay) devuelve `401`
- [ ] En `audit_logs` existe un registro `TOKEN_REFRESHED`

---

## 6. Logout

- [ ] `POST /api/v1/auth/logout` con refresh token valido devuelve `200`
- [ ] Ese refresh token ya no puede usarse para refresh (devuelve `401`)
- [ ] En `audit_logs` existe un registro `USER_LOGGED_OUT`

---

## 7. Endpoint /me

- [ ] `GET /api/v1/auth/me` SIN token devuelve `401 UNAUTHORIZED`
- [ ] `GET /api/v1/auth/me` CON token valido devuelve `200` con datos del usuario
- [ ] `GET /api/v1/auth/me` CON token invalido/expirado devuelve `401 INVALID_TOKEN`
- [ ] Los datos incluyen: `id`, `fullName`, `email`, `roles`

---

## 8. Google Auth

- [ ] `POST /api/v1/auth/google` con `idToken` valido devuelve `200` con tokens
- [ ] `POST /api/v1/auth/google` con `idToken` invalido devuelve `401 INVALID_ID_TOKEN`
- [ ] En `audit_logs` existe `GOOGLE_AUTH_SUCCESS` o `GOOGLE_AUTH_FAILED` segun el caso

---

## 9. Cambio de password

- [ ] `POST /api/v1/auth/change-password` SIN token devuelve `401`
- [ ] `POST /api/v1/auth/change-password` con password actual incorrecta devuelve `401 WRONG_PASSWORD`
- [ ] `POST /api/v1/auth/change-password` con password nueva < 8 caracteres devuelve `400 WEAK_PASSWORD`
- [ ] `POST /api/v1/auth/change-password` con datos correctos devuelve `200`
- [ ] Despues del cambio, los refresh tokens anteriores quedan revocados (no se puede refresh)
- [ ] Se puede hacer login con la nueva password
- [ ] En `audit_logs` existe un registro `PASSWORD_CHANGED`

---

## 10. Seguridad

- [ ] El refresh token NUNCA se almacena en texto plano en PostgreSQL (solo SHA-256 hash)
- [ ] `refresh_tokens.token_hash` contiene solo el hash, no el token original
- [ ] El backend genera sus propios JWT (no usa Firebase ni servicios externos)
- [ ] El access token expira en 1 hora (o el valor configurado)
- [ ] JWT firmado con HMAC-SHA256 usando `JWT_ACCESS_SECRET`
- [ ] No hay secretos hardcodeados, todo via variables de entorno

---

## 11. Roles (si aplica endpoint protegido)

- [ ] El rol `USER` se asigna automaticamente al registrar
- [ ] El middleware `RequireRole` funciona correctamente (probado con SUPER_ADMIN, ADMIN, USER)

---

## Resumen

| Seccion                    | Total | Pasaron |
|----------------------------|-------|---------|
| Compilacion e infraestr.   | 6     | ___/6   |
| Registro                   | 6     | ___/6   |
| Login                      | 6     | ___/6   |
| Bloqueo intentos fallidos  | 5     | ___/5   |
| Refresh token              | 5     | ___/5   |
| Logout                     | 3     | ___/3   |
| /me                        | 4     | ___/4   |
| Google Auth                | 3     | ___/3   |
| Cambio password            | 6     | ___/6   |
| Seguridad                  | 6     | ___/6   |
| Roles                      | 2     | ___/2   |
| **TOTAL**                  | **52**| **___/52** |
