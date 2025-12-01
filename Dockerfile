# Multi-stage Dockerfile para API Mobile (EduGo)
# - Builder: compila la aplicación con Go
# - Final: imagen ligera basada en Alpine
#
# Notas:
# - Usa variables de entorno proporcionadas por docker-compose o el entorno.
# - La orquestación de dependencias (esperar postgres, mongo, rabbit) es
#   responsabilidad de docker-compose via depends_on: condition: service_healthy
# - Incluye archivos de configuración en /root/config/ (patrón Viper)

FROM golang:1.25-alpine AS builder

# Dependencias del sistema necesarias en builder
RUN apk add --no-cache git ca-certificates

# Directorio de trabajo en build
WORKDIR /app

# Copiar go.mod / go.sum primero para cachear dependencias
COPY go.mod go.sum ./

# Evitar checksum issue en módulos internos (ajustar según tu entorno)
ENV GONOSUMDB=github.com/EduGoGroup/*

# Descargar dependencias
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar el binario (estático)
# Salida: /app/main
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/main.go

# -------------------------
# Etapa final: imagen ligera
# -------------------------
FROM alpine:latest

# Instalar utilidades necesarias (certificados)
RUN apk add --no-cache ca-certificates

# Crear directorios esperados
WORKDIR /root/

# Copiar binario desde la etapa builder
COPY --from=builder /app/main /root/main

# Copiar archivos de configuración (patrón Viper - igual que api-admin)
COPY config/ /root/config/

# Permisos de ejecución
RUN chmod +x /root/main

# Exponer el puerto que usa la aplicación dentro del contenedor
EXPOSE 8080

# Comando directo (sin entrypoint)
# La orquestación de dependencias es responsabilidad de docker-compose
CMD ["./main"]
