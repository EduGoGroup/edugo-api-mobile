package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// BenchmarkAuthHandler_Login mide el rendimiento del endpoint de login
func BenchmarkAuthHandler_Login(b *testing.B) {
	mockService := &MockAuthService{
		LoginFunc: func(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
			// Simular operación de DB con pequeño delay
			time.Sleep(10 * time.Millisecond)
			return &dto.LoginResponse{
				AccessToken:  "benchmark_access_token",
				RefreshToken: "benchmark_refresh_token",
				User: dto.UserInfo{
					ID:        "user-123",
					Email:     req.Email,
					FirstName: "Benchmark",
					LastName:  "User",
					FullName:  "Benchmark User",
				},
			}, nil
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/login", handler.Login)

	reqBody := `{"email":"benchmark@test.com","password":"password123"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

// BenchmarkAuthHandler_Login_Parallel mide el rendimiento con concurrencia
func BenchmarkAuthHandler_Login_Parallel(b *testing.B) {
	mockService := &MockAuthService{
		LoginFunc: func(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
			time.Sleep(10 * time.Millisecond)
			return &dto.LoginResponse{
				AccessToken:  "benchmark_access_token",
				RefreshToken: "benchmark_refresh_token",
				User: dto.UserInfo{
					ID:        "user-123",
					Email:     req.Email,
					FirstName: "Benchmark",
					LastName:  "User",
					FullName:  "Benchmark User",
				},
			}, nil
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/login", handler.Login)

	reqBody := `{"email":"benchmark@test.com","password":"password123"}`

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Fatalf("Expected status 200, got %d", w.Code)
			}
		}
	})
}

// BenchmarkMaterialHandler_CreateMaterial mide el rendimiento de creación de materiales
func BenchmarkMaterialHandler_CreateMaterial(b *testing.B) {
	mockService := &MockMaterialService{
		CreateMaterialFunc: func(ctx context.Context, req dto.CreateMaterialRequest, authorID string) (*dto.MaterialResponse, error) {
			time.Sleep(15 * time.Millisecond) // Simular DB insert
			return &dto.MaterialResponse{
				ID:    "material-benchmark",
				Title: req.Title,
			}, nil
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.POST("/materials", MockUserIDMiddleware("user-123"), handler.CreateMaterial)

	reqBody := `{"title":"Benchmark Material","description":"Testing performance"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/materials", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			b.Fatalf("Expected status 201, got %d", w.Code)
		}
	}
}

// BenchmarkMaterialHandler_GenerateUploadURL mide el rendimiento de generación de URLs presignadas
func BenchmarkMaterialHandler_GenerateUploadURL(b *testing.B) {
	mockService := &MockMaterialService{
		GetMaterialFunc: func(ctx context.Context, id string) (*dto.MaterialResponse, error) {
			time.Sleep(5 * time.Millisecond) // Simular DB query
			return &dto.MaterialResponse{ID: id, Title: "Benchmark Material"}, nil
		},
	}

	mockS3 := &MockS3Storage{
		GeneratePresignedUploadURLFunc: func(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
			time.Sleep(8 * time.Millisecond) // Simular llamada a AWS SDK
			return "https://s3.amazonaws.com/bucket/" + key + "?presigned", nil
		},
	}

	logger := NewTestLogger()
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.POST("/materials/:id/upload-url", handler.GenerateUploadURL)

	reqBody := `{"file_name":"benchmark.pdf","content_type":"application/pdf"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/materials/test-id/upload-url", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

// BenchmarkMaterialHandler_GenerateUploadURL_Parallel mide rendimiento con concurrencia
func BenchmarkMaterialHandler_GenerateUploadURL_Parallel(b *testing.B) {
	mockService := &MockMaterialService{
		GetMaterialFunc: func(ctx context.Context, id string) (*dto.MaterialResponse, error) {
			time.Sleep(5 * time.Millisecond)
			return &dto.MaterialResponse{ID: id, Title: "Benchmark Material"}, nil
		},
	}

	mockS3 := &MockS3Storage{
		GeneratePresignedUploadURLFunc: func(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
			time.Sleep(8 * time.Millisecond)
			return "https://s3.amazonaws.com/bucket/" + key + "?presigned", nil
		},
	}

	logger := NewTestLogger()
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.POST("/materials/:id/upload-url", handler.GenerateUploadURL)

	reqBody := `{"file_name":"benchmark.pdf","content_type":"application/pdf"}`

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/materials/test-id/upload-url", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Fatalf("Expected status 200, got %d", w.Code)
			}
		}
	})
}

// BenchmarkMaterialHandler_ListMaterials mide el rendimiento de listado
func BenchmarkMaterialHandler_ListMaterials(b *testing.B) {
	// Crear dataset de 50 materiales
	materials := make([]*dto.MaterialResponse, 50)
	for i := 0; i < 50; i++ {
		materials[i] = &dto.MaterialResponse{
			ID:    string(rune(i)),
			Title: "Material " + string(rune(i)),
		}
	}

	mockService := &MockMaterialService{
		ListMaterialsFunc: func(ctx context.Context, filters repository.ListFilters) ([]*dto.MaterialResponse, error) {
			time.Sleep(20 * time.Millisecond) // Simular query compleja
			return materials, nil
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.GET("/materials", handler.ListMaterials)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/materials?page=1&page_size=20", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

// BenchmarkMaterialHandler_GetMaterial mide el rendimiento de obtener un material
func BenchmarkMaterialHandler_GetMaterial(b *testing.B) {
	mockService := &MockMaterialService{
		GetMaterialFunc: func(ctx context.Context, id string) (*dto.MaterialResponse, error) {
			time.Sleep(5 * time.Millisecond)
			return &dto.MaterialResponse{
				ID:    id,
				Title: "Benchmark Material",
			}, nil
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id", handler.GetMaterial)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/materials/test-id", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

// BenchmarkAuthHandler_Refresh mide el rendimiento del refresh de tokens
func BenchmarkAuthHandler_Refresh(b *testing.B) {
	mockService := &MockAuthService{
		RefreshAccessTokenFunc: func(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
			time.Sleep(12 * time.Millisecond)
			return &dto.RefreshResponse{
				AccessToken: "new_access_token",
			}, nil
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/refresh", handler.Refresh)

	reqBody := `{"refresh_token":"valid_refresh_token"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/refresh", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

// BenchmarkHealthHandler_Check_Minimal mide el rendimiento del health check (sin DB)
func BenchmarkHealthHandler_Check_Minimal(b *testing.B) {
	// Este benchmark está simplificado porque HealthHandler requiere conexiones reales
	// Para benchmarks completos, usar testcontainers en ambiente de CI/CD
	b.Skip("Requiere conexiones reales a PostgreSQL y MongoDB")
}

// BenchmarkJSONSerialization mide el overhead de serialización JSON
func BenchmarkJSONSerialization(b *testing.B) {
	response := dto.LoginResponse{
		AccessToken:  "benchmark_access_token_very_long_string_to_simulate_real_jwt",
		RefreshToken: "benchmark_refresh_token_very_long_string_to_simulate_real_jwt",
		User: dto.UserInfo{
			ID:        "user-123",
			Email:     "benchmark@test.com",
			FirstName: "Benchmark",
			LastName:  "User",
			FullName:  "Benchmark User",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(response)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPathTraversalValidation mide el rendimiento de la validación de seguridad
func BenchmarkPathTraversalValidation(b *testing.B) {
	testFiles := []string{
		"document.pdf",
		"../../../etc/passwd",
		"uploads/../../secret.txt",
		"normal-file.pdf",
		"folder/file.pdf",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileName := testFiles[i%len(testFiles)]
		_ = strings.Contains(fileName, "..") ||
			strings.Contains(fileName, "/") ||
			strings.Contains(fileName, "\\")
	}
}

// BenchmarkErrorHandling mide el overhead de manejo de errores
func BenchmarkErrorHandling(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := errors.NewNotFoundError("resource not found")
		_, ok := errors.GetAppError(err)
		if !ok {
			b.Fatal("Expected app error")
		}
	}
}
