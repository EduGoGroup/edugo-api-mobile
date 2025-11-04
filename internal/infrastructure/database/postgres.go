package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-shared/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// InitPostgreSQL inicializa la conexión a PostgreSQL con la configuración proporcionada.
// Configura el pool de conexiones y verifica que la conexión esté activa.
func InitPostgreSQL(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error) {
	log.Info("inicializando conexión a PostgreSQL",
		zap.String("host", cfg.Database.Postgres.Host),
		zap.Int("port", cfg.Database.Postgres.Port),
		zap.String("database", cfg.Database.Postgres.Database),
	)

	connStr := cfg.Database.Postgres.GetConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("error abriendo conexión a PostgreSQL",
			zap.Error(err),
		)
		return nil, fmt.Errorf("error opening postgres connection: %w", err)
	}

	// Configurar pool de conexiones
	if cfg.Database.Postgres.MaxConnections > 0 {
		db.SetMaxOpenConns(cfg.Database.Postgres.MaxConnections)
		log.Info("pool de conexiones configurado",
			zap.Int("max_open_conns", cfg.Database.Postgres.MaxConnections),
		)
	}
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// Verificar conexión con timeout
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		log.Error("error verificando conexión a PostgreSQL",
			zap.Error(err),
		)
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	log.Info("PostgreSQL conectado exitosamente",
		zap.String("host", cfg.Database.Postgres.Host),
		zap.Int("port", cfg.Database.Postgres.Port),
		zap.String("database", cfg.Database.Postgres.Database),
	)

	return db, nil
}
