package s3

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

// S3Client maneja las operaciones con AWS S3
type S3Client struct {
	client     *s3.Client
	bucketName string
	region     string
	logger     logger.Logger
}

// S3Config contiene la configuración para el cliente S3
type S3Config struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string // Opcional, para usar con Localstack
}

// NewS3Client crea una nueva instancia del cliente S3
func NewS3Client(ctx context.Context, cfg S3Config, log logger.Logger) (*S3Client, error) {
	// Configurar credenciales estáticas
	credentialsProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		"", // Token de sesión (opcional)
	)

	// Opciones de configuración
	configOptions := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentialsProvider),
	}

	// Si hay un endpoint personalizado (Localstack), agregarlo usando BaseEndpoint
	if cfg.Endpoint != "" {
		configOptions = append(configOptions, config.WithBaseEndpoint(cfg.Endpoint))
	}

	// Cargar configuración de AWS
	awsConfig, err := config.LoadDefaultConfig(ctx, configOptions...)
	if err != nil {
		return nil, errors.NewInternalError("error cargando configuración de AWS", err)
	}

	// Crear cliente S3
	s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		// Forzar path-style para Localstack
		if cfg.Endpoint != "" {
			o.UsePathStyle = true
		}
	})

	log.Info("Cliente S3 inicializado correctamente",
		zap.String("region", cfg.Region),
		zap.String("bucket", cfg.BucketName),
	)

	return &S3Client{
		client:     s3Client,
		bucketName: cfg.BucketName,
		region:     cfg.Region,
		logger:     log,
	}, nil
}

// GeneratePresignedUploadURL genera una URL presignada para subir un archivo a S3
func (c *S3Client) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error) {
	// Crear presigner
	presignClient := s3.NewPresignClient(c.client)

	// Preparar input para PutObject
	putObjectInput := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	// Generar URL presignada
	presignedReq, err := presignClient.PresignPutObject(ctx, putObjectInput, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})
	if err != nil {
		c.logger.Error("error generando URL presignada para upload",
			zap.String("key", key),
			zap.Error(err),
		)
		return "", errors.NewInternalError("error generando URL presignada", err)
	}

	c.logger.Info("URL presignada para upload generada",
		zap.String("key", key),
		zap.Duration("expires_in", expiresIn),
	)

	return presignedReq.URL, nil
}

// GeneratePresignedDownloadURL genera una URL presignada para descargar un archivo de S3
func (c *S3Client) GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	// Crear presigner
	presignClient := s3.NewPresignClient(c.client)

	// Preparar input para GetObject
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	// Generar URL presignada
	presignedReq, err := presignClient.PresignGetObject(ctx, getObjectInput, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})
	if err != nil {
		c.logger.Error("error generando URL presignada para download",
			zap.String("key", key),
			zap.Error(err),
		)
		return "", errors.NewInternalError("error generando URL presignada", err)
	}

	c.logger.Info("URL presignada para download generada",
		zap.String("key", key),
		zap.Duration("expires_in", expiresIn),
	)

	return presignedReq.URL, nil
}

// Compile-time verification that S3Client implements S3Storage interface
var _ S3Storage = (*S3Client)(nil)
