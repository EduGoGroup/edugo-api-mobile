# ============================================
# Makefile - API Mobile (EduGo)
# ============================================

# Variables
APP_NAME=api-mobile
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=bin
COVERAGE_DIR=coverage
MAIN_PATH=./cmd/main.go

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"

# Colors
YELLOW=\033[1;33m
GREEN=\033[1;32m
BLUE=\033[1;34m
RED=\033[1;31m
RESET=\033[0m

# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

.DEFAULT_GOAL := help

# ============================================
# Main Targets
# ============================================

help: ## Mostrar esta ayuda
	@echo "$(BLUE)$(APP_NAME) - Comandos disponibles:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'

build: ## Compilar binario
	@echo "$(YELLOW)🔨 Compilando $(APP_NAME)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Binario: $(BUILD_DIR)/$(APP_NAME) ($(VERSION))$(RESET)"

build-debug: ## Compilar binario para debugging (sin optimizaciones, útil para Delve)
	@echo "$(YELLOW)🔨 Compilando $(APP_NAME) para debug (all=-N -l)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -gcflags "all=-N -l" $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Binario para debug: $(BUILD_DIR)/$(APP_NAME)$(RESET)"

run: ## Ejecutar en modo desarrollo
	@echo "$(YELLOW)🚀 Ejecutando $(APP_NAME) (ambiente: $(APP_ENV))...$(RESET)"
	@$(GOCMD) run $(MAIN_PATH)

dev: deps swagger run ## Desarrollo completo

# ============================================
# Testing
# ============================================

test: ## Ejecutar todos los tests
	@echo "$(YELLOW)🧪 Ejecutando tests...$(RESET)"
	@$(GOTEST) -v -race ./...
	@echo "$(GREEN)✓ Tests completados$(RESET)"

test-coverage: ## Tests con cobertura (HTML report)
	@echo "$(YELLOW)📊 Generando reporte de cobertura...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	@$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ Reporte: $(COVERAGE_DIR)/coverage.html$(RESET)"
	@echo "$(BLUE)💡 Abrir: open $(COVERAGE_DIR)/coverage.html$(RESET)"

test-unit: ## Solo tests unitarios (rápido)
	@echo "$(YELLOW)🧪 Ejecutando tests unitarios...$(RESET)"
	@$(GOTEST) -v -short -race ./internal/... -timeout 2m
	@echo "$(GREEN)✓ Tests unitarios completados$(RESET)"

test-unit-coverage: ## Tests unitarios con cobertura
	@echo "$(YELLOW)📊 Tests unitarios con cobertura...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -v -short -race -coverprofile=$(COVERAGE_DIR)/unit-coverage.out ./internal/... -timeout 2m
	@./scripts/filter-coverage.sh $(COVERAGE_DIR)/unit-coverage.out $(COVERAGE_DIR)/unit-coverage-filtered.out
	@$(GOCMD) tool cover -html=$(COVERAGE_DIR)/unit-coverage-filtered.out -o $(COVERAGE_DIR)/unit-coverage.html
	@echo "$(GREEN)✓ Reporte: $(COVERAGE_DIR)/unit-coverage.html$(RESET)"

test-integration: ## Tests de integración (con testcontainers) - HABILITADOS con RUN_INTEGRATION_TESTS=true
	@echo "$(YELLOW)🐳 Ejecutando tests de integración con testcontainers...$(RESET)"
	@RUN_INTEGRATION_TESTS=true $(GOTEST) -v -tags=integration ./test/integration/... -timeout 10m
	@echo "$(GREEN)✓ Tests de integración completados$(RESET)"

test-integration-verbose: ## Tests de integración con logs detallados
	@echo "$(YELLOW)🐳 Tests de integración (verbose)...$(RESET)"
	@RUN_INTEGRATION_TESTS=true $(GOTEST) -v -tags=integration ./test/integration/... -timeout 10m

test-integration-skip: ## Tests de integración DESHABILITADOS (skip automático)
	@echo "$(BLUE)⏭️  Tests de integración deshabilitados$(RESET)"
	@RUN_INTEGRATION_TESTS=false $(GOTEST) -v -tags=integration ./test/integration/... -timeout 5m
	@echo "$(BLUE)ℹ️  Tests skipped (esperado)$(RESET)"

test-integration-coverage: ## Tests de integración con coverage
	@echo "$(YELLOW)📊 Tests de integración con coverage...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@RUN_INTEGRATION_TESTS=true $(GOTEST) -tags=integration -coverprofile=$(COVERAGE_DIR)/integration-coverage.out -covermode=atomic ./test/integration/... -timeout 10m
	@$(GOCMD) tool cover -html=$(COVERAGE_DIR)/integration-coverage.out -o $(COVERAGE_DIR)/integration-coverage.html
	@echo "$(GREEN)✓ Reporte: $(COVERAGE_DIR)/integration-coverage.html$(RESET)"

test-all: test-unit test-integration ## Ejecutar todos los tests

test-watch: ## Watch mode para tests (requiere entr)
	@echo "$(YELLOW)👀 Watching tests...$(RESET)"
	@command -v entr > /dev/null || (echo "$(RED)❌ entr no instalado. Instalar con: brew install entr$(RESET)" && exit 1)
	@find . -name "*.go" | entr -c make test-unit

coverage-report: ## Generar reporte de cobertura completo con filtrado
	@echo "$(YELLOW)📊 Generando reporte de cobertura completo...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out ./... -timeout 5m
	@./scripts/filter-coverage.sh $(COVERAGE_DIR)/coverage.out $(COVERAGE_DIR)/coverage-filtered.out
	@$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage-filtered.out -o $(COVERAGE_DIR)/coverage.html
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage-filtered.out | tail -20
	@echo "$(GREEN)✓ Reporte: $(COVERAGE_DIR)/coverage.html$(RESET)"

coverage-check: ## Verificar que cobertura cumple umbral mínimo (60%)
	@echo "$(YELLOW)🎯 Verificando cobertura mínima...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out ./... -timeout 5m
	@./scripts/filter-coverage.sh $(COVERAGE_DIR)/coverage.out $(COVERAGE_DIR)/coverage-filtered.out
	@./scripts/check-coverage.sh $(COVERAGE_DIR)/coverage-filtered.out 60
	@echo "$(GREEN)✓ Cobertura cumple umbral mínimo$(RESET)"

docker-check: ## Verificar que Docker esté corriendo
	@echo "$(YELLOW)🐳 Verificando Docker...$(RESET)"
	@docker ps > /dev/null 2>&1 || (echo "$(RED)❌ Docker no está corriendo. Inicia Docker Desktop.$(RESET)" && exit 1)
	@echo "$(GREEN)✓ Docker está corriendo$(RESET)"

benchmark: ## Ejecutar benchmarks
	@echo "$(YELLOW)⚡ Ejecutando benchmarks...$(RESET)"
	@$(GOTEST) -bench=. -benchmem ./...

# ============================================
# Code Quality
# ============================================

fmt: ## Formatear código
	@echo "$(YELLOW)✨ Formateando código...$(RESET)"
	@$(GOFMT) -w .
	@echo "$(GREEN)✓ Código formateado$(RESET)"

vet: ## Análisis estático
	@echo "$(YELLOW)🔍 Ejecutando go vet...$(RESET)"
	@$(GOVET) ./...
	@echo "$(GREEN)✓ Análisis estático completado$(RESET)"

lint: ## Linter completo
	@echo "$(YELLOW)🔎 Ejecutando golangci-lint...$(RESET)"
	@golangci-lint run --timeout=5m || echo "$(YELLOW)⚠️  Instalar con: make tools$(RESET)"

audit: ## Auditoría de calidad completa
	@echo "$(BLUE)=== 🔐 AUDITORÍA ===$(RESET)"
	@echo "$(YELLOW)1. Verificando go.mod...$(RESET)"
	@$(GOMOD) verify
	@echo "$(YELLOW)2. Formato...$(RESET)"
	@test -z "$$($(GOFMT) -l .)" || (echo "$(RED)Sin formatear:$(RESET)" && $(GOFMT) -l .)
	@echo "$(YELLOW)3. Vet...$(RESET)"
	@$(GOVET) ./...
	@echo "$(YELLOW)4. Tests...$(RESET)"
	@$(GOTEST) -race -vet=off ./...
	@echo "$(GREEN)✓ Auditoría completada$(RESET)"

# ============================================
# Dependencies
# ============================================

deps: ## Descargar dependencias
	@echo "$(YELLOW)📦 Instalando dependencias...$(RESET)"
	@$(GOMOD) download
	@echo "$(GREEN)✓ Dependencias listas$(RESET)"

tidy: ## Limpiar go.mod
	@echo "$(YELLOW)🧹 Limpiando go.mod...$(RESET)"
	@$(GOMOD) tidy
	@echo "$(GREEN)✓ go.mod actualizado$(RESET)"

tools: ## Instalar herramientas
	@echo "$(YELLOW)🛠️  Instalando herramientas...$(RESET)"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Herramientas instaladas$(RESET)"

configctl: ## Build configctl CLI
	@echo "$(YELLOW)🔧 Building configctl...$(RESET)"
	@go build -o bin/configctl ./tools/configctl/
	@echo "$(GREEN)✓ configctl: bin/configctl$(RESET)"

config-validate: configctl ## Validar archivos de configuración
	@./bin/configctl validate

config-docs: configctl ## Generar documentación de configuración
	@./bin/configctl generate-docs

# ============================================
# Swagger
# ============================================

swagger: ## Regenerar Swagger
	@echo "$(YELLOW)📚 Regenerando Swagger...$(RESET)"
	@swag init -g cmd/main.go -o docs --parseInternal
	@echo "$(GREEN)✓ Swagger: http://localhost:8080/swagger/index.html$(RESET)"

# ============================================
# Desarrollo Local
# ============================================

dev-setup: docker-check ## Configurar ambiente de desarrollo con Docker
	@echo "$(YELLOW)🚀 Configurando ambiente de desarrollo...$(RESET)"
	@./test/scripts/setup_dev_env.sh

dev-teardown: ## Limpiar ambiente de desarrollo
	@echo "$(YELLOW)🧹 Limpiando ambiente de desarrollo...$(RESET)"
	@./test/scripts/teardown_dev_env.sh

dev-reset: dev-teardown dev-setup ## Resetear ambiente de desarrollo

dev-logs: ## Ver logs de contenedores de desarrollo
	@docker-compose -f docker-compose-dev.yml logs -f

dev-status: ## Ver estado de contenedores de desarrollo
	@echo "$(BLUE)📊 Estado de contenedores de desarrollo:$(RESET)"
	@docker-compose -f docker-compose-dev.yml ps

# ============================================
# Análisis de Tests
# ============================================

test-analyze: ## Analizar estructura de tests
	@echo "$(YELLOW)🔍 Analizando estructura de tests...$(RESET)"
	@echo "$(BLUE)Tests Unitarios:$(RESET)"
	@find internal -name "*_test.go" -type f | wc -l | xargs echo "  Archivos:"
	@echo "$(BLUE)Tests de Integración:$(RESET)"
	@find test/integration -name "*_test.go" -type f | wc -l | xargs echo "  Archivos:"

test-missing: ## Identificar módulos sin tests
	@echo "$(YELLOW)🔍 Buscando módulos sin tests...$(RESET)"
	@echo "$(RED)Módulos sin cobertura (0%):$(RESET)"
	@go test -coverprofile=/tmp/coverage-check.out ./... > /dev/null 2>&1 || true
	@go tool cover -func=/tmp/coverage-check.out | grep "0.0%" | head -20 || echo "$(GREEN)  Todos los módulos tienen alguna cobertura$(RESET)"

test-validate: test-unit test-integration ## Validar que todos los tests pasan
	@echo "$(GREEN)✅ Todos los tests pasan$(RESET)"

# ============================================
# Docker
# ============================================

docker-build: ## Build imagen
	@echo "$(YELLOW)🐳 Building...$(RESET)"
	@docker build -t edugo/$(APP_NAME):$(VERSION) .
	@echo "$(GREEN)✓ Imagen: edugo/$(APP_NAME):$(VERSION)$(RESET)"

docker-run: ## Run con compose
	@docker-compose up -d
	@echo "$(GREEN)✓ Corriendo en http://localhost:8080$(RESET)"

docker-stop: ## Stop compose
	@docker-compose down

docker-logs: ## Ver logs
	@docker-compose logs -f

# ============================================
# CI/CD
# ============================================

ci: audit test-coverage swagger ## CI pipeline
	@echo "$(GREEN)✅ CI completado$(RESET)"

pre-commit: fmt vet test ## Pre-commit hook

# ============================================
# Cleanup
# ============================================

clean: ## Limpiar todo
	@rm -rf $(BUILD_DIR) $(COVERAGE_DIR)
	@$(GOCMD) clean -cache -testcache
	@echo "$(GREEN)✓ Limpieza completa$(RESET)"

# ============================================
# Meta
# ============================================

all: clean deps fmt vet swagger test build ## Build completo
	@echo "$(GREEN)🎉 Build completo$(RESET)"

info: ## Info del proyecto
	@echo "$(BLUE)📋 $(APP_NAME)$(RESET)"
	@echo "  Versión: $(VERSION)"
	@echo "  Ambiente: $(APP_ENV)"
	@echo "  Go: $$($(GOCMD) version)"

.PHONY: help build build-debug run dev test test-coverage test-unit test-unit-coverage test-integration test-integration-verbose test-integration-coverage test-all test-watch coverage-report coverage-check test-analyze test-missing test-validate benchmark fmt vet lint audit deps tidy tools swagger docker-build docker-run docker-stop docker-logs dev-setup dev-teardown dev-reset dev-logs dev-status ci pre-commit clean all info configctl config-validate config-docs
