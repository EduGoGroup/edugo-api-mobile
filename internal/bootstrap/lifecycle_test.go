package bootstrap

import (
	"errors"
	"testing"

	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewLifecycleManager verifica que se cree correctamente un nuevo LifecycleManager
func TestNewLifecycleManager(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	assert.NotNil(t, lm, "LifecycleManager no debería ser nil")
	assert.NotNil(t, lm.logger, "Logger no debería ser nil")
	assert.NotNil(t, lm.cleanupFuncs, "cleanupFuncs no debería ser nil")
	assert.Equal(t, 0, len(lm.cleanupFuncs), "cleanupFuncs debería estar vacío inicialmente")
	assert.False(t, lm.cleaned, "cleaned debería ser false inicialmente")
}

// TestLifecycleManager_Register verifica que se registren correctamente las funciones de cleanup
func TestLifecycleManager_Register(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	cleanup := func() error {
		return nil
	}

	lm.Register("test-resource", cleanup)

	assert.Equal(t, 1, len(lm.cleanupFuncs), "Debería haber 1 función de cleanup registrada")
	assert.Equal(t, "test-resource", lm.cleanupFuncs[0].name, "El nombre del recurso debería coincidir")
}

// TestLifecycleManager_RegisterMultiple verifica que se registren múltiples funciones
func TestLifecycleManager_RegisterMultiple(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	lm.Register("resource1", func() error { return nil })
	lm.Register("resource2", func() error { return nil })
	lm.Register("resource3", func() error { return nil })

	assert.Equal(t, 3, len(lm.cleanupFuncs), "Debería haber 3 funciones de cleanup registradas")
	assert.Equal(t, "resource1", lm.cleanupFuncs[0].name)
	assert.Equal(t, "resource2", lm.cleanupFuncs[1].name)
	assert.Equal(t, "resource3", lm.cleanupFuncs[2].name)
}

// TestLifecycleManager_Cleanup_Success verifica que el cleanup se ejecute correctamente
func TestLifecycleManager_Cleanup_Success(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	called := false
	cleanup := func() error {
		called = true
		return nil
	}

	lm.Register("test-resource", cleanup)

	err := lm.Cleanup()

	assert.NoError(t, err, "No debería haber error en cleanup exitoso")
	assert.True(t, called, "La función de cleanup debería haber sido llamada")
	assert.True(t, lm.cleaned, "cleaned debería ser true después del cleanup")
}

// TestLifecycleManager_Cleanup_LIFO verifica que el cleanup se ejecute en orden LIFO
func TestLifecycleManager_Cleanup_LIFO(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	var order []string

	lm.Register("first", func() error {
		order = append(order, "first")
		return nil
	})
	lm.Register("second", func() error {
		order = append(order, "second")
		return nil
	})
	lm.Register("third", func() error {
		order = append(order, "third")
		return nil
	})

	err := lm.Cleanup()

	assert.NoError(t, err)
	require.Equal(t, 3, len(order), "Deberían haberse ejecutado 3 cleanups")
	assert.Equal(t, "third", order[0], "El tercer recurso debería limpiarse primero")
	assert.Equal(t, "second", order[1], "El segundo recurso debería limpiarse segundo")
	assert.Equal(t, "first", order[2], "El primer recurso debería limpiarse último")
}

// TestLifecycleManager_Cleanup_WithError verifica el manejo de errores durante cleanup
func TestLifecycleManager_Cleanup_WithError(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	expectedErr := errors.New("cleanup failed")
	lm.Register("failing-resource", func() error {
		return expectedErr
	})

	err := lm.Cleanup()

	assert.Error(t, err, "Debería haber error cuando un cleanup falla")

	var cleanupErr *CleanupError
	assert.True(t, errors.As(err, &cleanupErr), "El error debería ser de tipo CleanupError")
	assert.Equal(t, 1, cleanupErr.FailureCount, "Debería haber 1 fallo")
	assert.Equal(t, 0, cleanupErr.SuccessCount, "No debería haber éxitos")
}

// TestLifecycleManager_Cleanup_ContinuesOnError verifica que continúe limpiando después de errores
func TestLifecycleManager_Cleanup_ContinuesOnError(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	var cleanedResources []string

	lm.Register("resource1", func() error {
		cleanedResources = append(cleanedResources, "resource1")
		return nil
	})
	lm.Register("resource2", func() error {
		cleanedResources = append(cleanedResources, "resource2")
		return errors.New("resource2 failed")
	})
	lm.Register("resource3", func() error {
		cleanedResources = append(cleanedResources, "resource3")
		return nil
	})

	err := lm.Cleanup()

	assert.Error(t, err, "Debería haber error")

	var cleanupErr *CleanupError
	require.True(t, errors.As(err, &cleanupErr), "El error debería ser de tipo CleanupError")
	assert.Equal(t, 1, cleanupErr.FailureCount, "Debería haber 1 fallo")
	assert.Equal(t, 2, cleanupErr.SuccessCount, "Debería haber 2 éxitos")

	// Verificar que todos los recursos intentaron limpiarse (orden LIFO)
	require.Equal(t, 3, len(cleanedResources), "Deberían haberse intentado limpiar 3 recursos")
	assert.Equal(t, "resource3", cleanedResources[0])
	assert.Equal(t, "resource2", cleanedResources[1])
	assert.Equal(t, "resource1", cleanedResources[2])
}

// TestLifecycleManager_Cleanup_MultipleErrors verifica el manejo de múltiples errores
func TestLifecycleManager_Cleanup_MultipleErrors(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	lm.Register("resource1", func() error {
		return errors.New("error1")
	})
	lm.Register("resource2", func() error {
		return errors.New("error2")
	})
	lm.Register("resource3", func() error {
		return errors.New("error3")
	})

	err := lm.Cleanup()

	assert.Error(t, err)

	var cleanupErr *CleanupError
	require.True(t, errors.As(err, &cleanupErr), "El error debería ser de tipo CleanupError")
	assert.Equal(t, 3, cleanupErr.FailureCount, "Debería haber 3 fallos")
	assert.Equal(t, 0, cleanupErr.SuccessCount, "No debería haber éxitos")
	assert.Equal(t, 3, len(cleanupErr.Errors), "Debería haber 3 errores recolectados")
}

// TestLifecycleManager_Cleanup_Empty verifica el comportamiento con cleanup vacío
func TestLifecycleManager_Cleanup_Empty(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	err := lm.Cleanup()

	assert.NoError(t, err, "No debería haber error con cleanup vacío")
	assert.True(t, lm.cleaned, "cleaned debería ser true")
}

// TestLifecycleManager_Cleanup_Idempotent verifica que Cleanup sea idempotente
func TestLifecycleManager_Cleanup_Idempotent(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	callCount := 0
	lm.Register("resource", func() error {
		callCount++
		return nil
	})

	// Primera llamada
	err1 := lm.Cleanup()
	assert.NoError(t, err1)
	assert.Equal(t, 1, callCount, "Debería haberse llamado una vez")

	// Segunda llamada
	err2 := lm.Cleanup()
	assert.NoError(t, err2)
	assert.Equal(t, 1, callCount, "No debería haberse llamado de nuevo")

	// Tercera llamada
	err3 := lm.Cleanup()
	assert.NoError(t, err3)
	assert.Equal(t, 1, callCount, "No debería haberse llamado de nuevo")
}

// TestLifecycleManager_RegisterAfterCleanup verifica que no se registren funciones después del cleanup
func TestLifecycleManager_RegisterAfterCleanup(t *testing.T) {
	log := logger.NewZapLogger("debug", "json")
	lm := NewLifecycleManager(log)

	lm.Register("resource1", func() error { return nil })

	err := lm.Cleanup()
	assert.NoError(t, err)

	// Intentar registrar después del cleanup
	lm.Register("resource2", func() error { return nil })

	// Debería seguir teniendo solo 1 función registrada
	assert.Equal(t, 1, len(lm.cleanupFuncs), "No debería registrar funciones después del cleanup")
}

// TestCleanupError_Error verifica el mensaje de error de CleanupError
func TestCleanupError_Error(t *testing.T) {
	err := &CleanupError{
		Errors: []error{
			errors.New("error1"),
			errors.New("error2"),
		},
		SuccessCount: 3,
		FailureCount: 2,
	}

	errMsg := err.Error()
	assert.Contains(t, errMsg, "2 failures", "El mensaje debería mencionar los fallos")
	assert.Contains(t, errMsg, "5 resources", "El mensaje debería mencionar el total de recursos")
}

// TestCleanupError_Unwrap verifica que Unwrap retorne los errores subyacentes
func TestCleanupError_Unwrap(t *testing.T) {
	err1 := errors.New("error1")
	err2 := errors.New("error2")

	cleanupErr := &CleanupError{
		Errors: []error{err1, err2},
	}

	unwrapped := cleanupErr.Unwrap()
	assert.Equal(t, 2, len(unwrapped), "Debería retornar 2 errores")
	assert.Equal(t, err1, unwrapped[0])
	assert.Equal(t, err2, unwrapped[1])
}
