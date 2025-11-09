# Documento de Requisitos - Mejora de Estrategia de Testing

## Introducción

Este documento define los requisitos para analizar, evaluar y mejorar la estrategia de testing del proyecto edugo-api-mobile. El proyecto actualmente tiene una estructura de testing mixta con tests unitarios ubicados junto al código fuente y una carpeta separada para tests de integración. Se requiere un análisis exhaustivo de la estructura actual, identificación de problemas y definición de mejoras para alcanzar una estrategia de testing robusta y escalable.

## Glosario

- **Sistema de Testing**: El conjunto completo de tests unitarios, de integración y herramientas de cobertura del proyecto edugo-api-mobile
- **Test Unitario**: Prueba que verifica el comportamiento de una unidad individual de código (función, método, clase) de forma aislada
- **Test de Integración**: Prueba que verifica la interacción entre múltiples componentes del sistema, incluyendo bases de datos y servicios externos
- **Cobertura de Código**: Métrica que indica el porcentaje de código ejecutado durante los tests
- **Testcontainers**: Biblioteca que permite crear contenedores Docker efímeros para tests de integración
- **Mock**: Objeto simulado que imita el comportamiento de un objeto real para propósitos de testing
- **Seed Data**: Datos de prueba insertados en la base de datos para facilitar los tests
- **PostgreSQL**: Base de datos relacional principal del proyecto
- **MongoDB**: Base de datos NoSQL utilizada para almacenar assessments
- **RabbitMQ**: Sistema de mensajería para comunicación asíncrona
- **Helper de Testing**: Función utilitaria que facilita la configuración y limpieza de tests
- **Schema SQL**: Estructura de tablas y relaciones de la base de datos PostgreSQL
- **Exclusión de Cobertura**: Configuración que permite omitir ciertos archivos o paquetes del cálculo de cobertura

## Requisitos

### Requisito 1: Análisis de Estructura Actual de Testing

**Historia de Usuario:** Como desarrollador del proyecto, quiero entender la estructura actual de testing para identificar fortalezas y debilidades

#### Criterios de Aceptación

1. WHEN el análisis se ejecuta, THE Sistema de Testing SHALL generar un reporte que identifique todos los archivos de test existentes con su ubicación
2. WHEN el análisis se ejecuta, THE Sistema de Testing SHALL calcular la cobertura actual de código por paquete y módulo
3. WHEN el análisis se ejecuta, THE Sistema de Testing SHALL identificar módulos sin tests unitarios
4. WHEN el análisis se ejecuta, THE Sistema de Testing SHALL verificar la existencia y funcionalidad de tests de integración
5. WHEN el análisis se ejecuta, THE Sistema de Testing SHALL documentar la estructura de carpetas test/unit y su contenido actual

### Requisito 2: Evaluación de Calidad de Tests Existentes

**Historia de Usuario:** Como desarrollador del proyecto, quiero evaluar la calidad de los tests existentes para determinar si siguen buenas prácticas

#### Criterios de Aceptación

1. WHEN se evalúan los tests unitarios, THE Sistema de Testing SHALL verificar que utilizan el patrón Arrange-Act-Assert
2. WHEN se evalúan los tests unitarios, THE Sistema de Testing SHALL verificar que utilizan mocks apropiadamente para aislar dependencias
3. WHEN se evalúan los tests de integración, THE Sistema de Testing SHALL verificar que limpian el estado entre ejecuciones
4. WHEN se evalúan los tests de integración, THE Sistema de Testing SHALL verificar que utilizan testcontainers correctamente
5. WHEN se evalúan los tests, THE Sistema de Testing SHALL identificar tests que fallan o están deshabilitados

### Requisito 3: Configuración de Exclusiones de Cobertura

**Historia de Usuario:** Como desarrollador del proyecto, quiero configurar exclusiones de cobertura para omitir código generado y archivos de configuración

#### Criterios de Aceptación

1. WHEN se configura la cobertura, THE Sistema de Testing SHALL excluir archivos generados automáticamente del cálculo de cobertura
2. WHEN se configura la cobertura, THE Sistema de Testing SHALL excluir archivos de configuración y DTOs del cálculo de cobertura
3. WHEN se configura la cobertura, THE Sistema de Testing SHALL excluir el paquete cmd/main del cálculo de cobertura
4. WHEN se configura la cobertura, THE Sistema de Testing SHALL excluir mocks y helpers de testing del cálculo de cobertura
5. WHEN se genera un reporte de cobertura, THE Sistema de Testing SHALL mostrar solo los módulos relevantes según las exclusiones configuradas

### Requisito 4: Organización de Tests Unitarios

**Historia de Usuario:** Como desarrollador del proyecto, quiero una estructura clara para tests unitarios que sea fácil de mantener y escalar

#### Criterios de Aceptación

1. WHEN se organiza la estructura de tests, THE Sistema de Testing SHALL mantener tests unitarios junto al código fuente con sufijo _test.go
2. WHEN se organiza la estructura de tests, THE Sistema de Testing SHALL eliminar carpetas vacías en test/unit que solo contienen archivos .gitkeep
3. WHEN se crea un nuevo módulo, THE Sistema de Testing SHALL incluir tests unitarios en el mismo paquete
4. WHEN se ejecutan tests unitarios, THE Sistema de Testing SHALL excluir tests de integración mediante build tags
5. WHEN se documenta la estructura, THE Sistema de Testing SHALL proporcionar guías claras sobre dónde ubicar tests unitarios

### Requisito 5: Infraestructura de Tests de Integración

**Historia de Usuario:** Como desarrollador del proyecto, quiero una infraestructura robusta para tests de integración que sea reutilizable

#### Criterios de Aceptación

1. WHEN se inicializan tests de integración, THE Sistema de Testing SHALL levantar contenedores de PostgreSQL, MongoDB y RabbitMQ mediante testcontainers
2. WHEN se inicializan tests de integración, THE Sistema de Testing SHALL crear el schema completo de PostgreSQL automáticamente
3. WHEN se inicializan tests de integración, THE Sistema de Testing SHALL crear índices necesarios en MongoDB automáticamente
4. WHEN se inicializan tests de integración, THE Sistema de Testing SHALL proporcionar funciones helper para seed de datos de prueba
5. WHEN se finalizan tests de integración, THE Sistema de Testing SHALL limpiar todos los contenedores y recursos creados

### Requisito 6: Gestión de Datos de Prueba

**Historia de Usuario:** Como desarrollador del proyecto, quiero scripts y funciones para gestionar datos de prueba de forma centralizada

#### Criterios de Aceptación

1. WHEN se crean datos de prueba, THE Sistema de Testing SHALL proporcionar funciones helper para crear usuarios con contraseñas conocidas
2. WHEN se crean datos de prueba, THE Sistema de Testing SHALL proporcionar funciones helper para crear materiales asociados a usuarios
3. WHEN se crean datos de prueba, THE Sistema de Testing SHALL proporcionar funciones helper para crear assessments en MongoDB
4. WHEN se documentan datos de prueba, THE Sistema de Testing SHALL incluir comentarios con valores sin encriptar para contraseñas hasheadas
5. WHEN se limpian datos de prueba, THE Sistema de Testing SHALL proporcionar funciones para truncar todas las tablas en orden correcto

### Requisito 7: Configuración de RabbitMQ para Testing

**Historia de Usuario:** Como desarrollador del proyecto, quiero que RabbitMQ se configure automáticamente para tests y desarrollo

#### Criterios de Aceptación

1. WHEN se inicia RabbitMQ para tests, THE Sistema de Testing SHALL crear exchanges necesarios automáticamente
2. WHEN se inicia RabbitMQ para tests, THE Sistema de Testing SHALL crear colas necesarias automáticamente
3. WHEN se inicia RabbitMQ para tests, THE Sistema de Testing SHALL configurar bindings entre exchanges y colas
4. WHEN RabbitMQ no está disponible, THE Sistema de Testing SHALL utilizar un publisher mock sin fallar los tests
5. WHEN se solicita, THE Sistema de Testing SHALL proporcionar scripts para inicializar RabbitMQ en ambiente de desarrollo

### Requisito 8: Reutilización de Infraestructura para Desarrollo

**Historia de Usuario:** Como desarrollador del proyecto, quiero reutilizar la infraestructura de tests de integración para preparar mi ambiente de desarrollo

#### Criterios de Aceptación

1. WHEN se ejecuta un comando de setup, THE Sistema de Testing SHALL crear contenedores de PostgreSQL y MongoDB para desarrollo local
2. WHEN se ejecuta un comando de setup, THE Sistema de Testing SHALL ejecutar scripts de schema en PostgreSQL
3. WHEN se ejecuta un comando de setup, THE Sistema de Testing SHALL cargar datos de prueba en ambas bases de datos
4. WHEN se ejecuta un comando de setup, THE Sistema de Testing SHALL configurar RabbitMQ con exchanges y colas
5. WHEN se ejecuta un comando de teardown, THE Sistema de Testing SHALL detener y limpiar todos los contenedores de desarrollo

### Requisito 9: Mejora de Cobertura de Tests

**Historia de Usuario:** Como desarrollador del proyecto, quiero un plan claro para incrementar la cobertura de tests del proyecto

#### Criterios de Aceptación

1. WHEN se genera el plan de cobertura, THE Sistema de Testing SHALL identificar módulos críticos con cobertura menor al 50%
2. WHEN se genera el plan de cobertura, THE Sistema de Testing SHALL priorizar módulos de dominio y servicios para testing
3. WHEN se genera el plan de cobertura, THE Sistema de Testing SHALL proporcionar plantillas de tests para módulos sin cobertura
4. WHEN se genera el plan de cobertura, THE Sistema de Testing SHALL establecer metas de cobertura por módulo (mínimo 70% para servicios)
5. WHEN se genera el plan de cobertura, THE Sistema de Testing SHALL documentar casos de prueba faltantes para cada módulo

### Requisito 10: Documentación y Guías de Testing

**Historia de Usuario:** Como desarrollador del proyecto, quiero documentación clara sobre cómo escribir y ejecutar tests

#### Criterios de Aceptación

1. WHEN se consulta la documentación, THE Sistema de Testing SHALL proporcionar guías para escribir tests unitarios con ejemplos
2. WHEN se consulta la documentación, THE Sistema de Testing SHALL proporcionar guías para escribir tests de integración con ejemplos
3. WHEN se consulta la documentación, THE Sistema de Testing SHALL documentar todos los helpers disponibles y su uso
4. WHEN se consulta la documentación, THE Sistema de Testing SHALL incluir comandos make para ejecutar diferentes tipos de tests
5. WHEN se consulta la documentación, THE Sistema de Testing SHALL explicar la estrategia de testing del proyecto y sus principios

### Requisito 11: Validación de Tests Existentes

**Historia de Usuario:** Como desarrollador del proyecto, quiero verificar que todos los tests existentes funcionan correctamente

#### Criterios de Aceptación

1. WHEN se ejecutan tests unitarios, THE Sistema de Testing SHALL completar sin errores todos los tests en internal/application/service
2. WHEN se ejecutan tests unitarios, THE Sistema de Testing SHALL completar sin errores todos los tests en internal/infrastructure/http/handler
3. WHEN se ejecutan tests de integración, THE Sistema de Testing SHALL completar sin errores todos los tests en test/integration
4. WHEN se ejecutan tests de integración, THE Sistema de Testing SHALL verificar que testcontainers se levantan correctamente
5. WHEN se detectan tests fallidos, THE Sistema de Testing SHALL generar un reporte con detalles del fallo y sugerencias de corrección

### Requisito 12: Automatización de Testing en CI/CD

**Historia de Usuario:** Como desarrollador del proyecto, quiero que los tests se ejecuten automáticamente en el pipeline de CI/CD

#### Criterios de Aceptación

1. WHEN se ejecuta el pipeline de CI, THE Sistema de Testing SHALL ejecutar todos los tests unitarios
2. WHEN se ejecuta el pipeline de CI, THE Sistema de Testing SHALL ejecutar tests de integración si Docker está disponible
3. WHEN se ejecuta el pipeline de CI, THE Sistema de Testing SHALL generar reportes de cobertura
4. WHEN se ejecuta el pipeline de CI, THE Sistema de Testing SHALL fallar el build si la cobertura cae por debajo del umbral mínimo
5. WHEN se ejecuta el pipeline de CI, THE Sistema de Testing SHALL publicar reportes de cobertura en formato HTML
