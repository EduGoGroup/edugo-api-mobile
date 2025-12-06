package noop

import (
	"context"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-shared/logger"
	"go.uber.org/zap"
)

// NoopS3Storage es una implementación noop de S3Storage
// Se utiliza cuando S3 no está disponible o está configurado como opcional
type NoopS3Storage struct {
	logger logger.Logger
}

// NewNoopS3Storage crea una nueva instancia de NoopS3Storage
func NewNoopS3Storage(log logger.Logger) *NoopS3Storage {
	return &NoopS3Storage{logger: log}
}

// GeneratePresignedUploadURL simula la generación de una URL presignada para subir archivos
// Registra un mensaje de debug y retorna un error indicando que S3 no está disponible
func (s *NoopS3Storage) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
	s.logger.Debug("noop storage: presigned upload URL not generated (S3 not available)",
		zap.String("key", key),
		zap.String("content_type", contentType),
		zap.Duration("expires", expires),
	)
	return "", fmt.Errorf("s3 not available")
}

// GeneratePresignedDownloadURL simula la generación de una URL presignada para descargar archivos
// Registra un mensaje de debug y retorna un error indicando que S3 no está disponible
func (s *NoopS3Storage) GeneratePresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	s.logger.Debug("noop storage: presigned download URL not generated (S3 not available)",
		zap.String("key", key),
		zap.Duration("expires", expires),
	)
	return "", fmt.Errorf("s3 not available")
}
