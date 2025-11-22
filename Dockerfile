# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalar dependencias necesarias para compilar (si las hubiera)
RUN apk add --no-cache git

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o expirapp cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copiar el binario desde el builder
COPY --from=builder /app/expirapp .

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./expirapp"]
