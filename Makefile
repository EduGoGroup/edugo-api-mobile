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

test: ## Ejecutar todos los tests (unitarios + integración)
	@echo "$(YELLOW)🧪 Ejecutando todos los tests...$(RESET)"
	@$(GOTEST) -v -race -tags=integration ./... -timeout 10m
	@echo "$(GREEN)✓ Tests completados$(RESET)"

test-unit: ## Tests unitarios (rápido, sin integración)
	@echo "$(YELLOW)🧪 Ejecutando tests unitarios...$(RESET)"
	@$(GOTEST) -v -short -race ./internal/... -timeout 2m
	@echo "$(GREEN)✓ Tests unitarios completados$(RESET)"

test-integration: ## Tests de integración (con testcontainers)
	@echo "$(YELLOW)🐳 Ejecutando tests de integración...$(RESET)"
	@RUN_INTEGRATION_TESTS=true $(GOTEST) -v -tags=integration ./test/integration/... -timeout 10m
	@echo "$(GREEN)✓ Tests de integración completados$(RESET)"

test-all: test-unit test-integration ## Ejecutar unitarios + integración por separado

coverage-report: ## Reporte de cobertura completo (incluye integración)
	@echo "$(YELLOW)📊 Generando reporte de cobertura completo...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@RUN_INTEGRATION_TESTS=true $(GOTEST) -tags=integration -coverprofile=$(COVERAGE_DIR)/coverage.out ./... -timeout 10m || true
	@./scripts/filter-coverage.sh $(COVERAGE_DIR)/coverage.out $(COVERAGE_DIR)/coverage-filtered.out
	@$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage-filtered.out -o $(COVERAGE_DIR)/coverage.html
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage-filtered.out | tail -20
	@echo "$(GREEN)✓ Reporte: $(COVERAGE_DIR)/coverage.html$(RESET)"
	@echo "$(BLUE)💡 Abrir: open $(COVERAGE_DIR)/coverage.html$(RESET)"

coverage-check: ## Verificar cobertura mínima (60%)
	@echo "$(YELLOW)🎯 Verificando cobertura mínima...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -tags=integration -coverprofile=$(COVERAGE_DIR)/coverage.out ./... -timeout 10m
	@./scripts/filter-coverage.sh $(COVERAGE_DIR)/coverage.out $(COVERAGE_DIR)/coverage-filtered.out
	@./scripts/check-coverage.sh $(COVERAGE_DIR)/coverage-filtered.out 60
	@echo "$(GREEN)✓ Cobertura cumple umbral mínimo$(RESET)"

test-watch: ## Watch mode para tests (requiere entr)
	@echo "$(YELLOW)👀 Watching tests...$(RESET)"
	@command -v entr > /dev/null || (echo "$(RED)❌ entr no instalado. Instalar con: brew install entr$(RESET)" && exit 1)
	@find . -name "*.go" | entr -c make test-unit

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
# Ambiente de Desarrollo
# ============================================

dev-init: ## Inicializar ambiente de desarrollo (Docker + DB + seed data)
	@./scripts/dev-init.sh

dev-setup: ## Configurar ambiente local (Docker) - llama a dev-init
	@echo "$(YELLOW)🚀 Configurando ambiente...$(RESET)"
	@./scripts/dev-init.sh

dev-status: ## Mostrar estado del ambiente de desarrollo
	@./scripts/dev-status.sh

dev-reset: ## Resetear ambiente completo (⚠️  DESTRUCTIVO - borra datos)
	@./scripts/dev-reset.sh

dev-teardown: ## Limpiar ambiente local
	@echo "$(YELLOW)🧹 Limpiando ambiente...$(RESET)"
	@./test/scripts/teardown_dev_env.sh

# ============================================
# Análisis y Validación
# ============================================

test-stats: ## Estadísticas de tests
	@echo "$(BLUE)📊 Estadísticas de Tests:$(RESET)"
	@echo "  Unitarios:    $$(find internal -name "*_test.go" -type f | wc -l | xargs)"
	@echo "  Integración:  $$(find test/integration -name "*_test.go" -type f | wc -l | xargs)"
	@echo "  Total:        $$(find . -name "*_test.go" -type f | wc -l | xargs)"

# ============================================
# Docker
# ============================================

docker-build: ## Build imagen Docker
	@echo "$(YELLOW)🐳 Building imagen...$(RESET)"
	@docker build -t edugo/$(APP_NAME):$(VERSION) .
	@echo "$(GREEN)✓ Imagen: edugo/$(APP_NAME):$(VERSION)$(RESET)"

docker-up: ## Levantar servicios con compose
	@docker-compose up -d
	@echo "$(GREEN)✓ API: http://localhost:8080$(RESET)"

docker-down: ## Detener servicios
	@docker-compose down

docker-logs: ## Ver logs de servicios
	@docker-compose logs -f

# ============================================
# CI/CD y Hooks
# ============================================

ci: fmt vet test coverage-check ## Pipeline CI completo
	@echo "$(GREEN)✅ CI completado$(RESET)"

pre-commit: fmt vet test-unit ## Hook pre-commit (rápido)
	@echo "$(GREEN)✅ Pre-commit OK$(RESET)"

# ============================================
# Utilidades
# ============================================

clean: ## Limpiar archivos generados
	@echo "$(YELLOW)🧹 Limpiando...$(RESET)"
	@rm -rf $(BUILD_DIR) $(COVERAGE_DIR)
	@$(GOCMD) clean -cache -testcache
	@echo "$(GREEN)✓ Limpieza completa$(RESET)"

info: ## Información del proyecto
	@echo "$(BLUE)📋 $(APP_NAME)$(RESET)"
	@echo "  Versión:  $(VERSION)"
	@echo "  Ambiente: $(APP_ENV)"
	@echo "  Go:       $$($(GOCMD) version | cut -d' ' -f3)"

# ============================================
# Comandos Compuestos
# ============================================

all: clean deps fmt vet test build ## Build completo desde cero
	@echo "$(GREEN)🎉 Build completo$(RESET)"

quick: fmt test-unit build ## Build rápido (sin integración)
	@echo "$(GREEN)⚡ Build rápido completado$(RESET)"

.PHONY: help build build-debug run dev test test-unit test-integration test-all test-watch \
        coverage-report coverage-check test-stats benchmark \
        fmt vet lint audit deps tidy tools configctl config-validate config-docs \
        swagger docker-build docker-up docker-down docker-logs \
        dev-init dev-setup dev-status dev-reset dev-teardown \
        ci pre-commit clean info all quick
