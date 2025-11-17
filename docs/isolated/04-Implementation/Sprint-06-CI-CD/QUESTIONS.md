# Preguntas y Decisiones del Sprint 06

## Q001: ¿GitHub Actions o GitLab CI?
**Decisión por Defecto:** **GitHub Actions**

**Justificación:**
- Repos en GitHub
- Integración nativa
- Actions gratuitas para repos públicos

---

## Q002: ¿Multi-stage Dockerfile?
**Decisión por Defecto:** **Sí (builder + runtime)**

**Justificación:**
- Imagen más pequeña (alpine runtime)
- Build artifacts no en imagen final
- Mejor seguridad

---

## Q003: ¿Publicar a Docker Hub o GitHub Registry?
**Decisión por Defecto:** **GitHub Container Registry (ghcr.io)**

**Justificación:**
- Integrado con GitHub
- Privado gratis
- Mejor para org

---

**Sprint:** 06/06
