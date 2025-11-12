package adapter

import (
	"context"
	"time"

	infraS3 "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
	"github.com/EduGoGroup/edugo-shared/logger"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// StorageClientAdapter adapta *s3.Client (de shared/bootstrap) a s3.S3Storage (interfaz de api-mobile)
// Este adapter envuelve el cliente S3 de AWS SDK v2 y provee funcionalidad de presigned URLs
// que es específica de api-mobile y no está en shared/bootstrap
type StorageClientAdapter struct {
	client        *awsS3.Client
	presignClient *awsS3.PresignClient
	bucketName    string
	logger        logger.Logger
}

// NewStorageClientAdapter crea un nuevo adapter de storage client
// s3Client: cliente S3 retornado por shared/bootstrap
// bucketName: nombre del bucket a usar
// logger: logger para registrar eventos
func NewStorageClientAdapter(
	s3Client *awsS3.Client,
	bucketName string,
	log logger.Logger,
) infraS3.S3Storage {
	return &StorageClientAdapter{
		client:        s3Client,
		presignClient: awsS3.NewPresignClient(s3Client),
		bucketName:    bucketName,
		logger:        log,
	}
}

// GeneratePresignedUploadURL genera una URL presignada para subir archivos a S3
// Implementa la interfaz s3.S3Storage
func (a *StorageClientAdapter) GeneratePresignedUploadURL(
	ctx context.Context,
	key string,
	contentType string,
	expires time.Duration,
) (string, error) {
	// Crear input para PutObject
	input := &awsS3.PutObjectInput{
		Bucket:      &a.bucketName,
		Key:         &key,
		ContentType: &contentType,
	}

	// Generar URL presignada
	request, err := a.presignClient.PresignPutObject(ctx, input, func(opts *awsS3.PresignOptions) {
		opts.Expires = expires
	})

	if err != nil {
		a.logger.Error("failed to generate presigned upload URL",
			"bucket", a.bucketName,
			"key", key,
			"error", err,
		)
		return "", err
	}

	a.logger.Debug("presigned upload URL generated",
		"bucket", a.bucketName,
		"key", key,
		"expires", expires,
	)

	return request.URL, nil
}

// GeneratePresignedDownloadURL genera una URL presignada para descargar archivos de S3
// Implementa la interfaz s3.S3Storage
func (a *StorageClientAdapter) GeneratePresignedDownloadURL(
	ctx context.Context,
	key string,
	expires time.Duration,
) (string, error) {
	// Crear input para GetObject
	input := &awsS3.GetObjectInput{
		Bucket: &a.bucketName,
		Key:    &key,
	}

	// Generar URL presignada
	request, err := a.presignClient.PresignGetObject(ctx, input, func(opts *awsS3.PresignOptions) {
		opts.Expires = expires
	})

	if err != nil {
		a.logger.Error("failed to generate presigned download URL",
			"bucket", a.bucketName,
			"key", key,
			"error", err,
		)
		return "", err
	}

	a.logger.Debug("presigned download URL generated",
		"bucket", a.bucketName,
		"key", key,
		"expires", expires,
	)

	return request.URL, nil
}

// Verificar en compile-time que StorageClientAdapter implementa s3.S3Storage
var _ infraS3.S3Storage = (*StorageClientAdapter)(nil)
