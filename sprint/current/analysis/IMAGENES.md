# 🖼️ Guía de Imágenes de Diagramas

**Generado**: 2025-11-04  
**Total de imágenes**: 10 diagramas PNG

---

## 📐 Diagramas de Arquitectura

### architecture-1.png (42.7 KB)
**Fuente**: `architecture-phase-2.md`  
**Descripción**: Diagrama de arquitectura general de la Fase 2  
**Componentes mostrados**:
- HTTP Layer (Material Handler, Assessment Handler)
- Application Layer (Services)
- Infrastructure Layer (RabbitMQ Publisher, S3 Client, PostgreSQL, MongoDB)
- External Systems (RabbitMQ Server, AWS S3)

### architecture-2.png (35.2 KB)
**Fuente**: `architecture-phase-2.md`  
**Descripción**: Diagrama de flujo de eventos  
**Muestra**: Secuencia de creación de material con S3 upload y publicación de eventos

---

## 🗄️ Diagramas de Modelo de Datos

### data-model-1.png (70.1 KB)
**Fuente**: `data-model-phase-2.md`  
**Descripción**: Diagrama ER (Entidad-Relación) de PostgreSQL  
**Entidades mostradas**:
- MATERIALS ↔ MATERIAL_VERSIONS
- USERS ↔ USER_PROGRESS
- USERS ↔ REFRESH_TOKENS
- USERS ↔ LOGIN_ATTEMPTS

---

## 🔄 Diagramas de Procesos

### process-1.png (129.2 KB)
**Fuente**: `process-diagram-phase-2.md`  
**Descripción**: **Flujo principal de creación de material con S3 Upload**  
**Pasos mostrados**:
1. POST /materials con metadata
2. Validación de DTO
3. Generación de presigned URL (S3)
4. Insert metadata en PostgreSQL
5. Publicación de evento material_uploaded (RabbitMQ)
6. Retorno de response con presignedURL
7. Upload directo del cliente a S3

### process-2.png (138.0 KB)
**Fuente**: `process-diagram-phase-2.md`  
**Descripción**: **Flujo de registro de intento de assessment con evento**  
**Pasos mostrados**:
1. POST /assessments/:id/attempts
2. Validación de respuestas
3. Recuperación de assessment de MongoDB
4. Cálculo de puntaje
5. Persistencia de intento
6. Generación de feedback (aggregation)
7. Publicación de evento assessment_attempt_recorded

### process-3.png (116.2 KB)
**Fuente**: `process-diagram-phase-2.md`  
**Descripción**: **Flujo de actualización de progreso (UPSERT)**  
**Pasos mostrados**:
1. PUT /progress
2. Validación de porcentaje
3. UPSERT en PostgreSQL (INSERT ... ON CONFLICT ... DO UPDATE)
4. Lógica condicional: GREATEST para actualizar solo si nuevo % > actual
5. Cálculo automático de status
6. Establecimiento de completed_at

### process-4.png (73.4 KB)
**Fuente**: `process-diagram-phase-2.md`  
**Descripción**: **Flujo de consulta de estadísticas (Aggregation)**  
**Pasos mostrados**:
1. GET /users/:id/stats
2. Aggregation pipeline MongoDB (5 stages)
3. Cálculos: total_attempts, average_score, highest/lowest
4. Construcción de array de recent_attempts
5. Retorno de estadísticas completas

### process-5.png (38.5 KB)
**Fuente**: `process-diagram-phase-2.md`  
**Descripción**: **Flujo alternativo - Fallo de RabbitMQ**  
**Muestra**: Manejo de error cuando RabbitMQ no está disponible (evento no crítico, log warning pero continuar)

### process-6.png (45.5 KB)
**Fuente**: `process-diagram-phase-2.md`  
**Descripción**: **Flujo alternativo - Presigned URL expirada**  
**Muestra**: Proceso cuando cliente intenta usar URL expirada y necesita solicitar nueva URL

### process-7.png (38.6 KB)
**Fuente**: `process-diagram-phase-2.md`  
**Descripción**: **Flujo alternativo - Query compleja lenta (Timeout)**  
**Muestra**: Manejo de timeout en queries complejas con logging y error 503

---

## 📊 Resumen por Categoría

| Categoría | Cantidad | Tamaño Total |
|-----------|----------|--------------|
| Arquitectura | 2 | 77.9 KB |
| Modelo de Datos | 1 | 70.1 KB |
| Procesos | 7 | 579.4 KB |
| **TOTAL** | **10** | **727.4 KB** |

---

## 🔧 Cómo se Generaron

Estas imágenes fueron generadas automáticamente desde los archivos Markdown usando:

**Script**: `generate-diagrams.py`  
**Herramienta**: `mermaid-cli` (mmdc)  
**Formato**: PNG con fondo transparente y tema dark  

### Regenerar Imágenes

Si necesitas regenerar las imágenes (por ejemplo, después de modificar los diagramas en los archivos Markdown):

```bash
cd sprint/current/analysis
python3 generate-diagrams.py
```

---

## 📖 Archivos de Origen

- `architecture-phase-2.md` - Arquitectura del sistema
- `data-model-phase-2.md` - Modelo de datos y queries
- `process-diagram-phase-2.md` - Flujos de procesos
- `readme-phase-2.md` - Resumen ejecutivo

---

**Nota**: Las imágenes PNG son más pesadas que los diagramas Mermaid en texto, pero son más fáciles de visualizar en editores que no soportan Mermaid nativamente (como Zed).
