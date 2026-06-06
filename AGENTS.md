# AGENTS.md - Instrucciones para agentes IA del proyecto ApuesTec

## 1. Contexto general

Este repositorio contiene el proyecto **ApuesTec**, una plataforma educativa de predicciones deportivas del Mundial.

El sistema NO maneja dinero real, NO realiza apuestas monetarias, NO usa cuotas reales, NO procesa pagos y NO se integra con casas de apuestas. La aplicacion funciona mediante predicciones de marcadores, puntos, rankings, salas privadas, invitaciones temporales, roles, permisos y gamificacion.

## 2. Stack obligatorio

Frontend:

- Next.js.
- React.
- TypeScript.
- TailwindCSS.
- shadcn/ui.
- Recharts.

Backend:

- Go.
- Gin Gonic.

Persistencia e infraestructura:

- PostgreSQL como base de datos principal.
- Redis para rankings, cache, rate limiting e invitaciones temporales.
- Docker y Docker Compose.
- Nginx como reverse proxy y balanceador de carga.
- k6 para pruebas de estres.

## 3. Reglas obligatorias del proyecto

- Antes de crear, modificar o eliminar codigo, leer y respetar `PLAN_PROYECTO_APUESTEC.md`.
- No cambiar el stack tecnologico sin autorizacion explicita.
- No usar Firebase Auth.
- No usar API externa de partidos en el MVP.
- No usar dinero real, pagos, cuotas reales, casas de apuestas ni logica de apuestas monetarias.
- No crear codigo espagueti.
- No mezclar responsabilidades.
- No subir secretos ni archivos `.env` reales.
- El backend siempre debe emitir sus propios JWT.
- Google OAuth solo valida identidad.
- El refresh token no debe guardarse en `localStorage`.
- Antes de avanzar de fase, validar los criterios de aceptacion definidos en `PLAN_PROYECTO_APUESTEC.md`.

## 4. Arquitectura obligatoria

Backend:

- Usar Feature-Based Architecture.
- Usar Service Layer.
- Usar Repository Pattern.
- Usar DTO Pattern.
- Usar Middleware Pattern.
- No colocar logica de negocio en handlers.
- No acceder directamente a PostgreSQL desde handlers.
- No acceder directamente a Redis desde handlers.
- No duplicar logica de autorizacion.

Frontend:

- Usar Feature-Based Architecture.
- Usar Next.js App Router.
- Centralizar llamadas HTTP en services por feature.
- Proteger rutas segun autenticacion y rol.
- Mantener componentes pequenos, declarativos y reutilizables.

## 5. Reglas funcionales criticas

- Las salas del MVP solo usan estados `ACTIVE` y `CLOSED`.
- Las salas no se eliminan fisicamente en el MVP.
- Una sala solo puede tener una invitacion activa al mismo tiempo.
- Una nueva invitacion de sala revoca automaticamente la invitacion activa anterior de esa misma sala.
- Las invitaciones solo aceptan duraciones de 1, 3, 5, 10, 15 y 20 minutos.
- Los codigos de invitacion se guardan como hash en PostgreSQL.
- Redis debe usar claves basadas en hash del codigo y aplicar TTL.
- ApuesTec maneja un unico puntaje oficial global por usuario.
- Las salas no generan puntajes independientes ni billeteras separadas.
- El ranking de sala es una vista filtrada de miembros de la sala ordenados por su puntaje global acumulado.
- Los puntos no se duplican por pertenecer a varias salas.
- PostgreSQL es la fuente oficial de verdad.
- Redis es una capa de cache/optimizacion y los rankings deben poder reconstruirse desde PostgreSQL.

## 6. Modo de trabajo del agente

1. Antes de modificar archivos, revisar el objetivo de la fase actual.
2. No implementar modulos fuera de orden salvo autorizacion.
3. No instalar dependencias innecesarias.
4. No crear pantallas ni endpoints que no esten definidos en el contrato minimo.
5. Si existe una contradiccion entre `AGENTS.md` y `PLAN_PROYECTO_APUESTEC.md`, detenerse y reportarla antes de continuar.
6. Toda implementacion debe respetar separacion de responsabilidades.
7. Cada cambio debe mantener la documentacion actualizada.
