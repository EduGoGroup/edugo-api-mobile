# Dockerfile para API Mobile
# Puerto: 8080

FROM golang:alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates

# Establecer directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum
COPY go.mod go.sum ./

# Variable para bypass checksum (edugo-shared es público pero nuevo)
ENV GONOSUMDB=github.com/EduGoGroup/*

# Descargar dependencias
RUN go mod download

# Copiar código fuente completo
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Etapa final - imagen ligera
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binario compilado desde builder
COPY --from=builder /app/main .

# Exponer puerto
EXPOSE 8080

# Comando de inicio
CMD ["./main"]
