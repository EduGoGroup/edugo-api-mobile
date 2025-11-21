# Multi-stage Dockerfile para API Mobile (EduGo)
# - Builder: compila la aplicación con Go
# - Final: imagen ligera basada en Alpine con scripts de entrypoint
#
# Notas:
# - Usa variables de entorno proporcionadas por docker-compose o el entorno.
# - Copia los scripts en /scripts y establece ENTRYPOINT al script de entrada.
# - Instala netcat (nc) en la imagen final para que wait-for.sh pueda usarlo.
#
# Recomendación: construir con --no-cache si cambias dependencias o scripts.

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

# Instalar utilidades necesarias (certificados, netcat para wait-for, y bash para entrypoint)
RUN apk add --no-cache ca-certificates netcat-openbsd bash

# Crear directorios esperados
WORKDIR /root/

# Copiar binario desde la etapa builder
COPY --from=builder /app/main /root/main

# Copiar scripts directamente del contexto de build (no desde builder)
# Esto evita problemas de caché cuando los scripts se agregan después del COPY . .
COPY scripts/wait-for.sh /scripts/wait-for.sh
COPY scripts/docker-entrypoint.sh /scripts/docker-entrypoint.sh

# Permisos de ejecución
RUN chmod +x /scripts/wait-for.sh /scripts/docker-entrypoint.sh /root/main

# Exponer el puerto que usa la aplicación dentro del contenedor
EXPOSE 8080

# Usar el entrypoint que espera por dependencias y luego arranca la app
ENTRYPOINT ["/scripts/docker-entrypoint.sh"]

# Comando por defecto (puede ser sobreescrito al ejecutar el contenedor)
CMD ["./main"]
