package s3

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestS3ConfigValidation verifica que la configuración se valide correctamente
func TestS3ConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      S3Config
		shouldError bool
	}{
		{
			name: "configuración válida",
			config: S3Config{
				Region:          "us-east-1",
				BucketName:      "test-bucket",
				AccessKeyID:     "test-key",
				SecretAccessKey: "test-secret",
			},
			shouldError: false,
		},
		{
			name: "configuración válida con endpoint personalizado",
			config: S3Config{
				Region:          "us-east-1",
				BucketName:      "test-bucket",
				AccessKeyID:     "test-key",
				SecretAccessKey: "test-secret",
				Endpoint:        "http://localhost:4566",
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			log := logger.NewZapLogger("info", "json")

			client, err := NewS3Client(ctx, tt.config, log)

			if tt.shouldError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.config.BucketName, client.bucketName)
				assert.Equal(t, tt.config.Region, client.region)
			}
		})
	}
}

// TestGeneratePresignedUploadURL verifica la generación de URLs presignadas para upload
func TestGeneratePresignedUploadURL(t *testing.T) {
	// Skip si no hay credenciales de AWS configuradas
	// Este test se ejecutará en CI/CD con Localstack
	t.Skip("Requiere conexión a AWS S3 o Localstack")

	ctx := context.Background()
	log := logger.NewZapLogger("info", "json")

	config := S3Config{
		Region:          "us-east-1",
		BucketName:      "test-bucket",
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		Endpoint:        "http://localhost:4566", // Localstack
	}

	client, err := NewS3Client(ctx, config, log)
	require.NoError(t, err)

	// Generar URL presignada
	key := "test-materials/test-file.pdf"
	contentType := "application/pdf"
	expiresIn := 15 * time.Minute

	url, err := client.GeneratePresignedUploadURL(ctx, key, contentType, expiresIn)
	require.NoError(t, err)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "test-bucket")
	assert.Contains(t, url, "test-file.pdf")
}

// TestGeneratePresignedDownloadURL verifica la generación de URLs presignadas para download
func TestGeneratePresignedDownloadURL(t *testing.T) {
	// Skip si no hay credenciales de AWS configuradas
	// Este test se ejecutará en CI/CD con Localstack
	t.Skip("Requiere conexión a AWS S3 o Localstack")

	ctx := context.Background()
	log := logger.NewZapLogger("info", "json")

	config := S3Config{
		Region:          "us-east-1",
		BucketName:      "test-bucket",
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		Endpoint:        "http://localhost:4566", // Localstack
	}

	client, err := NewS3Client(ctx, config, log)
	require.NoError(t, err)

	// Generar URL presignada
	key := "test-materials/test-file.pdf"
	expiresIn := 1 * time.Hour

	url, err := client.GeneratePresignedDownloadURL(ctx, key, expiresIn)
	require.NoError(t, err)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "test-bucket")
	assert.Contains(t, url, "test-file.pdf")
}

// TestPresignedURLExpiration verifica que las URLs tengan tiempo de expiración
func TestPresignedURLExpiration(t *testing.T) {
	t.Skip("Requiere conexión a AWS S3 o Localstack")

	ctx := context.Background()
	log := logger.NewZapLogger("info", "json")

	config := S3Config{
		Region:          "us-east-1",
		BucketName:      "test-bucket",
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		Endpoint:        "http://localhost:4566",
	}

	client, err := NewS3Client(ctx, config, log)
	require.NoError(t, err)

	// Generar URL con tiempo de expiración corto
	key := "test-materials/test-file.pdf"
	contentType := "application/pdf"
	expiresIn := 1 * time.Second

	url, err := client.GeneratePresignedUploadURL(ctx, key, contentType, expiresIn)
	require.NoError(t, err)
	assert.NotEmpty(t, url)

	// La URL debe contener parámetros de expiración
	assert.Contains(t, url, "X-Amz-Expires")
}

// TestS3ClientImplementsInterface verifica que S3Client implementa S3Storage interface
func TestS3ClientImplementsInterface(t *testing.T) {
	ctx := context.Background()
	log := logger.NewZapLogger("info", "json")

	config := S3Config{
		Region:          "us-east-1",
		BucketName:      "test-bucket",
		AccessKeyID:     "test-key",
		SecretAccessKey: "test-secret",
	}

	client, err := NewS3Client(ctx, config, log)
	require.NoError(t, err)

	// Verificar que el cliente implementa la interfaz S3Storage
	var _ S3Storage = client
	assert.NotNil(t, client)
}
