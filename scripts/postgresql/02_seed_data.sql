-- =====================================================
-- Script: 02_seed_data.sql
-- Descripción: Datos mock para testing y desarrollo
-- Autor: EduGo Team
-- Fecha: 2025-11-05
-- Versión: 1.0.0
-- =====================================================

-- Limpiar datos existentes (opcional, solo para desarrollo)
-- TRUNCATE TABLE material_progress, materials, users CASCADE;

-- =====================================================
-- Insertar usuarios de prueba
-- =====================================================

-- ⚠️  DEVELOPMENT ONLY - DO NOT USE IN PRODUCTION
-- All test users use the same password for convenience: password123

-- Profesor 1: Juan Pérez
-- Email: juan.perez@edugo.com
-- Password: password123 (for development only)
INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'juan.perez@edugo.com',
    '$2a$12$IhuIXyaR1EyIIe3Jvdhut.Eb/lSj0n8gQR22jP2v7efg7xw1hubWS',  -- Hash bcrypt de "password123"
    'Juan',
    'Pérez',
    'teacher',
    true,
    NOW() - INTERVAL '60 days',
    NOW() - INTERVAL '60 days'
)
ON CONFLICT (id) DO NOTHING;

-- Profesor 2: María González
-- Email: maria.gonzalez@edugo.com
-- Password: password123 (for development only)
INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
VALUES (
    '22222222-2222-2222-2222-222222222222',
    'maria.gonzalez@edugo.com',
    '$2a$12$IhuIXyaR1EyIIe3Jvdhut.Eb/lSj0n8gQR22jP2v7efg7xw1hubWS',
    'María',
    'González',
    'teacher',
    true,
    NOW() - INTERVAL '45 days',
    NOW() - INTERVAL '45 days'
)
ON CONFLICT (id) DO NOTHING;

-- Estudiante 1: Carlos Rodríguez
-- Email: carlos.rodriguez@student.edugo.com
-- Password: password123 (for development only)
INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
VALUES (
    '33333333-3333-3333-3333-333333333333',
    'carlos.rodriguez@student.edugo.com',
    '$2a$12$IhuIXyaR1EyIIe3Jvdhut.Eb/lSj0n8gQR22jP2v7efg7xw1hubWS',
    'Carlos',
    'Rodríguez',
    'student',
    true,
    NOW() - INTERVAL '30 days',
    NOW() - INTERVAL '30 days'
)
ON CONFLICT (id) DO NOTHING;

-- Estudiante 2: Ana Martínez
-- Email: ana.martinez@student.edugo.com
-- Password: password123 (for development only)
INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
VALUES (
    '44444444-4444-4444-4444-444444444444',
    'ana.martinez@student.edugo.com',
    '$2a$12$IhuIXyaR1EyIIe3Jvdhut.Eb/lSj0n8gQR22jP2v7efg7xw1hubWS',
    'Ana',
    'Martínez',
    'student',
    true,
    NOW() - INTERVAL '25 days',
    NOW() - INTERVAL '25 days'
)
ON CONFLICT (id) DO NOTHING;

-- Estudiante 3: Luis Fernández
-- Email: luis.fernandez@student.edugo.com
-- Password: password123 (for development only)
INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
VALUES (
    '55555555-5555-5555-5555-555555555555',
    'luis.fernandez@student.edugo.com',
    '$2a$12$IhuIXyaR1EyIIe3Jvdhut.Eb/lSj0n8gQR22jP2v7efg7xw1hubWS',
    'Luis',
    'Fernández',
    'student',
    true,
    NOW() - INTERVAL '20 days',
    NOW() - INTERVAL '20 days'
)
ON CONFLICT (id) DO NOTHING;

-- =====================================================
-- Insertar materiales de prueba
-- =====================================================

-- IMPORTANTE: Estos materiales tienen diferentes updated_at para validar
-- el funcionamiento del índice de ordenamiento por fecha

-- Material 1: Introducción a Go (más reciente)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'Introducción a Go',
    'Material básico para aprender los fundamentos del lenguaje Go. Incluye tipos de datos, estructuras de control y funciones.',
    '11111111-1111-1111-1111-111111111111',  -- Juan Pérez
    'programming-101',
    'materials/11111111-1111-1111-1111-111111111111/intro-go.pdf',
    'https://s3.amazonaws.com/edugo-bucket/materials/intro-go.pdf',
    'published',
    'completed',
    false,
    NOW() - INTERVAL '15 days',
    NOW() - INTERVAL '1 day'  -- Actualizado hace 1 día (más reciente)
)
ON CONFLICT (id) DO NOTHING;

-- Material 2: Arquitectura Hexagonal (segundo más reciente)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'Arquitectura Hexagonal en Go',
    'Guía completa sobre Clean Architecture y diseño de software en Go. Patrones, principios SOLID y mejores prácticas.',
    '11111111-1111-1111-1111-111111111111',  -- Juan Pérez
    'architecture-201',
    'materials/11111111-1111-1111-1111-111111111111/clean-arch.pdf',
    'https://s3.amazonaws.com/edugo-bucket/materials/clean-arch.pdf',
    'published',
    'completed',
    false,
    NOW() - INTERVAL '12 days',
    NOW() - INTERVAL '3 days'  -- Actualizado hace 3 días
)
ON CONFLICT (id) DO NOTHING;

-- Material 3: Testing en Go (tercero más reciente)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    'Testing avanzado en Go',
    'Aprende a escribir tests unitarios, de integración y E2E en Go. Incluye mocking, testcontainers y benchmarking.',
    '22222222-2222-2222-2222-222222222222',  -- María González
    'testing-301',
    'materials/22222222-2222-2222-2222-222222222222/go-testing.pdf',
    'https://s3.amazonaws.com/edugo-bucket/materials/go-testing.pdf',
    'published',
    'completed',
    false,
    NOW() - INTERVAL '10 days',
    NOW() - INTERVAL '5 days'  -- Actualizado hace 5 días
)
ON CONFLICT (id) DO NOTHING;

-- Material 4: PostgreSQL Avanzado (cuarto)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    'PostgreSQL: Índices y Optimización',
    'Aprende a optimizar queries, crear índices efectivos y analizar planes de ejecución con EXPLAIN ANALYZE.',
    '22222222-2222-2222-2222-222222222222',  -- María González
    'database-401',
    'materials/22222222-2222-2222-2222-222222222222/postgres-optimization.pdf',
    'https://s3.amazonaws.com/edugo-bucket/materials/postgres-optimization.pdf',
    'published',
    'completed',
    false,
    NOW() - INTERVAL '8 days',
    NOW() - INTERVAL '7 days'  -- Actualizado hace 7 días
)
ON CONFLICT (id) DO NOTHING;

-- Material 5: Docker y Microservicios (quinto)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
    'Docker y Microservicios',
    'Introducción a contenedores Docker, Docker Compose y arquitectura de microservicios con Go.',
    '11111111-1111-1111-1111-111111111111',  -- Juan Pérez
    'devops-501',
    'materials/11111111-1111-1111-1111-111111111111/docker-microservices.pdf',
    'https://s3.amazonaws.com/edugo-bucket/materials/docker-microservices.pdf',
    'published',
    'completed',
    false,
    NOW() - INTERVAL '6 days',
    NOW() - INTERVAL '10 days'  -- Actualizado hace 10 días
)
ON CONFLICT (id) DO NOTHING;

-- Material 6: AWS S3 y Storage (sexto)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    'ffffffff-ffff-ffff-ffff-ffffffffffff',
    'AWS S3: Almacenamiento en la Nube',
    'Aprende a usar AWS S3 para almacenar archivos, configurar permisos y generar URLs firmadas.',
    '22222222-2222-2222-2222-222222222222',  -- María González
    'cloud-601',
    'materials/22222222-2222-2222-2222-222222222222/aws-s3.pdf',
    'https://s3.amazonaws.com/edugo-bucket/materials/aws-s3.pdf',
    'published',
    'completed',
    false,
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '12 days'  -- Actualizado hace 12 días
)
ON CONFLICT (id) DO NOTHING;

-- Material 7: RabbitMQ y Messaging (séptimo)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    '10000000-0000-0000-0000-000000000007',
    'RabbitMQ: Mensajería Asíncrona',
    'Introducción a RabbitMQ, colas de mensajes, exchanges y patrones de mensajería en sistemas distribuidos.',
    '11111111-1111-1111-1111-111111111111',  -- Juan Pérez
    'messaging-701',
    'materials/11111111-1111-1111-1111-111111111111/rabbitmq.pdf',
    'https://s3.amazonaws.com/edugo-bucket/materials/rabbitmq.pdf',
    'published',
    'completed',
    false,
    NOW() - INTERVAL '4 days',
    NOW() - INTERVAL '14 days'  -- Actualizado hace 14 días
)
ON CONFLICT (id) DO NOTHING;

-- Material 8: API REST Design (octavo)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    '10000000-0000-0000-0000-000000000008',
    'Diseño de APIs REST',
    'Mejores prácticas para diseñar APIs REST: verbos HTTP, códigos de estado, versionado y documentación.',
    '22222222-2222-2222-2222-222222222222',  -- María González
    'api-design-801',
    'materials/22222222-2222-2222-2222-222222222222/rest-api.pdf',
    'https://s3.amazonaws.com/edugo-bucket/materials/rest-api.pdf',
    'published',
    'completed',
    false,
    NOW() - INTERVAL '3 days',
    NOW() - INTERVAL '16 days'  -- Actualizado hace 16 días (más antiguo)
)
ON CONFLICT (id) DO NOTHING;

-- Material 9: En proceso (draft para variar)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    '10000000-0000-0000-0000-000000000009',
    'Seguridad en APIs',
    'Material en desarrollo sobre autenticación, autorización, JWT, OAuth 2.0 y mejores prácticas de seguridad.',
    '11111111-1111-1111-1111-111111111111',  -- Juan Pérez
    'security-901',
    NULL,  -- Aún no se ha subido el archivo
    NULL,
    'draft',
    'pending',
    false,
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '2 days'
)
ON CONFLICT (id) DO NOTHING;

-- Material 10: Archivado (para testing de soft delete)
INSERT INTO materials (id, title, description, author_id, subject_id, s3_key, s3_url, status, processing_status, is_deleted, created_at, updated_at)
VALUES (
    '10000000-0000-0000-0000-000000000010',
    'Material Obsoleto - No Usar',
    'Este material está archivado y no debe aparecer en listados públicos.',
    '22222222-2222-2222-2222-222222222222',  -- María González
    'legacy-000',
    'materials/22222222-2222-2222-2222-222222222222/obsolete.pdf',
    'https://s3.amazonaws.com/edugo-bucket/materials/obsolete.pdf',
    'archived',
    'completed',
    false,
    NOW() - INTERVAL '30 days',
    NOW() - INTERVAL '20 days'
)
ON CONFLICT (id) DO NOTHING;

-- =====================================================
-- Insertar progreso de lectura de estudiantes
-- =====================================================

-- Carlos Rodríguez: Ha leído varios materiales
INSERT INTO material_progress (material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at)
VALUES
    -- Completó "Introducción a Go"
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '33333333-3333-3333-3333-333333333333', 100, 50, 'completed', NOW() - INTERVAL '1 day', NOW() - INTERVAL '5 days', NOW() - INTERVAL '1 day'),
    -- En progreso "Arquitectura Hexagonal" (70%)
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '33333333-3333-3333-3333-333333333333', 70, 35, 'in_progress', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 hours'),
    -- Comenzó "Testing en Go" (20%)
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '33333333-3333-3333-3333-333333333333', 20, 10, 'in_progress', NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day')
ON CONFLICT (material_id, user_id) DO NOTHING;

-- Ana Martínez: Estudiante muy activa
INSERT INTO material_progress (material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at)
VALUES
    -- Completó "Introducción a Go"
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '44444444-4444-4444-4444-444444444444', 100, 50, 'completed', NOW() - INTERVAL '2 days', NOW() - INTERVAL '6 days', NOW() - INTERVAL '2 days'),
    -- Completó "Arquitectura Hexagonal"
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '44444444-4444-4444-4444-444444444444', 100, 50, 'completed', NOW() - INTERVAL '3 days', NOW() - INTERVAL '7 days', NOW() - INTERVAL '3 days'),
    -- En progreso "Testing en Go" (85%)
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '44444444-4444-4444-4444-444444444444', 85, 42, 'in_progress', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '4 days', NOW() - INTERVAL '1 hour'),
    -- Comenzó "PostgreSQL Avanzado" (40%)
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', '44444444-4444-4444-4444-444444444444', 40, 20, 'in_progress', NOW() - INTERVAL '6 hours', NOW() - INTERVAL '2 days', NOW() - INTERVAL '6 hours')
ON CONFLICT (material_id, user_id) DO NOTHING;

-- Luis Fernández: Estudiante nuevo, pocos materiales leídos
INSERT INTO material_progress (material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at)
VALUES
    -- En progreso "Introducción a Go" (50%)
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '55555555-5555-5555-5555-555555555555', 50, 25, 'in_progress', NOW() - INTERVAL '3 days', NOW() - INTERVAL '4 days', NOW() - INTERVAL '3 days'),
    -- Comenzó "Docker y Microservicios" (10%)
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '55555555-5555-5555-5555-555555555555', 10, 5, 'in_progress', NOW() - INTERVAL '5 days', NOW() - INTERVAL '6 days', NOW() - INTERVAL '5 days')
ON CONFLICT (material_id, user_id) DO NOTHING;

-- =====================================================
-- Verificación de datos insertados
-- =====================================================

DO $$
DECLARE
    user_count INTEGER;
    material_count INTEGER;
    progress_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO user_count FROM users;
    SELECT COUNT(*) INTO material_count FROM materials WHERE is_deleted = false;
    SELECT COUNT(*) INTO progress_count FROM material_progress;

    RAISE NOTICE '====================================';
    RAISE NOTICE 'Datos de prueba insertados exitosamente';
    RAISE NOTICE '====================================';
    RAISE NOTICE 'Usuarios creados: %', user_count;
    RAISE NOTICE 'Materiales creados: % (% publicados)', material_count, (SELECT COUNT(*) FROM materials WHERE status = 'published');
    RAISE NOTICE 'Registros de progreso: %', progress_count;
    RAISE NOTICE '====================================';
END $$;

-- =====================================================
-- Consultas de verificación rápida
-- =====================================================

-- Listar usuarios por rol
-- SELECT role, COUNT(*) as count FROM users GROUP BY role ORDER BY role;

-- Listar materiales por estado
-- SELECT status, COUNT(*) as count FROM materials WHERE is_deleted = false GROUP BY status ORDER BY status;

-- Ver materiales ordenados por updated_at (validar índice después)
-- SELECT title, updated_at FROM materials WHERE is_deleted = false AND status = 'published' ORDER BY updated_at DESC LIMIT 10;

-- Ver progreso de Carlos Rodríguez
-- SELECT m.title, mp.percentage, mp.status, mp.last_accessed_at FROM material_progress mp JOIN materials m ON mp.material_id = m.id WHERE mp.user_id = '33333333-3333-3333-3333-333333333333' ORDER BY mp.last_accessed_at DESC;

-- =====================================================
-- Resumen de credenciales de usuarios de prueba
-- =====================================================

DO $
DECLARE
    user_record RECORD;
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '====================================';
    RAISE NOTICE 'TEST USER CREDENTIALS (DEVELOPMENT ONLY)';
    RAISE NOTICE '====================================';
    RAISE NOTICE '';
    RAISE NOTICE '⚠️  WARNING: These credentials are for development/testing only!';
    RAISE NOTICE '   DO NOT use these in production environments.';
    RAISE NOTICE '';
    RAISE NOTICE 'All users share the same password: password123';
    RAISE NOTICE '';
    RAISE NOTICE '------------------------------------';
    RAISE NOTICE 'TEACHERS:';
    RAISE NOTICE '------------------------------------';

    FOR user_record IN
        SELECT email, first_name, last_name, role
        FROM users
        WHERE role = 'teacher'
        ORDER BY email
    LOOP
        RAISE NOTICE '  Email:    %', user_record.email;
        RAISE NOTICE '  Password: password123';
        RAISE NOTICE '  Name:     % %', user_record.first_name, user_record.last_name;
        RAISE NOTICE '';
    END LOOP;

    RAISE NOTICE '------------------------------------';
    RAISE NOTICE 'STUDENTS:';
    RAISE NOTICE '------------------------------------';

    FOR user_record IN
        SELECT email, first_name, last_name, role
        FROM users
        WHERE role = 'student'
        ORDER BY email
    LOOP
        RAISE NOTICE '  Email:    %', user_record.email;
        RAISE NOTICE '  Password: password123';
        RAISE NOTICE '  Name:     % %', user_record.first_name, user_record.last_name;
        RAISE NOTICE '';
    END LOOP;

    RAISE NOTICE '====================================';
    RAISE NOTICE 'Quick Login Test:';
    RAISE NOTICE '  curl -X POST http://localhost:8080/api/v1/auth/login \';
    RAISE NOTICE '    -H "Content-Type: application/json" \';
    RAISE NOTICE '    -d ''{"email":"juan.perez@edugo.com","password":"password123"}''';
    RAISE NOTICE '====================================';
    RAISE NOTICE '';
END $;
