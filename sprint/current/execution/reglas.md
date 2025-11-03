# Reglas del Proyecto

Este archivo define las reglas y estándares que el agente de ejecución debe seguir al implementar código. Es opcional pero altamente recomendado para mantener consistencia.

---

## 📝 Estándares de Código

### Convenciones de Nombres

#### Variables y Funciones
```javascript
// ✅ Correcto - camelCase
const userName = 'John';
function getUserById(id) { }

// ❌ Incorrecto
const user_name = 'John';
const UserName = 'John';
```

#### Clases y Componentes
```javascript
// ✅ Correcto - PascalCase
class UserController { }
class AuthService { }

// React Components
const UserProfile = () => { }
```

#### Constantes
```javascript
// ✅ Correcto - UPPER_SNAKE_CASE
const MAX_RETRY_ATTEMPTS = 3;
const API_BASE_URL = 'https://api.example.com';
```

#### Archivos
```
// ✅ Correcto
user-controller.js        (kebab-case)
UserController.js         (PascalCase para clases/componentes)
user.service.ts          (kebab-case)

// ❌ Incorrecto
user_controller.js
UserController.Service.js
```

### Estructura de Proyecto

```
project-root/
├── src/
│   ├── config/          # Configuraciones
│   ├── controllers/     # Controladores (API)
│   ├── models/          # Modelos de datos
│   ├── services/        # Lógica de negocio
│   ├── middleware/      # Middleware de Express/framework
│   ├── routes/          # Definición de rutas
│   ├── utils/           # Utilidades y helpers
│   ├── validators/      # Validadores de entrada
│   └── index.js         # Punto de entrada
├── tests/               # Tests
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── docs/                # Documentación
├── .env.example         # Ejemplo de variables de entorno
├── .gitignore
├── package.json
└── README.md
```

### Comentarios y Documentación

```javascript
// ✅ Correcto - JSDoc para funciones públicas
/**
 * Obtiene un usuario por su ID
 * @param {string} id - ID del usuario
 * @returns {Promise<User>} Usuario encontrado
 * @throws {NotFoundError} Si el usuario no existe
 */
async function getUserById(id) {
  // Comentario inline solo para lógica compleja
  const user = await db.users.findOne({ id });
  if (!user) {
    throw new NotFoundError('Usuario no encontrado');
  }
  return user;
}

// ❌ Evitar comentarios obvios
const name = 'John'; // Asignar nombre a John
```

### Manejo de Errores

```javascript
// ✅ Correcto - Manejo explícito
async function createUser(data) {
  try {
    const validatedData = await validateUserData(data);
    const user = await db.users.create(validatedData);
    logger.info(`Usuario creado: ${user.id}`);
    return user;
  } catch (error) {
    if (error instanceof ValidationError) {
      logger.warn('Datos de usuario inválidos', { error: error.message });
      throw new BadRequestError('Datos inválidos', error.details);
    }
    logger.error('Error creando usuario', { error });
    throw new InternalServerError('Error al crear usuario');
  }
}

// ❌ Incorrecto - Errores silenciados
async function createUser(data) {
  try {
    return await db.users.create(data);
  } catch (error) {
    console.log('Error');  // No es suficiente
  }
}
```

### Organización de Imports

```javascript
// ✅ Correcto - Agrupados y ordenados
// 1. Node.js built-ins
import path from 'path';
import fs from 'fs';

// 2. External dependencies
import express from 'express';
import { body, validationResult } from 'express-validator';

// 3. Internal modules
import { UserService } from './services/user.service.js';
import { authMiddleware } from './middleware/auth.js';
import { logger } from './utils/logger.js';

// ❌ Incorrecto - Desordenado
import { UserService } from './services/user.service.js';
import express from 'express';
import path from 'path';
import { logger } from './utils/logger.js';
```

---

## 🔄 Política de Commits

### Cuándo Hacer Commits

**SÍ hacer commit cuando**:
- ✅ Se completa una tarea atómica del plan
- ✅ El código compila sin errores
- ✅ Los tests existentes pasan
- ✅ La funcionalidad está verificada

**NO hacer commit cuando**:
- ❌ El código tiene errores de compilación
- ❌ Los tests están fallando
- ❌ La implementación está a medias
- ❌ No se ha validado la funcionalidad

### Formato de Mensajes de Commit

Usar **Conventional Commits**:

```
<tipo>(<alcance>): <descripción>

<cuerpo opcional>

<footer opcional>
```

#### Tipos Permitidos:
- `feat`: Nueva funcionalidad
- `fix`: Corrección de bug
- `docs`: Cambios en documentación
- `style`: Cambios de formato (no afectan lógica)
- `refactor`: Refactorización de código
- `test`: Agregar o modificar tests
- `chore`: Tareas de mantenimiento

#### Ejemplos:

```bash
# Nuevo feature
git commit -m "feat(auth): implementar registro de usuarios

- Crear modelo User con validaciones
- Implementar endpoint POST /api/auth/register
- Agregar hash de contraseñas con bcrypt
- Agregar tests unitarios

Completa Tarea 1.1"

# Corrección de bug
git commit -m "fix(auth): corregir validación de email duplicado

Resolver error donde emails duplicados causaban crash
en lugar de retornar error 400

Relacionado con Tarea 1.2"

# Refactorización
git commit -m "refactor(services): extraer lógica de validación a servicio separado"

# Tests
git commit -m "test(auth): agregar tests de integración para endpoints de autenticación"
```

### Estrategia de Commits

- **Commits atómicos**: Un commit por tarea o sub-tarea lógica
- **Commits frecuentes**: Mejor muchos commits pequeños que uno grande
- **Mensajes descriptivos**: Explicar QUÉ y POR QUÉ, no cómo

### Qué NO Incluir en Commits

❌ Archivos de configuración local (`.env`, IDE configs)
❌ Node_modules o dependencias
❌ Archivos de build generados
❌ Logs
❌ Archivos temporales

Asegúrate de tener un `.gitignore` apropiado:
```gitignore
# Dependencies
node_modules/
vendor/

# Environment variables
.env
.env.local

# Build outputs
dist/
build/
*.log

# IDE
.vscode/
.idea/
*.swp

# OS
.DS_Store
Thumbs.db
```

---

## 🧪 Testing Requerido

### Niveles de Testing

#### 1. Tests Unitarios (Obligatorio)
- **Qué**: Funciones/métodos individuales
- **Cuándo**: Para toda lógica de negocio
- **Cobertura mínima**: 80% en servicios y utilidades

```javascript
// Ejemplo con Jest
describe('UserService', () => {
  describe('validateEmail', () => {
    it('debe aceptar emails válidos', () => {
      expect(UserService.validateEmail('test@example.com')).toBe(true);
    });

    it('debe rechazar emails inválidos', () => {
      expect(UserService.validateEmail('invalid-email')).toBe(false);
    });
  });
});
```

#### 2. Tests de Integración (Recomendado)
- **Qué**: Interacción entre módulos
- **Cuándo**: Para flujos críticos (auth, pagos, etc.)

```javascript
// Ejemplo
describe('Auth Integration', () => {
  it('debe registrar usuario y poder hacer login', async () => {
    // Registrar
    const registerRes = await request(app)
      .post('/api/auth/register')
      .send({ email: 'test@example.com', password: 'Test123!' });

    expect(registerRes.status).toBe(201);

    // Login
    const loginRes = await request(app)
      .post('/api/auth/login')
      .send({ email: 'test@example.com', password: 'Test123!' });

    expect(loginRes.status).toBe(200);
    expect(loginRes.body.token).toBeDefined();
  });
});
```

#### 3. Tests E2E (Opcional)
- **Qué**: Flujos completos de usuario
- **Cuándo**: Para funcionalidades críticas del producto

### Ejecución de Tests

```bash
# Ejecutar todos los tests
npm test

# Ejecutar con cobertura
npm run test:coverage

# Ejecutar en modo watch (desarrollo)
npm run test:watch

# Ejecutar solo tests unitarios
npm run test:unit

# Ejecutar solo tests de integración
npm run test:integration
```

### Criterios de Calidad

**Antes de marcar una tarea como completa**:
- ✅ Tests unitarios escritos y pasando
- ✅ Cobertura >= 80% para código nuevo
- ✅ Tests de integración para flujos críticos
- ✅ Sin tests comentados o skipeados (`it.skip`, `xit`)

**Estructura de Tests**:
```
tests/
├── unit/
│   ├── services/
│   │   ├── user.service.test.js
│   │   └── auth.service.test.js
│   └── utils/
│       └── validators.test.js
├── integration/
│   ├── auth.test.js
│   └── users.test.js
└── e2e/
    └── user-journey.test.js
```

---

## 🛠️ Herramientas y Configuración

### Linting
- **ESLint**: Para JavaScript/TypeScript
- **Configuración**: Extender de `eslint:recommended` o Airbnb
- **Ejecución**: Antes de cada commit

```bash
npm run lint
npm run lint:fix  # Auto-corregir
```

### Formateo
- **Prettier**: Para formato consistente
- **Configuración**: `.prettierrc`

```json
{
  "semi": true,
  "trailingComma": "es5",
  "singleQuote": true,
  "printWidth": 100,
  "tabWidth": 2
}
```

### Pre-commit Hooks
- **Husky + lint-staged**: Validar antes de commit

```json
{
  "husky": {
    "hooks": {
      "pre-commit": "lint-staged"
    }
  },
  "lint-staged": {
    "*.js": ["eslint --fix", "prettier --write", "git add"],
    "*.{json,md}": ["prettier --write", "git add"]
  }
}
```

---

## 📋 Checklist antes de Completar una Tarea

Antes de marcar una tarea como completada, verifica:

- [ ] Código implementado según especificación
- [ ] Convenciones de nombres aplicadas
- [ ] Estructura de archivos correcta
- [ ] Comentarios JSDoc en funciones públicas
- [ ] Manejo de errores apropiado
- [ ] Tests unitarios escritos y pasando
- [ ] Tests de integración (si aplica)
- [ ] Código compilado sin errores
- [ ] Linter sin errores
- [ ] Formatter aplicado
- [ ] Sin código comentado innecesariamente
- [ ] Sin console.logs de debug
- [ ] Variables de entorno documentadas en .env.example
- [ ] README actualizado (si aplica)

---

## 🚫 Qué NO Hacer

❌ **No hacer commit sin validar**
❌ **No dejar código comentado en producción**
❌ **No usar console.log para logging** (usar biblioteca de logging)
❌ **No hardcodear valores** (usar configuración)
❌ **No ignorar warnings del linter**
❌ **No dejar TODOs sin contexto**
❌ **No mezclar cambios no relacionados en un commit**
❌ **No copiar/pegar código repetidamente** (extraer a función)

---

## 📝 Notas Adicionales

- Estas reglas pueden adaptarse según las necesidades específicas del sprint
- Si una regla bloquea el progreso razonablemente, documentar la excepción
- Priorizar código funcionando y limpio sobre perfección absoluta
- Cuando tengas duda, preguntar/documentar antes de asumir

---

_Estas reglas son aplicadas por el Agente de Ejecución durante el desarrollo_
