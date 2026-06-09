# Fase 6: Frontend Auth Base

## Objetivo

Implementar la base de autenticacion en el frontend: registro, inicio de sesion, proteccion de rutas,
dashboard basico, perfil con cambio de contrasena y cierre de sesion.

## Stack

- Frontend: Next.js 16 + React 19 + TypeScript
- Estilos: TailwindCSS v4 + shadcn/ui
- Backend consumido mediante Nginx en: http://localhost:8081/api/v1

## Arquitectura implementada

Se siguio el patron Feature-Based Architecture con separacion por capaz dentro de `features/auth/`:

```
frontend/
  app/
    page.tsx                     # Home publica (landing page)
    layout.tsx                   # Layout principal con AuthProvider
    login/page.tsx               # Pagina de inicio de sesion
    register/page.tsx            # Pagina de registro
    dashboard/page.tsx           # Dashboard protegido
    profile/page.tsx             # Perfil protegido con cambio de contrasena
    globals.css                  # Variables CSS modo oscuro, colores primarios
  features/auth/
    context/
      AuthContext.tsx            # Contexto de autenticacion + Provider + useAuth
    components/
      LoginForm.tsx              # Formulario de inicio de sesion
      RegisterForm.tsx           # Formulario de registro
      ProtectedRoute.tsx         # Componente de proteccion de rutas
      UserMenu.tsx               # Menu de usuario (nombre + cerrar sesion)
    services/
      auth.service.ts            # Servicio HTTP de autenticacion
    types/
      auth.types.ts              # Tipos TypeScript para auth
    utils/
      token-storage.ts           # Utilidades para localStorage de tokens
  shared/
    components/
      LoadingScreen.tsx          # Pantalla de carga
      ErrorMessage.tsx           # Mensaje de error
      SuccessMessage.tsx         # Mensaje de exito
    services/
      api-client.ts              # Cliente HTTP base con manejo de auth
```

## Endpoints consumidos

| Metodo | Ruta                          | Uso                         |
|--------|-------------------------------|-----------------------------|
| POST   | /api/v1/auth/register         | Registro de usuario         |
| POST   | /api/v1/auth/login            | Inicio de sesion            |
| GET    | /api/v1/auth/me               | Obtener usuario actual      |
| POST   | /api/v1/auth/logout           | Cerrar sesion               |
| POST   | /api/v1/auth/change-password  | Cambiar contrasena          |

## Formato de respuesta del backend

Todas las respuestas del backend usan un envoltorio estandar:

**Exito:**
```json
{ "success": true, "data": { ... }, "message": "..." }
```

**Error:**
```json
{ "success": false, "error": { "code": "...", "message": "..." } }
```

El `api-client.ts` extrae automaticamente `data` de las respuestas exitosas y lanza
errores con `code` y `message` desde la estructura de error.

## Correccion del bug de login

**Problema original:** El frontend asumia que la respuesta del backend era directamente
`{ accessToken, refreshToken, user }`, pero el backend envuelve toda respuesta en
`{ success: true, data: { accessToken, refreshToken, user }, message: "..." }`.

**Solucion:** Se modifico `api-client.ts` para:
1. Parsear el envoltorio `{ success, data, message }` del backend.
2. Si `success === false`, leer `error.code` y `error.message` y lanzar error estructurado.
3. Si `success === true`, retornar `data` directamente, eliminando el envoltorio.
4. Los servicios y componentes ahora trabajan con los datos limpios sin preocuparse del envoltorio.

## Diseno visual

- Modo oscuro por defecto (color-scheme: dark).
- Fondo near-black: `oklch(0.11 0.008 260)`.
- Color primario naranja: `oklch(0.68 0.18 50)` con degradado a amber.
- Bordes redondeados (radius: 0.75rem).
- Cards con sombras suaves y bordes sutiles.
- Inputs con fondo `surface-muted` y borde `border`.
- Botones con gradiente naranja/amber y sombra.
- Efectos hover y animaciones sutiles.
- Diseno responsive.

## Variables de diseno (CSS custom properties)

Definidas en `globals.css` como variables CSS y mapeadas a `@theme` de Tailwind:

```
--color-primary       Color principal naranja
--color-background    Fondo near-black
--color-surface       Color de superficie para cards
--color-surface-hover Color hover de superficie
--color-surface-muted Color muted para inputs/fondos secundarios
--color-border        Color de bordes
--color-muted         Texto muted
```

## Flujo de autenticacion

1. **Registro**: POST /auth/register -> guarda tokens -> redirige a /dashboard.
2. **Login**: POST /auth/login -> guarda tokens -> redirige a /dashboard.
3. **Validacion de sesion**: Al cargar la app, si existe accessToken, se valida via GET /auth/me.
4. **Proteccion de rutas**: /dashboard y /profile estan protegidos por ProtectedRoute.
   Si no hay sesion activa, redirige a /login.
5. **Logout**: POST /auth/logout con refreshToken, limpia tokens locales, redirige a /login.
6. **Cambio de contrasena**: POST /auth/change-password, al exito cierra sesion y redirige a /login.

## Manejo de errores

- Codigos de error del backend: `USER_LOCKED`, `INVALID_CREDENTIALS`, `VALIDATION_ERROR`.
- USER_LOCKED: "Usuario bloqueado por demasiados intentos fallidos."
- INVALID_CREDENTIALS: "Credenciales invalidas."
- Si /me responde 401, se limpia la sesion y redirige a /login.

## Home publica

La pagina principal (`/`) es ahora una landing page publica con:
- Hero section con titulo y descripcion de ApuesTec.
- Botones "Participar ahora" y "Iniciar sesion".
- Seccion "Como funciona" con 4 pasos.
- Seccion "Sistema de puntuacion" con tabla de reglas.
- Preview de proximos partidos.
- Preview de ranking.
- Header con navegacion condicional (autenticado vs invitado).
- Footer con disclaimer educativo.

## Archivos creados

```
frontend/app/login/page.tsx
frontend/app/register/page.tsx
frontend/app/dashboard/page.tsx
frontend/app/profile/page.tsx
frontend/app/page.tsx                       # Landing page publica
frontend/app/globals.css                    # Variables CSS modo oscuro
frontend/features/auth/context/AuthContext.tsx
frontend/features/auth/components/LoginForm.tsx
frontend/features/auth/components/RegisterForm.tsx
frontend/features/auth/components/ProtectedRoute.tsx
frontend/features/auth/components/UserMenu.tsx
frontend/features/auth/services/auth.service.ts
frontend/features/auth/types/auth.types.ts
frontend/features/auth/utils/token-storage.ts
frontend/shared/components/LoadingScreen.tsx
frontend/shared/components/ErrorMessage.tsx
frontend/shared/components/SuccessMessage.tsx
frontend/shared/services/api-client.ts
docs/fases/FASE_06_FRONTEND_AUTH_BASE.md
```

## Archivos modificados

```
frontend/app/layout.tsx              # AuthProvider, metadata, lang=es
frontend/.env.local                  # NEXT_PUBLIC_API_URL
frontend/.env.example                # NEXT_PUBLIC_API_URL
```

## Validacion

```bash
cd frontend
npm run build
npm run lint
docker compose config --quiet
docker compose up --build -d
```

## Criterios de aceptacion

- [x] npm run build pasa
- [x] npm run lint pasa
- [x] docker compose config --quiet pasa
- [x] / es una Home publica y no redirige a login
- [x] /register registra usuario real usando el backend
- [x] /login inicia sesion con usuario real y redirige a /dashboard
- [x] /dashboard esta protegido
- [x] /profile esta protegido
- [x] /me se consume correctamente desde frontend
- [x] logout cierra sesion y revoca refresh token
- [x] change-password funciona desde la interfaz
- [x] Si el usuario esta bloqueado, el frontend muestra mensaje claro
- [x] Si el token es invalido o no existe, redirige a login
- [x] Dashboard mantiene sesion al recargar la pagina
- [x] Diseno en modo oscuro con paleta naranja/amber
- [x] No se implementaron salas, predicciones ni rankings
