# Fase 5: Limpieza de Código Legacy

> **Estado:** ⏳ PENDIENTE
> **PR:** -
> **Branch:** `feature/legacy-cleanup` (por crear)
> **Estimado:** 2-4 horas

---

## Tareas

### ⏳ DEP-002: Limpiar repositorio legacy de Assessments
- **Problema:** Existen referencias a repositorios legacy de assessments
- **Solución:** Verificar que no hay dependencias y eliminar/marcar como deprecated
- **Commit:** Pendiente

**Pasos:**
1. Buscar referencias a repositorios legacy
2. Verificar que código nuevo no los usa
3. Eliminar o marcar claramente como `// DEPRECATED`
4. Actualizar documentación

---

### ⏳ DEBT-004: Documentar plan de consolidación de sistemas Assessment
- **Tipo:** Documentación
- **Problema:** Existen múltiples implementaciones de assessment
- **Solución:** Crear documento de migración con timeline
- **Commit:** Pendiente

**Entregable:**
- `documents/improvements/ASSESSMENT_CONSOLIDATION_PLAN.md`
- Timeline para eliminación de código legacy
- Dependencias entre sistemas
- Riesgos y mitigaciones

---

### ⏳ Eliminar código comentado restante
- **Problema:** Bloques de código comentado en varios archivos
- **Solución:** Buscar y eliminar o crear issues para funcionalidad faltante
- **Commit:** Pendiente

**Pasos:**
1. Buscar patrones de código comentado extenso
2. Evaluar si es código obsoleto o funcionalidad pendiente
3. Si obsoleto → eliminar
4. Si pendiente → crear issue y eliminar comentario
