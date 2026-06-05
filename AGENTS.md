# AGENTS.md - Instrucciones para agentes IA del proyecto ApuesTec

## 1. Contexto general

Este repositorio contiene el proyecto **ApuesTec**, una plataforma educativa de predicciones deportivas del Mundial.

El sistema NO maneja dinero real, NO realiza apuestas monetarias, NO usa cuotas reales y NO procesa pagos.  
La aplicación funciona mediante predicciones, puntos, rankings, salas privadas y gamificación.

El objetivo principal es desarrollar una aplicación web completa, segura, escalable y documentada, cumpliendo con los entregables académicos:

- Propuesta de arquitectura del sistema.
- Desarrollo de la aplicación.
- Contenerización con Docker.
- Escalamiento horizontal y balanceo de carga.
- Pruebas de estrés con k6.
- Sistema de puntuación.
- Sistema de salas/grupos con invitación.
- Roles, permisos, seguridad y auditoría.

---

## 2. Regla principal para cualquier agente

Antes de crear, modificar o eliminar código, el agente debe leer y respetar el archivo:

```text
PLAN_PROYECTO_APUESTEC.md

Ese archivo contiene las decisiones técnicas, funcionales y arquitectónicas del sistema.

No se debe cambiar el stack tecnológico, la arquitectura base ni las reglas de negocio sin justificarlo claramente.