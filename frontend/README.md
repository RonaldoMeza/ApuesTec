# ApuesTec — Frontend

Frontend de ApuesTec, plataforma educativa de predicciones deportivas del Mundial.

## Stack

- Next.js 16 + React 19 + TypeScript
- TailwindCSS v4
- shadcn/ui (Button, DropdownMenu, Avatar, Sonner)
- Consume API mediante Nginx en `http://localhost:8081/api/v1`

## Estructura

```
frontend/
  app/
    page.tsx             Landing page pública
    login/page.tsx       Inicio de sesión
    register/page.tsx    Registro de usuarios
    dashboard/page.tsx   Dashboard protegido
    profile/page.tsx     Perfil protegido con cambio de contraseña
    layout.tsx           Layout raíz con AuthProvider y Toaster
    globals.css          Variables CSS modo oscuro
  features/auth/         Módulo de autenticación (feature-based)
    context/             AuthContext + Provider + useAuth hook
    components/          LoginForm, RegisterForm, ProtectedRoute
    services/            auth.service.ts (llamadas HTTP)
    types/               Interfaces TypeScript
    utils/               token-storage.ts (localStorage)
  shared/
    components/          AppLayout, UserNav, LoadingScreen
    services/            api-client.ts (cliente HTTP con manejo de respuestas)
  components/ui/         Componentes shadcn/ui (button, dropdown-menu, avatar, sonner)
```

## Páginas

| Ruta        | Descripción                             | Acceso     |
|-------------|-----------------------------------------|------------|
| `/`         | Landing pública con hero, pasos, reglas | Público    |
| `/login`    | Inicio de sesión                        | Público    |
| `/register` | Registro de usuario                     | Público    |
| `/dashboard`| Panel principal del usuario             | Protegido  |
| `/profile`  | Perfil y cambio de contraseña           | Protegido  |

## Ejecución

### Con Docker (todo incluido)

```bash
docker compose up --build -d
# Frontend en http://localhost:8081
```

### Local + infraestructura en Docker

```powershell
# Terminal 1: infraestructura
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Terminal 2: frontend
cd frontend
npm run dev
# Abrir http://localhost:3000 (Nginx en 8081 también funciona)
```

## Variables de entorno

```env
NEXT_PUBLIC_APP_NAME=ApuesTec
NEXT_PUBLIC_APP_URL=http://localhost:8081
NEXT_PUBLIC_API_URL=http://localhost:8081/api/v1
NEXT_PUBLIC_GOOGLE_CLIENT_ID=...
```

## Diseño

- Modo oscuro por defecto
- Paleta naranja/ámbar como color primario
- Componentes con bordes redondeados, sombras y efectos hover
- Responsive y minimalista
