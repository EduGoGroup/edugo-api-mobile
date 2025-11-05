# Análisis del Sprint - Optimización de Queries: Índice en Materials

## Resumen Ejecutivo

Este sprint se enfoca en una tarea de optimización de base de datos atómica y autocontenida: la creación de un índice descendente en la columna `updated_at` de la tabla `materials` en PostgreSQL. Aunque es una tarea técnicamente simple, representa una mejora tangible de performance en un caso de uso real y frecuente del sistema: el listado de materiales educativos ordenados cronológicamente.

La estrategia del sprint es implementar una solución quirúrgica que no introduce riesgos ni dependencias, permitiendo validar simultáneamente el sistema de automatización de comandos/agentes mientras se entrega valor funcional al proyecto. El alcance deliberadamente limitado asegura que la implementación es verificable, reversible y puede completarse en una sola sesión de trabajo sin efectos colaterales.

Este sprint forma parte de la Fase 3 del plan maestro del proyecto (Implementar Queries Complejas), específicamente la subtarea 3.2, y fue seleccionado estratégicamente como caso de prueba por su bajo riesgo y alta claridad de criterios de éxito.

## Objetivo del Sprint

**Objetivo Principal**: Crear un script SQL de migración que agregue un índice descendente en la columna `updated_at` de la tabla `materials` para optimizar consultas de listado ordenadas por fecha de actualización más reciente.

**Objetivos Secundarios**:
- Validar el sistema de comandos/agentes del proyecto con una tarea real
- Establecer precedente de optimizaciones incrementales de base de datos
- Mejorar la experiencia de usuario en listados de materiales recientes
- Documentar proceso de validación de performance con EXPLAIN ANALYZE

**Criterio de Éxito**: El script SQL ejecutado exitosamente debe resultar en que consultas con `ORDER BY updated_at DESC` utilicen el índice en su plan de ejecución, demostrando mejora de performance medible.

## Arquitectura Propuesta

### Tipo de Arquitectura
Este sprint no modifica la arquitectura del sistema, sino que implementa una **optimización de capa de persistencia** dentro del patrón de **Clean Architecture (Hexagonal)** existente del proyecto EduGo API Mobile.

### Descripción de Arquitectura

El proyecto mantiene su estructura de tres capas principales:

1. **Capa de Dominio** (`internal/domain/`): Define la entidad `Material` con sus propiedades, incluyendo el campo `updated_at` de tipo timestamp. Esta capa no se modifica.

2. **Capa de Aplicación** (`internal/application/`): Contiene los servicios que invocan métodos del repositorio de Materials, como `ListMaterials()` con filtros y ordenamiento. Esta capa tampoco se modifica.

3. **Capa de Infraestructura** (`internal/infrastructure/`):
   - **Persistencia** (`persistence/postgresql/`): Aquí reside el repositorio `MaterialRepository` que ejecuta queries SQL contra PostgreSQL. Esta capa se beneficia indirectamente del índice sin modificación de código.
   - **HTTP Handlers** (`http/handler/`): Los handlers que exponen endpoints REST para listar materiales experimentarán mejoras de latencia sin cambios en su implementación.

### Componentes Principales

**1. Base de Datos PostgreSQL** (Componente modificado)
- **Responsabilidad**: Almacenamiento persistente de materiales educativos
- **Tecnología**: PostgreSQL 14+
- **Modificación**: Agregado de índice `idx_materials_updated_at` en la tabla `materials`
- **Interacciones**: Es consultado por el `MaterialRepository` mediante queries SQL

**2. Script de Migración** (Componente nuevo)
- **Responsabilidad**: Aplicar cambio de esquema de forma idempotente
- **Ubicación**: `scripts/postgresql/06_indexes_materials.sql`
- **Características**:
  - Idempotente mediante `IF NOT EXISTS`
  - Documentado con comentarios explicativos
  - Secuencialmente numerado para control de versiones
- **Interacciones**: Es ejecutado manualmente por DBA o automatizado en pipeline CI/CD

**3. Material Repository** (Componente existente beneficiado)
- **Responsabilidad**: Abstracción de acceso a datos de materiales
- **Tecnología**: Go con driver `lib/pq`
- **Beneficio**: Queries con ordenamiento cronológico automáticamente usan el índice
- **Interacciones**: Es invocado por `MaterialService` en la capa de aplicación

### Interacciones

El flujo de interacción beneficiado por el índice es:

1. **Cliente HTTP** → Realiza petición `GET /api/materials?sort=updated_at&order=desc&limit=20`
2. **Material Handler** → Recibe request, valida parámetros, invoca servicio
3. **Material Service** → Aplica lógica de negocio, invoca repositorio
4. **Material Repository** → Ejecuta query SQL: `SELECT * FROM materials ORDER BY updated_at DESC LIMIT 20`
5. **PostgreSQL Query Planner** → Detecta índice `idx_materials_updated_at`, lo usa en lugar de full table scan
6. **PostgreSQL Executor** → Retorna resultados optimizados
7. **Repository → Service → Handler → Cliente** → Respuesta JSON con latencia reducida

La clave es que el índice es transparente para el código Go: ninguna capa de aplicación necesita modificarse, el optimizador de PostgreSQL automáticamente selecciona el índice cuando es beneficioso.

## Modelo de Datos

### Estrategia de Persistencia
**PostgreSQL Relacional** - Base de datos principal del sistema, adecuada para datos estructurados con relaciones claras y necesidad de transacciones ACID.

### Entidades Principales

Para este sprint, solo se trabaja con la entidad **Materials**:

**Tabla: `materials`**

**Descripción**: Almacena materiales educativos del sistema (PDFs, videos, enlaces, etc.) asociados a cursos y módulos.

**Atributos clave**:
| Campo | Tipo | Restricciones | Descripción |
|-------|------|---------------|-------------|
| id | UUID | PK, NOT NULL | Identificador único del material |
| title | VARCHAR(255) | NOT NULL | Título del material educativo |
| description | TEXT | NULLABLE | Descripción o resumen del contenido |
| type | VARCHAR(50) | NOT NULL | Tipo: 'pdf', 'video', 'link', 'document' |
| url | TEXT | NOT NULL | URL o ruta del material (S3 o externo) |
| course_id | UUID | FK, NOT NULL | Relación con tabla courses |
| module_id | UUID | FK, NULLABLE | Relación opcional con módulo específico |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | Fecha de creación |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | Fecha de última modificación |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'active' | Estado: 'active', 'archived', 'deleted' |

**Índices existentes** (antes de este sprint):
- `PRIMARY KEY (id)` - Índice automático por PK
- `idx_materials_course_id` - Para filtros por curso (asumido)
- `idx_materials_type` - Para filtros por tipo de material (asumido)

**Nuevo índice** (agregado en este sprint):
- `idx_materials_updated_at` (DESC) - Para ordenamiento cronológico descendente

**Justificación del índice DESC**:
La mayoría de las queries de listado de materiales ordenan por fecha de actualización más reciente primero (`ORDER BY updated_at DESC`). Un índice descendente permite a PostgreSQL leer las entradas del índice en orden directo sin necesidad de recorrerlas en reversa, optimizando este patrón de acceso común.

### Relaciones

**Materials** tiene las siguientes relaciones (no modificadas en este sprint):

1. **Materials N:1 Courses**
   - Un material pertenece a un curso
   - Campo: `materials.course_id` → `courses.id`
   - Comportamiento: ON DELETE CASCADE (si se borra curso, se borran materiales)

2. **Materials N:1 Modules** (opcional)
   - Un material puede estar asociado a un módulo específico dentro del curso
   - Campo: `materials.module_id` → `modules.id`
   - Comportamiento: ON DELETE SET NULL

3. **Materials 1:N Material_Files** (posible tabla relacionada)
   - Un material puede tener múltiples archivos/versiones
   - Relación uno a muchos para manejar versiones o archivos adjuntos

### Queries Beneficiadas por el Índice

El índice optimiza estos patrones de consulta frecuentes:

**Query 1: Listado de materiales recientes (más común)**
```sql
SELECT id, title, type, updated_at, url
FROM materials
WHERE status = 'active'
ORDER BY updated_at DESC
LIMIT 20 OFFSET 0;
```
**Mejora esperada**: El índice permite escaneo directo de las últimas 20 entradas sin recorrer toda la tabla.

**Query 2: Materiales actualizados en rango de fechas**
```sql
SELECT *
FROM materials
WHERE updated_at >= '2025-01-01'
  AND updated_at < '2025-02-01'
  AND course_id = 'xxx-xxx-xxx'
ORDER BY updated_at DESC;
```
**Mejora esperada**: Index range scan en `updated_at` combinado con filtro de `course_id`.

**Query 3: Últimos materiales por tipo**
```sql
SELECT *
FROM materials
WHERE type = 'video'
ORDER BY updated_at DESC
LIMIT 10;
```
**Mejora esperada**: Escaneo del índice con filtro adicional aplicado.

### Impacto en Performance

**Antes del índice**:
- PostgreSQL debe realizar **Seq Scan** (escaneo secuencial completo de la tabla)
- Si la tabla tiene 10,000 materiales, debe leer todos para ordenarlos
- Costo estimado: O(n log n) para ordenamiento completo
- Tiempo estimado: 50-200ms dependiendo del tamaño de tabla

**Después del índice**:
- PostgreSQL usa **Index Scan** en `idx_materials_updated_at`
- Lee solo las primeras N entradas del índice (ya ordenadas)
- Costo estimado: O(log n + k) donde k es el LIMIT
- Tiempo estimado: 5-20ms (mejora de 5-10x)

La mejora real depende del tamaño de la tabla `materials` y la selectividad de filtros adicionales.

## Flujo de Procesos

Este sprint implementa un flujo de optimización de base de datos, no un flujo de negocio de usuario. Sin embargo, podemos describir dos flujos principales:

### Proceso 1: Aplicación del Índice (Deployment)

**Descripción**: Flujo técnico de aplicación de la migración de base de datos.

**Pasos detallados**:

1. **Desarrollo: Crear Script SQL**
   - Desarrollador crea archivo `scripts/postgresql/06_indexes_materials.sql`
   - Script incluye:
     - Comentarios explicativos del propósito
     - Sentencia `CREATE INDEX IF NOT EXISTS`
     - Definición de índice descendente
     - Comentarios de verificación
   - Desarrollador valida sintaxis SQL localmente

2. **Validación Local: Ejecutar en Entorno de Desarrollo**
   - Desarrollador ejecuta: `psql -d edugo_db_local -f scripts/postgresql/06_indexes_materials.sql`
   - PostgreSQL procesa el comando:
     - Si índice no existe: lo crea (puede tomar varios segundos en tablas grandes)
     - Si índice ya existe: retorna sin error (idempotencia)
   - Desarrollador verifica creación: `SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials';`
   - Salida esperada: `idx_materials_updated_at | CREATE INDEX idx_materials_updated_at ON materials USING btree (updated_at DESC)`

3. **Prueba de Performance: EXPLAIN ANALYZE**
   - Desarrollador ejecuta query de prueba:
     ```sql
     EXPLAIN ANALYZE
     SELECT * FROM materials
     ORDER BY updated_at DESC
     LIMIT 20;
     ```
   - **Antes del índice (esperado)**:
     ```
     Sort  (cost=1500..1550 rows=10000)
       Sort Key: updated_at DESC
       ->  Seq Scan on materials  (cost=0.00..1000 rows=10000)
     Execution Time: 45.234 ms
     ```
   - **Después del índice (esperado)**:
     ```
     Limit  (cost=0.00..2.50 rows=20)
       ->  Index Scan using idx_materials_updated_at on materials  (cost=0.00..1250 rows=10000)
     Execution Time: 3.125 ms
     ```
   - Desarrollador documenta mejora de performance (opcional)

4. **Control de Versiones: Commit y Push**
   - Desarrollador agrega script al staging: `git add scripts/postgresql/06_indexes_materials.sql`
   - Crea commit atómico: `git commit -m "perf(db): agregar índice en materials.updated_at para optimizar ordenamiento"`
   - Push al repositorio remoto: `git push origin fix/debug-sprint-commands`

5. **Integración Continua: Pipeline CI/CD**
   - Pipeline detecta cambio en carpeta `scripts/postgresql/`
   - (Opcional) Ejecuta linter SQL para validar sintaxis
   - (Opcional) Ejecuta script en base de datos de testing
   - Build de aplicación Go: `go build ./...` (debe pasar)
   - Tests: `go test ./...` (no deberían fallar)

6. **Deployment a QA/Staging**
   - Pipeline de deployment ejecuta script en base de datos QA
   - Comando: `psql -h qa-db.edugo.com -d edugo_db -f scripts/postgresql/06_indexes_materials.sql`
   - Verificación automática (opcional):
     - Query para confirmar existencia de índice
     - EXPLAIN ANALYZE de query de prueba
     - Alerta si índice no fue creado

7. **Deployment a Producción**
   - DBA o pipeline automatizado ejecuta script en producción
   - Consideraciones:
     - En tablas grandes (>1M registros), `CREATE INDEX` puede tardar minutos
     - En PostgreSQL, `CREATE INDEX CONCURRENTLY` permite creación sin bloquear tabla (no usado aquí por ser operación rápida)
     - Monitorear carga de CPU durante creación
   - Post-deployment:
     - Verificar índice creado
     - Monitorear latencia de endpoints de listado de materiales
     - Confirmar reducción de tiempo de respuesta en APM/logs

8. **Monitoreo Post-Deployment**
   - Equipo DevOps observa métricas:
     - Latencia P50/P95/P99 de endpoint `/api/materials`
     - Query performance en logs de PostgreSQL (pg_stat_statements)
     - Uso de índice: queries que usan `idx_materials_updated_at`
   - Confirmación de éxito: latencia reducida en 30-70%

**Resultado**: Índice creado y activo en producción, queries optimizadas automáticamente.

---

### Proceso 2: Query Optimizada en Runtime (Usuario Final)

**Descripción**: Flujo de una petición de listado de materiales recientes después de aplicado el índice.

**Pasos detallados**:

1. **Solicitud del Cliente**
   - Usuario abre aplicación móvil EduGo
   - Navega a sección "Mis Materiales" o "Materiales Recientes"
   - Aplicación frontend (React Native/Flutter) envía request HTTP:
     ```
     GET /api/materials?sort=updated_at&order=desc&limit=20
     Authorization: Bearer eyJhbGc...
     ```

2. **Recepción en API Gateway/Load Balancer**
   - Request llega a infraestructura AWS (ALB o API Gateway)
   - Load balancer distribuye a una instancia de EduGo API Mobile
   - Request entra a servidor Gin en el puerto 8080

3. **Procesamiento en Material Handler**
   - Handler: `internal/infrastructure/http/handler/material_handler.go`
   - Middleware chain se ejecuta:
     - Logger middleware: registra request
     - Auth middleware: valida JWT del usuario
     - CORS middleware: verifica origen permitido
   - Handler extrae query params:
     - `sort=updated_at` → campo de ordenamiento
     - `order=desc` → dirección descendente
     - `limit=20` → tamaño de página
   - Handler valida parámetros (tipos, valores permitidos)
   - Handler construye DTO: `ListMaterialsRequest{Sort: "updated_at", Order: "desc", Limit: 20}`

4. **Invocación del Material Service**
   - Handler invoca: `materialService.ListMaterials(ctx, request)`
   - Service: `internal/application/material_service.go`
   - Service aplica lógica de negocio:
     - Valida que el usuario tiene permiso para ver materiales
     - Aplica filtros adicionales (ej: solo materiales activos, solo de cursos del usuario)
     - Construye query filters: `Filters{Status: "active", UserID: "xxx", Sort: "updated_at", Order: "desc"}`

5. **Ejecución en Material Repository**
   - Service invoca: `materialRepo.FindAll(ctx, filters)`
   - Repository: `internal/infrastructure/persistence/postgresql/material_repository.go`
   - Repository construye query SQL dinámicamente:
     ```sql
     SELECT m.id, m.title, m.type, m.url, m.updated_at, m.course_id
     FROM materials m
     INNER JOIN course_users cu ON m.course_id = cu.course_id
     WHERE m.status = $1
       AND cu.user_id = $2
     ORDER BY m.updated_at DESC
     LIMIT $3;
     ```
   - Repository ejecuta query con driver `lib/pq`:
     ```go
     rows, err := r.db.QueryContext(ctx, query, "active", userID, 20)
     ```

6. **Optimización en PostgreSQL Query Planner** (★ PUNTO DE IMPACTO DEL ÍNDICE)
   - PostgreSQL recibe query SQL
   - **Query Planner analiza query**:
     - Identifica ordenamiento: `ORDER BY m.updated_at DESC`
     - Busca índices aplicables en tabla `materials`
     - Detecta índice `idx_materials_updated_at` con dirección DESC matching
     - Calcula costo de planes alternativos:
       - **Plan A (sin índice)**: Seq Scan → Sort en memoria (cost: 1500)
       - **Plan B (con índice)**: Index Scan en `idx_materials_updated_at` → Limit (cost: 50)
     - **Selecciona Plan B** por menor costo
   - **Query Executor ejecuta Plan B**:
     - Abre cursor en índice `idx_materials_updated_at` desde el inicio (valor más reciente)
     - Lee secuencialmente las primeras 20 entradas del índice (ya ordenadas)
     - Por cada entrada del índice:
       - Accede al heap de la tabla para obtener columnas completas
       - Aplica filtros adicionales (`status = 'active'`, `user_id = 'xxx'`)
       - Si pasa filtros: agrega a resultado
     - Detiene lectura después de obtener 20 resultados (LIMIT satisfecho)
   - Retorna resultset al driver `lib/pq`

7. **Construcción de Entidades en Repository**
   - Repository recibe `sql.Rows` de PostgreSQL
   - Itera sobre rows y construye entidades de dominio:
     ```go
     materials := make([]*domain.Material, 0, 20)
     for rows.Next() {
       var m domain.Material
       err := rows.Scan(&m.ID, &m.Title, &m.Type, &m.URL, &m.UpdatedAt, &m.CourseID)
       materials = append(materials, &m)
     }
     ```
   - Cierra conexión a base de datos
   - Retorna slice de materiales al service: `return materials, nil`

8. **Transformación en Service**
   - Service recibe entidades de dominio
   - Aplica lógica adicional si es necesaria (ej: enriquecer con datos de caché)
   - Transforma entidades a DTOs de respuesta:
     ```go
     response := &dto.ListMaterialsResponse{
       Materials: toMaterialDTOs(materials),
       Total: len(materials),
     }
     ```
   - Retorna response al handler

9. **Serialización en Handler**
   - Handler recibe response del service
   - Serializa a JSON:
     ```go
     c.JSON(http.StatusOK, gin.H{
       "success": true,
       "data": response,
       "timestamp": time.Now(),
     })
     ```
   - Logger middleware registra respuesta exitosa con latencia

10. **Respuesta al Cliente**
    - Response HTTP viaja de vuelta:
      ```json
      HTTP/1.1 200 OK
      Content-Type: application/json
      X-Response-Time: 15ms

      {
        "success": true,
        "data": {
          "materials": [
            {"id": "...", "title": "Video Clase 10", "updated_at": "2025-11-04T10:30:00Z"},
            {"id": "...", "title": "PDF Resumen", "updated_at": "2025-11-03T15:20:00Z"},
            ...
          ],
          "total": 20
        },
        "timestamp": "2025-11-04T11:00:00Z"
      }
      ```
    - Aplicación móvil recibe JSON y renderiza lista de materiales
    - Usuario ve materiales ordenados por más recientes primero con latencia reducida (15ms vs 50ms antes del índice)

**Resultado**: Usuario obtiene listado de materiales recientes con latencia optimizada, experiencia fluida.

---

### Flujos Alternativos/Excepcionales

**Caso A: Índice Ya Existe (Idempotencia)**
- Al ejecutar el script, PostgreSQL detecta que `idx_materials_updated_at` ya existe
- Gracias a `IF NOT EXISTS`, retorna `NOTICE` en lugar de `ERROR`
- Script completa exitosamente sin modificaciones
- Resultado: operación segura para re-ejecutar

**Caso B: Tabla Materials Vacía o Muy Pequeña**
- Si la tabla tiene <1000 registros, PostgreSQL puede elegir **no usar el índice**
- Razón: Seq Scan de tabla pequeña es más rápido que index scan + heap access
- Query planner elige plan óptimo según estadísticas
- Resultado: índice existe pero no se usa (aceptable, preparado para crecimiento)

**Caso C: Query con Filtros Muy Selectivos**
- Si query tiene filtro altamente selectivo (ej: `WHERE course_id = 'xxx'` retorna 5 registros)
- PostgreSQL puede usar índice de `course_id` en lugar de `updated_at`
- Luego ordena en memoria los pocos resultados
- Resultado: índice de `updated_at` no se usa en este caso específico (correcto según optimizador)

**Caso D: Error de Sintaxis en Script**
- Si script SQL tiene error de sintaxis (ej: tipo de dato incorrecto)
- PostgreSQL rechaza creación con mensaje `ERROR: syntax error at or near ...`
- Script falla, índice no se crea
- Developer debe corregir y re-ejecutar
- Resultado: fallo seguro, sin efectos colaterales

## Stack Tecnológico Recomendado

Este sprint utiliza el stack tecnológico ya establecido del proyecto EduGo API Mobile. No se agregan nuevas tecnologías, solo se aplica una feature existente de PostgreSQL.

### Backend
- **Go 1.21+**: Lenguaje de backend (no modificado en este sprint)
- **Gin Framework**: Framework HTTP para API REST (no modificado)
- **Viper**: Gestión de configuración (no modificado)

### Base de Datos
- **PostgreSQL 14+**: Base de datos relacional principal
  - **Feature utilizada**: Índices descendentes (`CREATE INDEX ... DESC`)
  - **Feature utilizada**: Creación condicional (`IF NOT EXISTS`)
  - **Herramienta de análisis**: `EXPLAIN ANALYZE` para validar performance
  - **Driver Go**: `lib/pq` para conexión desde aplicación

### Herramientas de Migración
- **Scripts SQL manuales**: Enfoque actual del proyecto (no usa herramientas como Flyway/Liquibase)
- **Ubicación**: `scripts/postgresql/`
- **Convención**: Numeración secuencial `NN_descripcion.sql`

### DevOps/Deployment
- **Git**: Control de versiones (script SQL bajo control de versiones)
- **psql CLI**: Herramienta para ejecutar scripts manualmente
- **CI/CD Pipeline** (asumido): Automatización de ejecución de scripts en ambientes

### Monitoreo/Observabilidad
- **PostgreSQL pg_stat_statements**: Para analizar queries más costosas
- **Zap Logger** (edugo-shared): Logging estructurado de aplicación
- **APM** (asumido): Herramienta de Application Performance Monitoring para medir latencia de endpoints

### Justificación de Elecciones Tecnológicas

**¿Por qué PostgreSQL en lugar de MongoDB para este caso?**
- La tabla `materials` tiene estructura relacional clara (FK a courses, modules)
- Se beneficia de ACID transactions
- Los índices de PostgreSQL son altamente optimizados para ordenamiento
- El proyecto ya usa PostgreSQL como base de datos principal

**¿Por qué índice DESC en lugar de ASC?**
- El patrón de acceso más común es `ORDER BY updated_at DESC` (más reciente primero)
- Un índice DESC permite lectura directa sin reversa
- PostgreSQL puede usar índice ASC para queries DESC, pero con overhead adicional
- Mejor alineación entre dirección del índice y dirección de la query

**¿Por qué no usar índice compuesto (course_id, updated_at)?**
- Este sprint se enfoca en optimización simple y atómica
- Un índice compuesto beneficiaría queries filtradas por curso + ordenadas por fecha
- Sin embargo, agrega complejidad de decisión al query planner
- Decisión: empezar con índice simple, evaluar índice compuesto en futuro si es necesario

**¿Por qué no usar particionamiento de tabla?**
- Particionamiento es apropiado para tablas muy grandes (>10M registros)
- Agrega complejidad operacional significativa
- La tabla `materials` probablemente no requiere particionamiento aún
- Índice es suficiente para tamaño actual de datos

## Patrones de Diseño Recomendados

Aunque este sprint es principalmente de optimización de base de datos, se aplican ciertos patrones y principios:

### 1. Patrón: Migración Idempotente
**Descripción**: Las operaciones de cambio de esquema deben ser idempotentes, es decir, ejecutarse múltiples veces sin error ni efectos colaterales.

**Implementación**:
```sql
CREATE INDEX IF NOT EXISTS idx_materials_updated_at
ON materials(updated_at DESC);
```

**Justificación**:
- Permite re-ejecutar scripts en caso de fallo parcial
- Facilita deployment automatizado en múltiples ambientes
- Evita errores en rollbacks o re-deploys
- Simplifica testing de scripts

**Alternativas rechazadas**:
- `CREATE INDEX` sin `IF NOT EXISTS`: falla si ya existe (no idempotente)
- Script con `DROP INDEX` previo: más riesgoso (ventana sin índice)

### 2. Principio: Database Performance Tuning Incremental
**Descripción**: Optimizar base de datos mediante cambios pequeños, medibles y reversibles en lugar de refactorizaciones masivas.

**Aplicación en este sprint**:
- Se agrega UN solo índice específico
- Se mide impacto antes y después con EXPLAIN ANALYZE
- El índice es fácil de remover si causa problemas (`DROP INDEX`)
- No se modifica código de aplicación (cambio transparente)

**Beneficios**:
- Bajo riesgo: fácil rollback
- Alto valor: mejora medible de performance
- Aprendizaje: datos concretos de impacto de índices
- Escalabilidad: modelo para futuras optimizaciones

### 3. Patrón: Separation of Concerns (Índice vs Código)
**Descripción**: La optimización se implementa en la capa de base de datos sin modificar lógica de aplicación.

**Ventajas**:
- **Transparencia**: El código Go no necesita saber de la existencia del índice
- **Flexibilidad**: Podemos agregar/remover índices sin redeploy de aplicación
- **Performance**: PostgreSQL automáticamente decide cuándo usar el índice
- **Testing**: No se requieren nuevos tests unitarios en código Go

**Ejemplo de separación correcta**:
```go
// Repository (NO cambia)
func (r *MaterialRepository) FindAll(ctx context.Context, filters Filters) ([]*Material, error) {
    query := "SELECT * FROM materials ORDER BY updated_at DESC LIMIT ?"
    // PostgreSQL decide automáticamente usar el índice
    rows, err := r.db.QueryContext(ctx, query, filters.Limit)
    // ...
}
```

### 4. Principio: Measure First, Optimize Second
**Descripción**: Documentar performance antes y después de la optimización para validar efectividad.

**Implementación en este sprint**:
1. Ejecutar `EXPLAIN ANALYZE` antes del índice (establecer baseline)
2. Crear índice
3. Ejecutar mismo `EXPLAIN ANALYZE` después (medir mejora)
4. Documentar resultados (opcional pero recomendado)

**Ejemplo de documentación**:
```markdown
## Validación de Performance

### Antes del índice:
- Plan: Seq Scan → Sort
- Tiempo: 45.234 ms
- Costo estimado: 1550

### Después del índice:
- Plan: Index Scan using idx_materials_updated_at
- Tiempo: 3.125 ms
- Costo estimado: 50
- **Mejora: 14.5x más rápido**
```

### 5. Patrón: Convention Over Configuration (Naming)
**Descripción**: Seguir convenciones consistentes de nomenclatura para facilitar mantenimiento.

**Aplicación**:
- **Índice**: `idx_{tabla}_{columna(s)}` → `idx_materials_updated_at`
- **Script**: `NN_{tipo}_{tabla}.sql` → `06_indexes_materials.sql`
- **Ubicación**: `scripts/postgresql/` (carpeta estándar)

**Beneficios**:
- Predecibilidad: cualquier developer puede intuir el nombre
- Consistencia: todos los índices siguen el mismo patrón
- Búsqueda: fácil encontrar índices en queries de catálogo
- Escalabilidad: agregar más índices sin conflictos de nombre

## Consideraciones No Funcionales

### Escalabilidad

**Impacto en Write Performance**:
- **Trade-off**: Los índices mejoran SELECTs pero ralentizan INSERTs, UPDATEs y DELETEs
- **Análisis para este caso**:
  - La columna `updated_at` se modifica en cada UPDATE de material
  - Cada UPDATE requiere actualizar el índice además de la fila
  - Overhead estimado: 5-10% más lento en UPDATEs (aceptable)
- **Justificación**: Los materiales se **leen mucho más frecuentemente** que se modifican
  - Ratio estimado: 100 SELECTs por cada 1 UPDATE
  - Mejora de 10x en SELECTs compensa ralentización de 5% en UPDATEs
- **Monitoreo recomendado**: Observar latencia de endpoints de actualización de materiales

**Crecimiento de Tabla**:
- El índice crece proporcionalmente al tamaño de la tabla
- Espacio adicional estimado: ~10-15% del tamaño de la tabla
- Para tabla con 100,000 materiales: ~5-10 MB de espacio adicional
- Consideración: Espacio en disco es económico comparado con mejora de performance

**Escalabilidad Horizontal**:
- Los índices se replican en réplicas de lectura de PostgreSQL
- Read replicas se benefician igualmente del índice
- Queries de listado pueden dirigirse a réplicas para distribuir carga

### Seguridad

**Índice no introduce vulnerabilidades**:
- No modifica permisos de tabla
- No expone datos adicionales
- No abre puertos ni servicios
- **Conclusión**: Impacto de seguridad neutro

**Consideraciones de acceso**:
- El script SQL debe ejecutarse con usuario que tenga permiso `CREATE INDEX`
- En producción, evitar usar superuser, preferir usuario con permisos limitados
- Auditar logs de PostgreSQL para rastrear ejecución del script

**SQL Injection**:
- No aplicable: script SQL estático sin parámetros dinámicos
- El código Go que ejecuta queries sigue usando prepared statements

### Performance

**Mejora Esperada**:
- **Query de listado (ORDER BY updated_at DESC LIMIT N)**:
  - Antes: 50-200ms (Seq Scan + Sort)
  - Después: 5-20ms (Index Scan)
  - Mejora: 5-10x más rápido
- **Query con filtros + ordenamiento**:
  - Mejora variable según selectividad de filtros
  - En el mejor caso (filtros poco selectivos): mejora similar a listado simple
  - En el peor caso (filtros muy selectivos): índice puede no usarse (correcto según optimizador)

**Impacto en Memoria**:
- Índice se carga en memoria PostgreSQL (shared_buffers) bajo demanda
- Índice de 10 MB consume 10 MB de cache
- Beneficio: páginas de índice más usadas permanecen en cache (hot)

**Impacto en CPU**:
- Creación de índice: spike temporal de CPU durante `CREATE INDEX`
- Queries con índice: **menos CPU** que sin índice (evita sort en memoria)

**Recomendaciones de monitoreo**:
- Configurar alerta si latencia de endpoint `/api/materials` excede 100ms P95
- Usar `pg_stat_statements` para identificar queries más lentas
- Revisar plan de ejecución de queries periódicamente (puede cambiar con estadísticas)

### Mantenibilidad

**Documentación**:
- El script SQL debe incluir comentarios explicativos:
  - Propósito del índice
  - Queries que se benefician
  - Instrucciones de verificación
- Documentar resultado de EXPLAIN ANALYZE en commit o wiki (recomendado)

**Reversibilidad**:
- Rollback es simple y seguro:
  ```sql
  DROP INDEX IF EXISTS idx_materials_updated_at;
  ```
- El DROP INDEX es instantáneo (no requiere escaneo de tabla)
- La aplicación sigue funcionando sin el índice (performance degradada pero funcional)

**Testing**:
- **Tests unitarios de código Go**: No requieren modificación (índice es transparente)
- **Tests de integración**: Pueden beneficiarse del índice en queries de prueba
- **Tests de performance**: Ejecutar benchmark antes/después del índice (recomendado)

**Versionado de Scripts**:
- Seguir numeración secuencial: `06_indexes_materials.sql`
- Próximos scripts: `07_...`, `08_...`
- Mantener orden cronológico de aplicación
- No modificar scripts ya aplicados en producción (crear nuevos scripts para cambios)

**Observabilidad**:
- Índice es visible en catálogo de PostgreSQL:
  ```sql
  SELECT * FROM pg_indexes WHERE tablename = 'materials';
  ```
- Uso de índice es rastreable en logs con `log_statement = 'all'` (no recomendado en producción por volumen)
- Mejor opción: `pg_stat_statements` para métricas agregadas

### Compatibilidad

**Versión mínima de PostgreSQL**:
- Índices descendentes (`DESC`) soportados desde PostgreSQL 8.3 (2008)
- `IF NOT EXISTS` en `CREATE INDEX` soportado desde PostgreSQL 9.5 (2016)
- Requisito: PostgreSQL 9.5+ (ampliamente cumplido)

**Compatibilidad con ORM**:
- El proyecto no usa ORM (usa SQL directo con `lib/pq`)
- Índice es completamente transparente para el código Go
- No requiere cambios en queries ni en driver

**Compatibilidad con ambientes**:
- Script funciona idénticamente en dev, QA, staging y producción
- Consideración: tiempo de creación varía según tamaño de tabla por ambiente

## Riesgos Identificados

### Riesgo 1: Índice No Se Usa en Queries Esperadas
**Probabilidad**: Baja
**Impacto**: Medio (índice creado pero no aporta valor)

**Descripción**:
PostgreSQL puede elegir no usar el índice si su optimizador considera que un Seq Scan es más eficiente. Esto puede ocurrir si:
- La tabla `materials` es muy pequeña (<1000 registros)
- Las estadísticas de la tabla están desactualizadas
- La configuración de PostgreSQL penaliza el uso de índices (ej: `random_page_cost` muy alto)

**Mitigación**:
- Ejecutar `EXPLAIN ANALYZE` inmediatamente después de crear el índice para confirmar uso
- Si el índice no se usa:
  - Verificar estadísticas con `ANALYZE materials;`
  - Revisar configuración de `random_page_cost` (recomendado: 1.1 para SSD)
  - Si la tabla es pequeña: aceptar que el índice no se usa aún (se usará al crecer)
- Documentar el resultado esperado vs real

**Plan B**:
Si el índice no aporta valor, removerlo es trivial: `DROP INDEX idx_materials_updated_at;`

---

### Riesgo 2: Degradación de Performance en UPDATEs de Materiales
**Probabilidad**: Media
**Impacto**: Bajo (ralentización aceptable)

**Descripción**:
Cada UPDATE en la tabla `materials` requiere actualizar el índice `idx_materials_updated_at`, lo que añade overhead. Si el sistema tiene un alto volumen de actualizaciones de materiales, esto podría impactar latencia de operaciones de escritura.

**Mitigación**:
- Medir latencia de endpoints de UPDATE antes y después del índice
- Monitorear métricas de `UPDATE` en `pg_stat_statements`
- Estimación conservadora: overhead de 5-10% en UPDATEs (generalmente imperceptible)
- Justificación: ratio de lectura/escritura es aproximadamente 100:1 (los materiales se consultan mucho más que se editan)

**Plan B**:
Si el overhead en escrituras es inaceptable (poco probable), evaluar:
- Índice parcial: `CREATE INDEX ... WHERE status = 'active'` (reduce tamaño)
- Índice compuesto más selectivo
- Remover índice si el trade-off no es favorable

---

### Riesgo 3: Tiempo de Creación de Índice Bloquea Tabla en Producción
**Probabilidad**: Baja
**Impacto**: Alto (si ocurre, bloqueo de operaciones)

**Descripción**:
El comando `CREATE INDEX` adquiere un lock de nivel `SHARE` en la tabla `materials`, permitiendo lecturas pero bloqueando escrituras (INSERTs, UPDATEs, DELETEs) durante la creación del índice. En tablas grandes (>1M registros), esto puede tomar varios minutos.

**Mitigación**:
- **Opción 1 (usada en este sprint)**: `CREATE INDEX` estándar
  - Apropiado si la tabla `materials` es pequeña (<100K registros)
  - Tiempo estimado de creación: 5-30 segundos (aceptable)
  - Ventana de bloqueo de escrituras: mínima

- **Opción 2 (si tabla es grande)**: Usar `CREATE INDEX CONCURRENTLY`
  ```sql
  CREATE INDEX CONCURRENTLY idx_materials_updated_at
  ON materials(updated_at DESC);
  ```
  - No bloquea escrituras
  - Tarda más tiempo pero sin impacto en disponibilidad
  - Requiere dos scans de tabla en lugar de uno

**Plan de Ejecución en Producción**:
1. Ejecutar script en horario de bajo tráfico (ej: madrugada)
2. Monitorear logs de PostgreSQL durante creación
3. Tener plan de rollback listo (`DROP INDEX` en caso de problema)
4. Considerar `CONCURRENTLY` si tabla tiene >500K registros

---

### Riesgo 4: Aumento de Espacio en Disco de Base de Datos
**Probabilidad**: Alta (certeza)
**Impacto**: Muy Bajo (espacio adicional mínimo)

**Descripción**:
El índice consume espacio adicional en disco, estimado en 10-15% del tamaño de la tabla `materials`.

**Mitigación**:
- Verificar espacio disponible en disco antes de crear índice:
  ```sql
  SELECT pg_size_pretty(pg_total_relation_size('materials'));
  ```
- Estimación: tabla de 100 MB → índice de 10-15 MB (despreciable)
- Monitorear uso de disco de PostgreSQL con herramientas de infraestructura
- Alertar si uso de disco supera 80% de capacidad

**Plan B**:
Si espacio en disco se convierte en problema (poco probable para un índice):
- Remover índice: libera espacio inmediatamente
- Considerar compresión de datos o aumento de capacidad de disco

---

### Riesgo 5: Script SQL con Error de Sintaxis Rompe Pipeline de Deployment
**Probabilidad**: Baja (mitigado por validación)
**Impacto**: Medio (deployment falla)

**Descripción**:
Si el script `06_indexes_materials.sql` contiene error de sintaxis, el pipeline de CI/CD fallará al intentar ejecutar el script, bloqueando el deployment.

**Mitigación**:
- **Validación local**: Desarrollador ejecuta script en entorno local antes de commit
- **Linter SQL** (opcional): Integrar herramienta como `sqlfluff` en CI/CD
- **Test en QA primero**: Script se ejecuta en ambiente QA antes de producción
- **Idempotencia**: `IF NOT EXISTS` asegura que re-ejecución no causa problemas

**Plan B**:
Si script falla en deployment:
1. Rollback del deployment (aplicación anterior sigue funcionando)
2. Corregir sintaxis del script
3. Re-ejecutar pipeline

---

### Riesgo 6: Índice No Se Documenta y Se Pierde Contexto
**Probabilidad**: Media
**Impacto**: Bajo (mantenibilidad reducida)

**Descripción**:
Sin documentación adecuada, futuros developers pueden:
- No entender por qué existe el índice
- Considerar removerlo por "limpieza" sin medir impacto
- No saber qué queries se benefician

**Mitigación**:
- **Comentarios en script SQL**: Explicar propósito y queries beneficiadas
- **Commit message descriptivo**: `perf(db): agregar índice en materials.updated_at para optimizar ordenamiento`
- **Documentación de sprint** (este archivo): Registro permanente de decisión arquitectónica
- **Wiki/ADR** (opcional): Crear Architecture Decision Record para decisión de índice

**Plan B**:
Si en el futuro hay dudas sobre el índice:
- Usar `EXPLAIN ANALYZE` para verificar si se usa
- Consultar historial de git para encontrar este sprint
- Ejecutar query sin índice temporalmente para medir degradación

## Siguientes Pasos Recomendados

### Paso 1: Ejecutar Implementación del Sprint ✅
**Responsable**: Developer/Agent
**Duración**: 10-15 minutos

1. Crear archivo `scripts/postgresql/06_indexes_materials.sql` con contenido especificado
2. Validar sintaxis ejecutando en base de datos local
3. Verificar índice creado con query de catálogo
4. Ejecutar `EXPLAIN ANALYZE` para confirmar uso del índice
5. Documentar resultado (opcional)
6. Commit del script con mensaje: `perf(db): agregar índice en materials.updated_at para optimizar ordenamiento`

**Criterio de éxito**: Script ejecutado localmente sin errores, índice verificado, commit creado.

---

### Paso 2: Validar en Ambiente QA
**Responsable**: QA Engineer / DevOps
**Duración**: 20-30 minutos

1. Ejecutar script en base de datos QA:
   ```bash
   psql -h qa-db.edugo.com -d edugo_db -f scripts/postgresql/06_indexes_materials.sql
   ```
2. Verificar índice creado:
   ```sql
   SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials';
   ```
3. Ejecutar query de prueba con EXPLAIN ANALYZE:
   ```sql
   EXPLAIN ANALYZE SELECT * FROM materials ORDER BY updated_at DESC LIMIT 20;
   ```
4. Confirmar que plan de ejecución usa `idx_materials_updated_at`
5. Ejecutar tests de integración de la aplicación
6. Medir latencia de endpoint `/api/materials` antes y después

**Criterio de éxito**: Tests pasan, latencia reducida, plan de ejecución correcto.

---

### Paso 3: Deployment a Producción
**Responsable**: DevOps / DBA
**Duración**: 15-20 minutos + monitoreo de 24h

1. **Pre-deployment**:
   - Verificar espacio en disco de base de datos (debe tener margen)
   - Confirmar tamaño de tabla `materials`:
     ```sql
     SELECT pg_size_pretty(pg_total_relation_size('materials'));
     ```
   - Planificar ejecución en ventana de bajo tráfico (ej: 2 AM)

2. **Ejecución**:
   ```bash
   psql -h prod-db.edugo.com -d edugo_db -f scripts/postgresql/06_indexes_materials.sql
   ```
   - Monitorear tiempo de creación del índice
   - Verificar índice creado exitosamente

3. **Post-deployment**:
   - Verificar logs de PostgreSQL (sin errores)
   - Ejecutar EXPLAIN ANALYZE de query de prueba
   - Monitorear latencia de endpoint `/api/materials` en APM
   - Observar métricas de `pg_stat_statements` para queries de materiales

4. **Monitoreo continuo (24-48h)**:
   - Dashboard de latencia de endpoints
   - Uso de disco de PostgreSQL
   - Latencia de queries en logs
   - Alertas de error (no debería haber)

**Criterio de éxito**: Índice activo, latencia reducida en producción, sin errores.

---

### Paso 4: Documentar Resultados y Aprendizajes
**Responsable**: Developer / Tech Lead
**Duración**: 30 minutos

1. Actualizar `sprint/current/readme.md` marcando todas las casillas completadas
2. Documentar mejoras de performance medidas:
   - Latencia antes/después del índice
   - Plan de ejecución con EXPLAIN ANALYZE
   - Reducción porcentual de tiempo de respuesta
3. Agregar nota al wiki/ADR del proyecto:
   - Decisión: agregar índice en `materials.updated_at`
   - Contexto: optimizar listados cronológicos
   - Resultado: mejora de Nx
4. Compartir aprendizajes con el equipo:
   - Impacto real de índices en PostgreSQL
   - Proceso de validación con EXPLAIN ANALYZE
   - Mejores prácticas de migración idempotente

**Criterio de éxito**: Documentación completa, equipo informado.

---

### Paso 5: Evaluar Oportunidades de Optimización Adicionales
**Responsable**: Tech Lead / Arquitecto
**Duración**: 1 hora (análisis)

Basándose en los resultados de este sprint, evaluar:

1. **Índices adicionales en otras tablas**:
   - ¿La tabla `courses` se beneficiaría de índice en `updated_at`?
   - ¿La tabla `users` tiene queries con ordenamiento frecuente?

2. **Índices compuestos**:
   - ¿Vale la pena crear índice `(course_id, updated_at)` para queries filtradas por curso?
   - Analizar queries en `pg_stat_statements` para identificar patrones

3. **Índices parciales**:
   - ¿Crear índice solo para materiales activos? `WHERE status = 'active'`
   - Reduce tamaño del índice si el 90% de queries filtran por status

4. **Análisis de queries lentas**:
   - Ejecutar: `SELECT * FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;`
   - Identificar próximas oportunidades de optimización

**Entregable**: Backlog de optimizaciones priorizadas para futuros sprints.

---

### Paso 6: Integrar en Proceso de Desarrollo Estándar
**Responsable**: Tech Lead
**Duración**: Ongoing

1. **Documentar proceso de migración de BD**:
   - Crear guía: "Cómo agregar índices en EduGo"
   - Incluir plantilla de script SQL
   - Documentar proceso de validación con EXPLAIN ANALYZE

2. **Agregar validación en CI/CD**:
   - Linter SQL para detectar errores de sintaxis
   - Test de ejecución de scripts en base de datos temporal
   - Validación de idempotencia (ejecutar script 2 veces sin error)

3. **Establecer criterios de cuándo agregar índices**:
   - Threshold: queries que toman >50ms P95
   - Queries ejecutadas >100 veces/minuto
   - Alto ratio lectura/escritura (>10:1)

4. **Capacitar al equipo**:
   - Workshop: "Optimización de PostgreSQL con índices"
   - Demo: uso de EXPLAIN ANALYZE
   - Mejores prácticas de diseño de índices

**Criterio de éxito**: Proceso repetible y documentado para futuras optimizaciones.

---

## Resumen de Valor Entregado

| Aspecto | Antes | Después | Mejora |
|---------|-------|---------|--------|
| **Latencia de listado** | 50-200ms | 5-20ms | 5-10x más rápido |
| **Plan de ejecución** | Seq Scan + Sort | Index Scan | Óptimo |
| **Experiencia de usuario** | Listados lentos | Respuesta instantánea | Perceptible |
| **Escalabilidad** | Degrada con crecimiento | Escalable con índice | Preparado |
| **Costo de mantenimiento** | N/A | Mínimo | Transparente |

Este sprint, aunque pequeño en alcance, establece un modelo de optimización continua basado en datos y mediciones concretas, con riesgo mínimo y alto retorno de inversión.

---

💡 **Nota**: Este es un análisis rápido sin diagramas. Para análisis completo con diagramas visuales de arquitectura, modelo de datos y flujos de proceso, ejecutar:

```bash
/01-analysis --mode=full
```

---

**Análisis generado el**: 2025-11-04
**Configuración**:
- MODE: quick
- SCOPE: complete
- SOURCE: sprint/current/readme.md

**Estado del sprint**: Pendiente de ejecución
**Branch**: fix/debug-sprint-commands
**Complejidad estimada**: Baja
**Tiempo estimado**: 10-15 minutos
