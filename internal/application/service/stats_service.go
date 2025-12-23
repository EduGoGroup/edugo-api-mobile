package service

import (
	"context"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

type MaterialStats struct {
	TotalViews    int     `json:"total_views"`
	AvgProgress   float64 `json:"avg_progress"`
	TotalAttempts int     `json:"total_attempts"`
	AvgScore      float64 `json:"avg_score"`
}

type StatsService interface {
	GetMaterialStats(ctx context.Context, materialID string) (*MaterialStats, error)
	GetGlobalStats(ctx context.Context) (*dto.GlobalStatsDTO, error)
}

type statsService struct {
	logger          logger.Logger
	materialStats   repository.MaterialStats   // ISP: Solo necesita estadísticas
	assessmentStats repository.AssessmentStats // ISP: Solo necesita estadísticas
	progressStats   repository.ProgressStats   // ISP: Solo necesita estadísticas
}

func NewStatsService(
	logger logger.Logger,
	materialStats repository.MaterialStats, // ISP: Solo necesita estadísticas
	assessmentStats repository.AssessmentStats, // ISP: Solo necesita estadísticas
	progressStats repository.ProgressStats, // ISP: Solo necesita estadísticas
) StatsService {
	return &statsService{
		logger:          logger,
		materialStats:   materialStats,
		assessmentStats: assessmentStats,
		progressStats:   progressStats,
	}
}

func (s *statsService) GetMaterialStats(ctx context.Context, materialID string) (*MaterialStats, error) {
	_, err := valueobject.MaterialIDFromString(materialID)
	if err != nil {
		return nil, errors.NewValidationError("invalid material_id")
	}

	// Mock stats por ahora (en prod hacer queries a DB)
	stats := &MaterialStats{
		TotalViews:    150,
		AvgProgress:   67.5,
		TotalAttempts: 45,
		AvgScore:      78.3,
	}

	return stats, nil
}

// GetGlobalStats obtiene estadísticas globales del sistema ejecutando queries en paralelo
// Usa goroutines con sync.WaitGroup para optimizar performance
func (s *statsService) GetGlobalStats(ctx context.Context) (*dto.GlobalStatsDTO, error) {
	startTime := time.Now()

	s.logger.Info("iniciando obtención de estadísticas globales")

	// Variables para almacenar resultados de cada query
	var (
		totalMaterials   int64
		totalAssessments int64
		avgScore         float64
		activeUsers      int64
		avgProgress      float64
		queryErrors      []error
		mu               sync.Mutex // Proteger acceso a queryErrors
		wg               sync.WaitGroup
	)

	// Ejecutar 5 queries en paralelo usando goroutines
	wg.Add(5)

	// Goroutine 1: Contar materiales publicados (PostgresSQL)
	go func() {
		defer wg.Done()
		count, err := s.materialStats.CountPublishedMaterials(ctx)
		if err != nil {
			mu.Lock()
			queryErrors = append(queryErrors, err)
			mu.Unlock()
			s.logger.Error("error al contar materiales publicados", "error", err)
			return
		}
		totalMaterials = count
	}()

	// Goroutine 2: Contar evaluaciones completadas (MongoDB)
	go func() {
		defer wg.Done()
		count, err := s.assessmentStats.CountCompletedAssessments(ctx)
		if err != nil {
			mu.Lock()
			queryErrors = append(queryErrors, err)
			mu.Unlock()
			s.logger.Error("error al contar evaluaciones completadas", "error", err)
			return
		}
		totalAssessments = count
	}()

	// Goroutine 3: Calcular promedio de puntajes (MongoDB)
	go func() {
		defer wg.Done()
		avg, err := s.assessmentStats.CalculateAverageScore(ctx)
		if err != nil {
			mu.Lock()
			queryErrors = append(queryErrors, err)
			mu.Unlock()
			s.logger.Error("error al calcular promedio de puntajes", "error", err)
			return
		}
		avgScore = avg
	}()

	// Goroutine 4: Contar usuarios activos (PostgresSQL - últimos 30 días)
	go func() {
		defer wg.Done()
		count, err := s.progressStats.CountActiveUsers(ctx)
		if err != nil {
			mu.Lock()
			queryErrors = append(queryErrors, err)
			mu.Unlock()
			s.logger.Error("error al contar usuarios activos", "error", err)
			return
		}
		activeUsers = count
	}()

	// Goroutine 5: Calcular promedio de progreso (PostgresSQL)
	go func() {
		defer wg.Done()
		avg, err := s.progressStats.CalculateAverageProgress(ctx)
		if err != nil {
			mu.Lock()
			queryErrors = append(queryErrors, err)
			mu.Unlock()
			s.logger.Error("error al calcular promedio de progreso", "error", err)
			return
		}
		avgProgress = avg
	}()

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Sí hubo algún error en las queries, retornar error
	if len(queryErrors) > 0 {
		s.logger.Error("errores al obtener estadísticas globales",
			"total_errors", len(queryErrors))
		return nil, errors.NewInternalError("error al obtener estadísticas del sistema", queryErrors[0])
	}

	// Calcular tiempo de ejecución total
	elapsed := time.Since(startTime).Milliseconds()

	// Construir DTO con resultados
	stats := &dto.GlobalStatsDTO{
		TotalPublishedMaterials:   totalMaterials,
		TotalCompletedAssessments: totalAssessments,
		AverageAssessmentScore:    avgScore,
		ActiveUsersLast30Days:     activeUsers,
		AverageProgress:           avgProgress,
		GeneratedAt:               time.Now(),
	}

	s.logger.Info("estadísticas globales obtenidas exitosamente",
		"total_materials", totalMaterials,
		"total_assessments", totalAssessments,
		"avg_score", avgScore,
		"active_users", activeUsers,
		"avg_progress", avgProgress,
		"elapsed_ms", elapsed)

	return stats, nil
}
