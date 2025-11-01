# Roadmap: Social Login para EduGo

**Contexto**: Usuario confirmó que querrá social login en el futuro
**Estrategia**: Implementar auth básico ahora, agregar social login gradualmente

---

## 🎯 Fases de Implementación

### **FASE 1: Autenticación Básica (AHORA - 1 semana)**

```
Implementar:
✅ Email/Password con bcrypt
✅ Refresh tokens
✅ JWT estándar
✅ Logout/Revocación
✅ Rate limiting

Usuarios pueden:
- Registrarse con email/password
- Login
- Logout
- Cambiar contraseña

Esfuerzo: 5-7 días
Costo: $0
```

---

### **FASE 2: Preparar Infraestructura para Social (2-3 meses después)**

```
Refactorizar:
✅ Separar "creación de usuario" de "autenticación"
✅ Tabla users: agregar campos
   - auth_provider (email, google, facebook, apple)
   - provider_id (ID del usuario en el provider)
   - email_verified (boolean)

Preparar:
✅ AuthService.findOrCreateUserByEmail()
✅ AuthService.linkSocialAccount()

Esfuerzo: 1-2 días
Costo: $0
```

---

### **FASE 3: Implementar Social Login (Cuando sea necesario)**

#### **3.1 Google Sign-In** (El más común)

**En apps móviles**:
```kotlin
// Android
GoogleSignIn.getClient(this, GoogleSignInOptions.DEFAULT_SIGN_IN)
    .signInIntent
    .let { startActivityForResult(it, RC_SIGN_IN) }

// Obtienes: idToken de Google
```

```swift
// iOS
GIDSignIn.sharedInstance.signIn(with: config, presenting: self) { user, error in
    let idToken = user?.authentication.idToken
    // Envías a tu API
}
```

**En tu API**:
```go
// POST /v1/auth/login/google
func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
    // 1. Recibir token de Google
    var req struct {
        GoogleToken string `json:"google_token"`
    }
    c.BindJSON(&req)

    // 2. Verificar token con Google
    payload, err := verifyGoogleToken(req.GoogleToken)
    if err != nil {
        return unauthorized
    }

    // 3. Extraer datos
    email := payload.Claims["email"]
    name := payload.Claims["name"]
    googleID := payload.Subject

    // 4. Buscar o crear usuario
    user, err := h.authService.FindOrCreateUserFromGoogle(email, name, googleID)

    // 5. Generar TU token JWT (igual que login normal)
    token, err := h.jwtManager.GenerateToken(user.ID, user.Email, user.Role)

    // 6. Retornar
    c.JSON(200, LoginResponse{Token: token})
}

func verifyGoogleToken(token string) (*oauth2.Payload, error) {
    // Llamar a Google API para verificar
    // https://oauth2.googleapis.com/tokeninfo?id_token={token}
    // O usar librería: google.golang.org/api/oauth2/v2
}
```

**Esfuerzo**: 2-3 días por provider
**Costo**: $0 (Google Sign-In es gratis)

---

#### **3.2 Apple Sign-In** (Requerido para iOS)

**Nota importante**: Apple REQUIERE ofrecer "Sign in with Apple" si ofreces otros social logins en iOS.

**Implementación similar a Google**:
```go
// POST /v1/auth/login/apple
func (h *AuthHandler) LoginWithApple(c *gin.Context) {
    // Verificar token de Apple
    // Proceso similar a Google
}
```

**Esfuerzo**: 2 días
**Costo**: $0

---

#### **3.3 Facebook Login** (Opcional)

**Esfuerzo**: 2 días
**Costo**: $0

---

## 🔄 Flujo Completo de Autenticación (Con Social Login)

### **Opciones de Login para Usuarios**:

```
┌─────────────────────────────────────────┐
│      Pantalla de Login                  │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  Email    ____________________    │ │
│  │  Password ____________________    │ │
│  │           [  Login  ]             │ │
│  └───────────────────────────────────┘ │
│                                         │
│         ─── o continuar con ───         │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  [🔵 Continuar con Google]      │   │
│  └─────────────────────────────────┘   │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  [ Continuar con Apple]        │   │
│  └─────────────────────────────────┘   │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  [🔵 Continuar con Facebook]    │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

---

## 💾 Cambios en Base de Datos para Social Login

### **Tabla `user` actualizada**:

```sql
ALTER TABLE "user"
ADD COLUMN auth_provider VARCHAR(20) DEFAULT 'email',
ADD COLUMN provider_id VARCHAR(255),
ADD COLUMN email_verified BOOLEAN DEFAULT FALSE,
ADD COLUMN profile_picture_url TEXT;

-- Índice para búsqueda rápida
CREATE INDEX idx_provider_id ON "user"(auth_provider, provider_id);

-- Constraint: email puede ser null para algunos providers
ALTER TABLE "user"
ALTER COLUMN email DROP NOT NULL;
-- Apple puede ocultar email, entonces provider_id es clave
```

### **Posibles valores de `auth_provider`**:

| Valor | Descripción |
|-------|-------------|
| `email` | Registro tradicional con email/password |
| `google` | Google Sign-In |
| `apple` | Apple Sign-In |
| `facebook` | Facebook Login |

### **Ejemplo de registros**:

```sql
-- Usuario tradicional
email: jhoan@edugo.com
password_hash: $2a$12$...
auth_provider: email
provider_id: NULL

-- Usuario con Google
email: jhoan@gmail.com
password_hash: NULL (no tiene password)
auth_provider: google
provider_id: 117482651234567890123 (Google user ID)
email_verified: TRUE (Google ya lo verificó)

-- Usuario con Apple (email oculto)
email: privaterelay@icloud.com (email privado de Apple)
password_hash: NULL
auth_provider: apple
provider_id: 001234.a1b2c3d4e5f6... (Apple user ID)
```

---

## 🔐 Consideraciones de Seguridad

### **1. Vincular Cuentas (Account Linking)**

**Problema**: Usuario se registró con email, luego intenta login con Google usando mismo email.

**Solución**:
```go
func (s *AuthService) LoginWithGoogle(googleEmail, googleID string) (*User, error) {
    // Buscar por Google ID primero
    user, err := s.userRepo.FindByProviderID("google", googleID)
    if err == nil {
        return user, nil // Ya existe, login normal
    }

    // No existe por Google ID, buscar por email
    user, err = s.userRepo.FindByEmail(googleEmail)
    if err == nil {
        // Usuario ya existe con email/password
        // Opción 1: Error (pedir que use password)
        // Opción 2: Vincular automáticamente (riesgoso)
        // Opción 3: Pedir confirmación (recomendado)

        return nil, errors.New("email_already_registered_with_password")
    }

    // No existe, crear nuevo usuario
    return s.createUserFromGoogle(googleEmail, googleID)
}
```

**UX Recomendado**:
```
Usuario intenta: Login with Google (jhoan@edugo.com)
Sistema detecta: Email ya registrado con password

Mostrar:
"Ya tienes una cuenta con este email. ¿Quieres vincularla con Google?
[ Sí, vincular ]  [ No, usar password ]"

Si elige "Sí, vincular":
1. Pedir password actual (verificar identidad)
2. Vincular provider_id de Google
3. Ahora puede usar ambos métodos
```

---

### **2. Verificación de Email**

**Con email/password**:
```
1. Usuario se registra
2. Envías email de verificación
3. Usuario hace click en link
4. email_verified = TRUE
```

**Con Google/Facebook**:
```
1. Usuario hace login con Google
2. Google ya verificó el email
3. email_verified = TRUE (automático)
```

**Con Apple**:
```
1. Usuario hace login con Apple
2. Apple puede ocultar email real
3. email_verified = TRUE pero email puede ser relay
```

---

## 📱 SDKs Necesarios

### **Para Android**:

```gradle
// build.gradle
dependencies {
    // Google Sign-In
    implementation 'com.google.android.gms:play-services-auth:20.7.0'

    // Facebook Login (opcional)
    implementation 'com.facebook.android:facebook-login:16.1.3'
}
```

### **Para iOS**:

```swift
// Podfile
pod 'GoogleSignIn', '~> 7.0'
pod 'FBSDKLoginKit' // Facebook (opcional)
```

**Apple Sign-In**: Incluido en iOS, no requiere dependencias.

---

## ⏱️ Estimación de Esfuerzo

| Tarea | Esfuerzo | Cuándo |
|-------|----------|--------|
| **Auth básico** | 5-7 días | AHORA |
| **Preparar DB para social** | 1 día | 2-3 meses |
| **Google Sign-In** | 2-3 días | Cuando usuarios lo pidan |
| **Apple Sign-In** | 2 días | Cuando Google esté listo |
| **Facebook** | 2 días | Opcional |
| **Account linking** | 1-2 días | Con primer social provider |
| **TOTAL** | **13-17 días** | **Distribuido en 6-12 meses** |

---

## 💰 Costos

| Provider | Costo | Límites |
|----------|-------|---------|
| **Google Sign-In** | $0 | Ilimitado |
| **Apple Sign-In** | $0 | Ilimitado |
| **Facebook Login** | $0 | Ilimitado |

**Todos los social logins son GRATIS**, solo pagas desarrollo/implementación.

---

## 🎯 Estrategia Recomendada para EduGo

### **Fase 1 (Ahora - Mes 1)**:
```
✅ Implementar email/password con bcrypt
✅ Refresh tokens
✅ JWT + middleware
✅ Endpoints básicos

Resultado: Usuarios pueden registrarse y usar la app
```

### **Fase 2 (Mes 3-4)**:
```
✅ Agregar campos a tabla users (auth_provider, etc.)
✅ Refactorizar servicios para soportar múltiples providers
✅ Preparar infraestructura

Resultado: Listo para agregar social login cuando se necesite
```

### **Fase 3 (Mes 6-12)**:
```
✅ Implementar Google Sign-In (el más demandado)
✅ Implementar Apple Sign-In (requerido por App Store)
✅ Account linking
✅ Opcional: Facebook

Resultado: Múltiples opciones de login para usuarios
```

---

## 🚦 Criterios para Implementar Social Login

**Implementar cuando**:
- ✅ Usuarios lo pidan frecuentemente (>10 solicitudes)
- ✅ Tasa de registro sea baja (<30%)
- ✅ Competidores lo ofrezcan
- ✅ Tengas tiempo de desarrollo (2-3 días)

**NO es urgente si**:
- ❌ Usuarios están OK con email/password
- ❌ Tasa de registro es buena
- ❌ Tienes otras prioridades (features de negocio)

---

## 📝 Conclusión

**Para tu pregunta**: "Si en su momento querré social login"

**Respuesta**:
✅ Perfecto, la implementación propia que planteé es **100% compatible** con agregar social login después.

**Plan**:
1. **Ahora**: Implementar auth básico (email/password) en 1 semana
2. **2-3 meses**: Preparar DB (1 día)
3. **6-12 meses**: Agregar Google + Apple (4-5 días cuando lo necesites)

**Ventajas de este approach**:
- ✅ Empiezas rápido con lo básico
- ✅ No pagas $350/mes de Auth0 innecesariamente
- ✅ Agregas social login solo cuando usuarios lo pidan
- ✅ Sin vendor lock-in
- ✅ Código modular (agregar providers es fácil)

**¿Procedo con Fase 1 (auth básico)?**
En 5-7 días tendrás autenticación robusta, y social login lo agregas cuando lo necesites.

---

**Última actualización**: 2025-10-31
**Próxima revisión**: Cuando se implemente Fase 1
