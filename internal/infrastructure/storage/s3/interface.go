package s3

import (
	"context"
	"time"
)

// S3Storage define las operaciones de almacenamiento en S3
type S3Storage interface {
	// GeneratePresignedUploadURL genera una URL presignada para subir archivos
	GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error)

	// GeneratePresignedDownloadURL genera una URL presignada para descargar archivos
	GeneratePresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error)
}
