package bootstrap

import (
	"fmt"
	"sync"

	"github.com/EduGoGroup/edugo-shared/logger"
	"go.uber.org/zap"
)

// CleanupFunc es una función que cierra un recurso
// Retorna error si el cierre falla, pero no debe hacer panic
type CleanupFunc func() error

// resourceCleanup encapsula una función de cleanup con su nombre
// para propósitos de logging y debugging
type resourceCleanup struct {
	name    string
	cleanup CleanupFunc
}

// LifecycleManager gestiona el ciclo de vida de recursos de infraestructura
// Registra funciones de cleanup y las ejecuta en orden inverso (LIFO)
// cuando se invoca Cleanup()
type LifecycleManager struct {
	cleanupFuncs []resourceCleanup
	logger       logger.Logger
	mu           sync.Mutex // Protege cleanupFuncs de acceso concurrent
	cleaned      bool       // Previene múltiples ejecuciones de Cleanup
}

// NewLifecycleManager crea un nuevo gestor de ciclo de vida
func NewLifecycleManager(log logger.Logger) *LifecycleManager {
	return &LifecycleManager{
		cleanupFuncs: make([]resourceCleanup, 0),
		logger:       log,
		cleaned:      false,
	}
}

// Register registra una función de cleanup para un recurso
// Las funciones se ejecutarán en orden inverso (LIFO) durante Cleanup()
// Es seguro llamar desde múltiples goroutines
func (lm *LifecycleManager) Register(name string, cleanup CleanupFunc) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lm.cleaned {
		lm.logger.Warn("attempted to register cleanup after cleanup was already executed",
			zap.String("resource", name),
		)
		return
	}

	lm.cleanupFuncs = append(lm.cleanupFuncs, resourceCleanup{
		name:    name,
		cleanup: cleanup,
	})

	lm.logger.Debug("registered cleanup function",
		zap.String("resource", name),
		zap.Int("total_cleanups", len(lm.cleanupFuncs)),
	)
}

// Cleanup ejecuta todas las funciones de cleanup registradas en orden inverso (LIFO)
// Continúa ejecutando cleanups incluso si algunos fallan, recolectando todos los errores
// Retorna un error agregado si algún cleanup falló
// Es seguro llamar múltiples veces (solo ejecuta una vez)
func (lm *LifecycleManager) Cleanup() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	// Prevenir múltiples ejecuciones
	if lm.cleaned {
		lm.logger.Debug("cleanup already executed, skipping")
		return nil
	}
	lm.cleaned = true

	if len(lm.cleanupFuncs) == 0 {
		lm.logger.Debug("no cleanup functions registered")
		return nil
	}

	lm.logger.Info("starting resource cleanup",
		zap.Int("resource_count", len(lm.cleanupFuncs)),
	)

	var errors []error
	successCount := 0
	failureCount := 0

	// Ejecutar en orden inverso (LIFO) - último en registrarse, primero en cerrarse
	for i := len(lm.cleanupFuncs) - 1; i >= 0; i-- {
		rc := lm.cleanupFuncs[i]

		lm.logger.Debug("cleaning up resource",
			zap.String("resource", rc.name),
			zap.Int("remaining", i),
		)

		if err := rc.cleanup(); err != nil {
			failureCount++
			errors = append(errors, fmt.Errorf("failed to cleanup %s: %w", rc.name, err))
			lm.logger.Error("cleanup failed for resource",
				zap.String("resource", rc.name),
				zap.Error(err),
			)
		} else {
			successCount++
			lm.logger.Debug("successfully cleaned up resource",
				zap.String("resource", rc.name),
			)
		}
	}

	lm.logger.Info("resource cleanup completed",
		zap.Int("success", successCount),
		zap.Int("failures", failureCount),
		zap.Int("total", len(lm.cleanupFuncs)),
	)

	// Si hubo errores, retornar un error agregado
	if len(errors) > 0 {
		return &CleanupError{
			Errors:       errors,
			SuccessCount: successCount,
			FailureCount: failureCount,
		}
	}

	return nil
}

// CleanupError representa múltiples errores de cleanup
type CleanupError struct {
	Errors       []error
	SuccessCount int
	FailureCount int
}

// Error implementa la interfaz error
func (e *CleanupError) Error() string {
	return fmt.Sprintf("cleanup completed with %d failures out of %d resources: %v",
		e.FailureCount,
		e.SuccessCount+e.FailureCount,
		e.Errors,
	)
}

// Unwrap retorna los errores subyacentes para compatibilidad con errors.Is/As
func (e *CleanupError) Unwrap() []error {
	return e.Errors
}
