# Dependencias del Sprint 01 - Schema de Base de Datos

## Dependencias Técnicas Previas

### PostgreSQL
- [ ] **PostgreSQL 15+** instalado y corriendo
- [ ] Usuario con permisos `CREATE TABLE`, `CREATE INDEX`, `ALTER TABLE`
- [ ] Base de datos `edugo_test` creada para testing
- [ ] Base de datos `edugo_dev` creada para desarrollo

```bash
# Verificar versión de PostgreSQL
psql --version
# Output esperado: psql (PostgreSQL) 15.x o superior

# Verificar que servidor está corriendo
pg_isready -h localhost -p 5432
# Output esperado: localhost:5432 - accepting connections

# Crear bases de datos si no existen
psql -U postgres -c "CREATE DATABASE edugo_dev;"
psql -U postgres -c "CREATE DATABASE edugo_test;"
```

### Función gen_uuid_v7()
- [ ] Función `gen_uuid_v7()` disponible en PostgreSQL
- [ ] Migración `01_base_tables.sql` ejecutada previamente

```bash
# Verificar que función existe
psql -U postgres -d edugo_test -c "SELECT gen_uuid_v7();"
# Output esperado: Un UUID tipo v7 (ej: 01936d9a-7f8e-7000-a000-123456789abc)

# Si no existe, crearla
psql -U postgres -d edugo_test -c "
CREATE OR REPLACE FUNCTION gen_uuid_v7()
RETURNS UUID
AS \$\$
DECLARE
    unix_ts_ms BIGINT;
    uuid_bytes BYTEA;
BEGIN
    unix_ts_ms := (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT;
    uuid_bytes := E'\\\\x' || 
                  lpad(to_hex((unix_ts_ms >> 32)::INT), 8, '0') ||
                  lpad(to_hex((unix_ts_ms & 4294967295)::INT), 8, '0') ||
                  encode(gen_random_bytes(8), 'hex');
    RETURN CAST(substring(uuid_bytes::TEXT FROM 3) AS UUID);
END;
\$\$ LANGUAGE plpgsql VOLATILE;
"
```

---

## Dependencias de Tablas Existentes

### Tabla: materials
- [ ] Tabla `materials` existe en PostgreSQL
- [ ] Columnas requeridas: `id`, `title`, `content`, `processing_status`, `created_at`
- [ ] Constraint: `processing_status` tipo ENUM o VARCHAR con valores ('pending', 'processing', 'completed', 'failed')

```bash
# Verificar que tabla materials existe
psql -U postgres -d edugo_test -c "\d materials"

# Verificar columnas requeridas
psql -U postgres -d edugo_test -c "
    SELECT column_name, data_type 
    FROM information_schema.columns 
    WHERE table_name = 'materials' 
    AND column_name IN ('id', 'title', 'processing_status');
"

# Verificar que hay datos de prueba
psql -U postgres -d edugo_test -c "SELECT COUNT(*) FROM materials;"
# Output esperado: Al menos 3 materiales
```

**Si tabla no existe, crearla:**
```sql
CREATE TABLE materials (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    processing_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insertar materiales de prueba
INSERT INTO materials (title, content, processing_status) VALUES
('Introducción a Pascal', 'Contenido sobre Pascal...', 'completed'),
('Fundamentos de Python', 'Contenido sobre Python...', 'completed'),
('Algoritmos de Ordenamiento', 'Contenido sobre algoritmos...', 'completed');
```

### Tabla: users
- [ ] Tabla `users` existe en PostgreSQL
- [ ] Columnas requeridas: `id`, `email`, `name`, `role`, `created_at`
- [ ] Constraint: `role` tipo ENUM o VARCHAR con valores ('student', 'teacher', 'admin')

```bash
# Verificar que tabla users existe
psql -U postgres -d edugo_test -c "\d users"

# Verificar columnas requeridas
psql -U postgres -d edugo_test -c "
    SELECT column_name, data_type 
    FROM information_schema.columns 
    WHERE table_name = 'users' 
    AND column_name IN ('id', 'email', 'role');
"

# Verificar que hay estudiantes de prueba
psql -U postgres -d edugo_test -c "SELECT COUNT(*) FROM users WHERE role = 'student';"
# Output esperado: Al menos 1 estudiante
```

**Si tabla no existe, crearla:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'student',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insertar usuarios de prueba
INSERT INTO users (email, name, role) VALUES
('student1@edugo.com', 'Juan Pérez', 'student'),
('student2@edugo.com', 'María González', 'student'),
('teacher1@edugo.com', 'Prof. Rodriguez', 'teacher');
```

---

## Dependencias de MongoDB (Opcional para Sprint 01)

### MongoDB 7.0+
- [ ] MongoDB instalado y corriendo (para contexto, no requerido en Sprint 01)
- [ ] Colección `material_assessment` existente (se valida en Sprints posteriores)

```bash
# Verificar MongoDB (opcional)
mongosh --version
# Output esperado: 2.x.x o superior

# Verificar conexión (opcional)
mongosh "mongodb://localhost:27017" --eval "db.runCommand({ ping: 1 })"
```

**Nota:** Sprint 01 solo crea el campo `mongo_document_id` en tabla `assessment`. La integración real con MongoDB se hace en Sprint 03.

---

## Herramientas de Desarrollo

### psql CLI
- [ ] Cliente `psql` disponible en PATH

```bash
# Verificar psql
which psql
# Output esperado: /usr/local/bin/psql o similar

psql --version
# Output esperado: psql (PostgreSQL) 15.x
```

### Editores SQL (Opcional)
- DBeaver, pgAdmin, TablePlus, o similar para inspección visual
- No requerido, pero útil para debugging

---

## Variables de Entorno

### Ambiente de Desarrollo
```bash
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="edugo_dev"
export DB_SSLMODE="disable"

# Para tests
export DB_TEST_NAME="edugo_test"
```

**Archivo `.env.local` (ejemplo):**
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/edugo_dev?sslmode=disable
DATABASE_TEST_URL=postgres://postgres:postgres@localhost:5432/edugo_test?sslmode=disable
```

---

## Dependencias de Código (Post-Sprint 01)

**Nota:** Sprint 01 es puramente SQL. Las siguientes dependencias Go se requieren en sprints posteriores:

### Go Packages (Sprints 02-04)
```bash
# GORM (ORM)
go get -u gorm.io/gorm@v1.25.5
go get -u gorm.io/driver/postgres@v1.5.4

# MongoDB Driver (Sprint 03)
go get go.mongodb.org/mongo-driver/mongo@v1.13.1

# Testcontainers (Sprint 05)
go get github.com/testcontainers/testcontainers-go@v0.27.0
go get github.com/testcontainers/testcontainers-go/modules/postgres@v0.27.0
```

---

## Dependencias de Migraciones Previas

### Orden de Ejecución de Migraciones
La migración `06_assessments.sql` debe ejecutarse **después** de:

1. **01_base_tables.sql** - Tablas base (materials, users, función gen_uuid_v7())
2. **02_auth.sql** - Autenticación y roles (si existe)
3. **03_schools.sql** - Escuelas y unidades académicas (si existe)
4. **04_materials.sql** - Extensiones de materials (si existe)
5. **05_progress.sql** - Sistema de progreso (si existe)

```bash
# Verificar migraciones ejecutadas (PostgreSQL 15+ con pg_stat_statements)
psql -U postgres -d edugo_test -c "
    SELECT schemaname, tablename, tableowner 
    FROM pg_tables 
    WHERE schemaname = 'public' 
    ORDER BY tablename;
"

# Migraciones mínimas requeridas
# - Tabla materials
# - Tabla users
# - Función gen_uuid_v7()
```

---

## Verificación de Dependencias

### Script de Verificación Completa
Crear archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/verify_sprint01_deps.sql`

```sql
-- Verificación de Dependencias - Sprint 01
DO $$
DECLARE
    materials_count INTEGER;
    users_count INTEGER;
    has_gen_uuid_v7 BOOLEAN;
BEGIN
    RAISE NOTICE '=== VERIFICANDO DEPENDENCIAS SPRINT 01 ===';
    
    -- 1. Verificar función gen_uuid_v7()
    SELECT EXISTS (
        SELECT 1 FROM pg_proc WHERE proname = 'gen_uuid_v7'
    ) INTO has_gen_uuid_v7;
    
    IF NOT has_gen_uuid_v7 THEN
        RAISE EXCEPTION '❌ FALTA: Función gen_uuid_v7() no encontrada';
    ELSE
        RAISE NOTICE '✅ OK: Función gen_uuid_v7() existe';
    END IF;
    
    -- 2. Verificar tabla materials
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'materials') THEN
        RAISE EXCEPTION '❌ FALTA: Tabla materials no existe';
    ELSE
        SELECT COUNT(*) INTO materials_count FROM materials;
        RAISE NOTICE '✅ OK: Tabla materials existe (% filas)', materials_count;
    END IF;
    
    -- 3. Verificar tabla users
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'users') THEN
        RAISE EXCEPTION '❌ FALTA: Tabla users no existe';
    ELSE
        SELECT COUNT(*) INTO users_count FROM users WHERE role = 'student';
        RAISE NOTICE '✅ OK: Tabla users existe (% estudiantes)', users_count;
    END IF;
    
    -- 4. Verificar versión PostgreSQL
    IF (SELECT current_setting('server_version_num')::INTEGER) < 150000 THEN
        RAISE WARNING '⚠️  PostgreSQL < 15 detectado. Recomendado: 15+';
    ELSE
        RAISE NOTICE '✅ OK: PostgreSQL 15+ detectado';
    END IF;
    
    RAISE NOTICE '=== TODAS LAS DEPENDENCIAS SATISFECHAS ===';
END $$;
```

**Ejecutar verificación:**
```bash
psql -U postgres -d edugo_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/verify_sprint01_deps.sql
```

---

## Checklist de Pre-Ejecución

Antes de ejecutar `06_assessments.sql`:

- [ ] PostgreSQL 15+ instalado y corriendo
- [ ] Bases de datos `edugo_dev` y `edugo_test` creadas
- [ ] Función `gen_uuid_v7()` disponible
- [ ] Tabla `materials` existe con datos de prueba
- [ ] Tabla `users` existe con al menos 1 estudiante
- [ ] Variables de entorno configuradas
- [ ] Permisos de usuario PostgreSQL correctos
- [ ] Backup de BD realizado (si es producción)

**Comando de backup (producción):**
```bash
pg_dump -U postgres -d edugo_prod > backup_pre_sprint01_$(date +%Y%m%d_%H%M%S).sql
```

---

## Troubleshooting

### Error: "función gen_uuid_v7() no existe"
**Solución:** Ejecutar migración 01_base_tables.sql o crear función manualmente (ver sección "Función gen_uuid_v7()")

### Error: "relación 'materials' no existe"
**Solución:** Crear tabla materials (ver sección "Tabla: materials")

### Error: "relación 'users' no existe"
**Solución:** Crear tabla users (ver sección "Tabla: users")

### Error: "permiso denegado para crear tabla"
**Solución:** 
```bash
# Otorgar permisos al usuario
psql -U postgres -d edugo_test -c "GRANT CREATE ON SCHEMA public TO tu_usuario;"
psql -U postgres -d edugo_test -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO tu_usuario;"
```

### Error: "base de datos no existe"
**Solución:**
```bash
# Crear base de datos
psql -U postgres -c "CREATE DATABASE edugo_test;"
```

---

## Dependencias Post-Sprint 01

Archivos generados en Sprint 01 serán usados por:

- **Sprint 02:** Entities Go necesitan coincidir con schema SQL
- **Sprint 03:** Repositorios GORM usarán estas tablas
- **Sprint 05:** Tests de integración usarán seeds

**Próximos Sprints dependen de:**
- Schema PostgreSQL creado en Sprint 01
- Seeds de datos de prueba
- Integridad referencial validada

---

**Generado con:** Claude Code  
**Sprint:** 01/06  
**Última actualización:** 2025-11-14
