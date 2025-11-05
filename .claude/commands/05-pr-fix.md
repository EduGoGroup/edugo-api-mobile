# 05-pr-fix: Revisar y Corregir PR Automáticamente

Analiza un Pull Request, clasifica comentarios de reviewers, verifica pipelines y aplica correcciones automáticas según criterios predefinidos.

## 🎯 Descripción

Este comando invoca al agente especializado **flow-pr-fixer** para:
1. Conectar al PR activo (o al especificado)
2. Obtener comentarios de reviewers (Copilot, Claude, humanos)
3. Verificar estado de pipelines/checks (CI/CD)
4. Clasificar comentarios en categorías:
   - 🟢 **2.1** - Corrección inmediata
   - 🔵 **2.2** - Traducciones/docs (excluir)
   - 🟡 **2.3** - Deuda técnica (resolver después)
   - ⚪ **2.4** - No relevantes
   - 🟣 **2.5** - Dudosos (requieren decisión)
5. Aplicar correcciones automáticas (si --auto-fix)
6. Generar informe estructurado

---

## 📝 Sintaxis

```bash
/05-pr-fix [--pr=NUMBER] [--auto-fix] [--branch=NAME]
```

### Parámetros

| Parámetro | Descripción | Requerido | Default |
|-----------|-------------|-----------|---------|
| `--pr=NUMBER` | Número del PR a revisar | No | PR del branch actual |
| `--auto-fix` | Aplicar correcciones inmediatas sin confirmación | No | `false` |
| `--branch=NAME` | Branch del PR (si no se especifica número) | No | Branch actual |

---

## 📋 Ejemplos de Uso

### Ejemplo 1: Revisión básica del PR actual
```bash
/05-pr-fix
```
**Resultado**: Revisa el PR del branch actual, clasifica comentarios, genera informe y pregunta si aplicar correcciones.

---

### Ejemplo 2: Revisión con auto-corrección
```bash
/05-pr-fix --auto-fix
```
**Resultado**: Revisa el PR, clasifica y **aplica automáticamente** las correcciones inmediatas (Categoría 2.1).

---

### Ejemplo 3: Revisar PR específico
```bash
/05-pr-fix --pr=123
```
**Resultado**: Revisa el PR #123 aunque no sea del branch actual.

---

### Ejemplo 4: Revisión completa con auto-fix de PR específico
```bash
/05-pr-fix --pr=456 --auto-fix
```
**Resultado**: Revisa PR #456 y aplica todas las correcciones inmediatas automáticamente.

---

## 🔄 Flujo de Ejecución

```
1. Parsear parámetros del comando
2. Invocar agente flow-pr-fixer con contexto:
   - PR number (si se especificó)
   - Auto-fix flag
   - Branch name (si se especificó)
3. El agente ejecuta su flujo completo:
   - Conectar al PR
   - Verificar pipelines
   - Obtener comentarios
   - Clasificar comentarios
   - Aplicar correcciones (si auto-fix)
   - Generar informe
4. Mostrar informe al usuario
5. Proponer próximos pasos
```

---

## 🎯 Clasificación de Comentarios

El agente clasifica cada comentario en estas categorías:

### 🟢 Categoría 2.1: Corrección Inmediata
- Errores de sintaxis, typos, formato
- Imports no usados, variables sin usar
- Linting errors menores
- **Acción**: Se corrigen automáticamente (si --auto-fix)

### 🔵 Categoría 2.2: Traducciones/Documentación
- Traducciones de texto
- Mejoras de comentarios o docs
- **Acción**: Se excluyen (fuera de scope)

### 🟡 Categoría 2.3: Deuda Técnica
- Refactorizaciones complejas
- Mejoras de arquitectura
- Optimizaciones no críticas
- **Acción**: Se documentan para resolución posterior

### ⚪ Categoría 2.4: No Relevantes
- Preferencias personales sin impacto
- Ya implementados
- Contradicen guías del proyecto
- **Acción**: Se listan con razón de descarte

### 🟣 Categoría 2.5: Dudosos
- Comentarios ambiguos
- Múltiples opciones de implementación
- Requieren decisión de arquitectura
- **Acción**: Se presentan con opciones al usuario

---

## 📊 Informe Generado

El agente genera un informe markdown estructurado con:

```markdown
# 📊 Informe de Revisión de PR #[NUMBER]

## 🔍 Resumen Ejecutivo
[Cantidad de comentarios por categoría]

## ✅ Estado de Pipelines
[Estado de todos los checks: build, linting, tests]

## 🟢 Categoría 2.1: Correcciones Inmediatas
[Lista de correcciones aplicadas o pendientes]

## 🔵 Categoría 2.2: Excluidos (Traducción/Docs)
[Lista de comentarios excluidos con razón]

## 🟡 Categoría 2.3: Deuda Técnica
[Lista con justificación, impacto, esfuerzo, prioridad]

## ⚪ Categoría 2.4: No Relevantes
[Lista con razón de descarte]

## 🟣 Categoría 2.5: Requieren Decisión
[Lista con opciones: inmediato, deuda técnica, descartar]

## 🎯 Próximos Pasos Sugeridos
[Acciones recomendadas]
```

---

## 🚨 Casos de Uso Especiales

### Si no hay PR activo en el branch
```
❌ Error: No se encontró un PR activo para este branch
💡 Sugerencia: Especifica el número con --pr=NUMBER
```

### Si hay errores de pipeline críticos
```
⚠️ Warning: El PR tiene 3 checks fallidos
El agente intentará identificar y corregir los errores automáticamente
```

### Si auto-fix rompe el build
```
❌ Error: La corrección en [archivo] causó un error de compilación
La corrección fue revertida y movida a Categoría 2.5 (Requiere decisión)
```

---

## 🔧 Requisitos

- **GitHub CLI** (`gh`) instalado y autenticado, O
- **MCP GitHub** configurado en settings
- **Permisos** de lectura/escritura en el repositorio
- **Branch** asociado a un PR abierto (si no se especifica --pr)

---

## 📚 Documentación Relacionada

- Agente: `.claude/agents/flow-pr-fixer.md`
- GitHub CLI: https://cli.github.com/
- Convenciones del proyecto: `.claude/CLAUDE.md`

---

## 💡 Tip: Flujo Recomendado

```bash
# 1. Crear PR y esperar reviews
git push
gh pr create

# 2. Revisar comentarios sin aplicar correcciones
/05-pr-fix

# 3. Revisar informe y decidir sobre comentarios dudosos
[Leer informe generado]

# 4. Aplicar correcciones aprobadas
/05-pr-fix --auto-fix

# 5. Crear documento de deuda técnica (si es necesario)
[Usar informe para crear tech-debt.md]

# 6. Commit y push
git add .
git commit -m "fix: aplicar correcciones de PR review"
git push
```

---

**Versión**: 1.0
**Última actualización**: 2025-11-05
**Responsable**: Claude Code + Jhoan Medina
