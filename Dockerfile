# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Instalar dependencias necesarias para el build
RUN apk add --no-cache git

WORKDIR /app

# Copiar archivos de dependencias y descargar
COPY go.mod go.sum ./
RUN go mod download

# Copiar el codigo fuente
COPY . .

# Construir el binario optimizado para Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:latest

# Instalar certificados CA para conexiones TLS salientes
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario desde la etapa de builder
COPY --from=builder /app/server .

# Exponer el puerto configurado en el MVP
EXPOSE 8080

# Comando de inicio
CMD ["./server"]
