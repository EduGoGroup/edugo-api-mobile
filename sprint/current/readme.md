# Sprint: API REST de Saludos (Prueba Rápida)

> **⚡ Ejercicio simple para validar el flujo del sistema**
> Tiempo estimado: 5-10 minutos por fase

---

## Descripción del Proyecto
Crear una API REST minimalista que responda con mensajes de saludo personalizados. Este es un proyecto de prueba para validar que todos los comandos y agentes funcionan correctamente.

## Objetivo del Sprint
Implementar una API REST simple con:
- Endpoint de saludo básico
- Endpoint de saludo personalizado
- Respuestas en formato JSON

---

## Requerimientos Funcionales

### 1. Endpoint de Saludo Básico
- [ ] **GET /api/hello** - Retorna un saludo genérico
  - Respuesta: `{ "message": "Hello, World!" }`

### 2. Endpoint de Saludo Personalizado
- [ ] **GET /api/hello/:name** - Retorna un saludo personalizado
  - Parámetro: `name` (nombre de la persona)
  - Respuesta: `{ "message": "Hello, {name}!" }`
  - Ejemplo: `GET /api/hello/Juan` → `{ "message": "Hello, Juan!" }`

### 3. Endpoint de Estado
- [ ] **GET /api/status** - Retorna el estado del servidor
  - Respuesta: `{ "status": "ok", "timestamp": "2025-10-31T14:30:00Z" }`

---

## Requerimientos No Funcionales

### Simplicidad
- Código minimalista y fácil de entender
- Sin base de datos (todo en memoria)
- Sin autenticación (acceso público)

### Performance
- Respuestas en menos de 100ms
- Servidor debe iniciar en menos de 3 segundos

---

## Stack Tecnológico Sugerido

### Backend
- **Lenguaje**: Node.js
- **Framework**: Express.js
- **Puerto**: 3000

### Estructura Mínima
```
proyecto/
├── src/
│   ├── index.js          # Punto de entrada
│   └── routes/
│       └── hello.js      # Rutas de saludo
├── package.json
└── README.md
```

---

## Casos de Uso

### Caso 1: Saludo Genérico
```bash
# Request
curl http://localhost:3000/api/hello

# Response
{
  "message": "Hello, World!"
}
```

### Caso 2: Saludo Personalizado
```bash
# Request
curl http://localhost:3000/api/hello/Maria

# Response
{
  "message": "Hello, Maria!"
}
```

### Caso 3: Verificar Estado
```bash
# Request
curl http://localhost:3000/api/status

# Response
{
  "status": "ok",
  "timestamp": "2025-10-31T14:30:00Z"
}
```

---

## Entregables Esperados

### Código
- [ ] Servidor Express funcionando
- [ ] Tres endpoints implementados
- [ ] Manejo básico de errores

### Documentación
- [ ] README con instrucciones de ejecución
- [ ] Ejemplos de uso de cada endpoint

### Validación
- [ ] Servidor inicia sin errores
- [ ] Los tres endpoints responden correctamente
- [ ] Mensajes JSON bien formateados

---

## Criterios de Aceptación

El sprint se considerará exitoso cuando:

1. ✅ El servidor Express inicia en el puerto 3000
2. ✅ GET /api/hello retorna `{ "message": "Hello, World!" }`
3. ✅ GET /api/hello/Juan retorna `{ "message": "Hello, Juan!" }`
4. ✅ GET /api/status retorna el estado con timestamp
5. ✅ No hay errores en consola
6. ✅ El código es simple y limpio

---

## Datos de Prueba

### Ejemplos de nombres para probar
- Juan
- Maria
- Carlos
- Ana
- 世界 (caracteres Unicode)

---

## Restricciones

### Simplicidad Máxima
- NO usar base de datos
- NO usar autenticación
- NO usar middleware complejo
- Solo lo esencial para funcionar

### Tiempo
- Este ejercicio debe completarse rápido
- Enfoque en validar el flujo, no en complejidad

---

## Notas Importantes

### Este es un Ejercicio de Validación
- ✅ Sirve para probar que el sistema de comandos/agentes funciona
- ✅ Es rápido de implementar
- ✅ Fácil de validar
- ✅ Demuestra el flujo completo

### NO es el Ejercicio de Presentación
Para la presentación final, usa `Sprint/readme.futuro.md` que tiene un proyecto más completo y visual.

---

## Comandos a Ejecutar

```bash
# 1. Análisis
/01-analisis

# 2. Planificación
/02-planificacion

# 3. Ejecución
/03-ejecucion

# 4. Revisión
/04-revision

# 5. Prueba manual (según guía de validación)
cd proyecto
npm install
npm start
curl http://localhost:3000/api/hello
curl http://localhost:3000/api/hello/Juan
curl http://localhost:3000/api/status
```

---

## Resultado Esperado del Análisis

El agente de análisis debería generar:
- **Arquitectura**: Simple servidor Express con rutas
- **Modelo de datos**: No aplica (sin persistencia)
- **Diagrama de proceso**: Request → Router → Response
- **Resumen**: API REST minimalista

---

## Resultado Esperado de la Planificación

El agente planificador debería generar algo como:

### Fase 1: Configuración Inicial
- Tarea 1.1: Inicializar proyecto npm
- Tarea 1.2: Instalar Express
- Tarea 1.3: Crear estructura de carpetas

### Fase 2: Implementación de Endpoints
- Tarea 2.1: Crear endpoint GET /api/hello
- Tarea 2.2: Crear endpoint GET /api/hello/:name
- Tarea 2.3: Crear endpoint GET /api/status

### Fase 3: Validación
- Tarea 3.1: Probar endpoints manualmente
- Tarea 3.2: Verificar formato JSON
- Tarea 3.3: Documentar en README

---

_Ejercicio de prueba rápida para validar el Sistema de Flujo de Desarrollo_
_Para la presentación final, usar `Sprint/readme.futuro.md`_
