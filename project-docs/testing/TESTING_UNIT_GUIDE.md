# Guía de Tests Unitarios

## 🎯 Objetivo

Los tests unitarios verifican el comportamiento de funciones y métodos individuales de forma aislada, usando mocks para las dependencias.

## 📍 Ubicación

Los tests unitarios se ubican **junto al código fuente** con el sufijo `_test.go`:

```
internal/application/service/
├── auth_service.go
└── auth_service_test.go     ← Tests unitarios aquí
```

## ✅ Patrón AAA

Todos los tests deben seguir el patrón **Arrange-Act-Assert**:

```go
func TestCalculateScore(t *testing.T) {
    // Arrange - Preparar datos
    input := UserResponse{Answer: "A"}
    expected := 100.0
    
    // Act - Ejecutar función
    score := CalculateScore(input)
    
    // Assert - Verificar resultado
    assert.Equal(t, expected, score)
}
```

## 🧩 Uso de Mocks

### Con testify/mock

```go
// Crear mock
mockRepo := new(MockUserRepository)
mockRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

// Usar en service
service := NewAuthService(mockRepo, jwtManager)
result, err := service.GetUser(userID)

// Verificar llamadas
mockRepo.AssertExpectations(t)
mockRepo.AssertCalled(t, "FindByID", mock.Anything, userID)
```

## 📊 Table-Driven Tests

Para múltiples casos de prueba:

```go
func TestEmailValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"válido", "test@example.com", false},
        {"inválido", "invalid", true},
        {"vacío", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            err := ValidateEmail(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 🎯 Ejemplos por Tipo

### Value Objects

```go
func TestNewEmail_Valid(t *testing.T) {
    t.Parallel()
    
    email, err := NewEmail("test@example.com")
    
    require.NoError(t, err)
    assert.Equal(t, "test@example.com", email.String())
    assert.False(t, email.IsZero())
}
```

### Entities

```go
func TestNewMaterial_Validation(t *testing.T) {
    t.Parallel()
    
    material, err := NewMaterial("", "desc", authorID, "")
    
    require.Error(t, err)
    assert.Contains(t, err.Error(), "title is required")
}
```

### Services

```go
func TestAuthService_Login(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    mockJWT := new(MockJWTManager)
    mockRepo.On("FindByEmail", mock.Anything, email).Return(user, nil)
    mockJWT.On("GenerateTokens", user).Return(tokens, nil)
    
    service := NewAuthService(mockRepo, mockJWT)
    
    // Act
    result, err := service.Login(ctx, email, password)
    
    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, result.AccessToken)
    mockRepo.AssertExpectations(t)
}
```

## ⚡ Comandos

```bash
# Ejecutar tests unitarios
make test-unit

# Con cobertura
make test-unit-coverage

# Watch mode (desarrollo)
make test-watch
```

## 📚 Librerías Recomendadas

- `github.com/stretchr/testify/assert` - Assertions
- `github.com/stretchr/testify/require` - Assertions que detienen el test
- `github.com/stretchr/testify/mock` - Mocking

---

**Ver también**:
- [TESTING_INTEGRATION_GUIDE.md](./TESTING_INTEGRATION_GUIDE.md)
- [TESTING_GUIDE.md](./TESTING_GUIDE.md)
