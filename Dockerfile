# syntax=docker/dockerfile:1

# ==========================================
# Etapa 1: Compilación (Builder)
# ==========================================
FROM golang:alpine AS builder

# Instalar dependencias necesarias para compilar y certificados raíz
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Descargar módulos primero para aprovechar el caché de capas de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar binario estático sin CGO y con optimizaciones de tamaño (-w -s)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/clienthub ./cmd/main.go

# ==========================================
# Etapa 2: Imagen Final Ligera (Runner)
# ==========================================
FROM alpine:3.21

# Instalar certificados CA, zona horaria y curl para health checks
RUN apk add --no-cache ca-certificates tzdata curl

# Crear usuario y grupo sin privilegios por seguridad
RUN addgroup -S -g 10001 appgroup && \
    adduser -S -u 10001 -G appgroup appuser

WORKDIR /app

# Copiar el binario compilado y archivos de configuración
COPY --from=builder /app/bin/clienthub /app/clienthub
COPY --from=builder /app/config /app/config

# Asegurar permisos correctos
RUN chown -R appuser:appgroup /app

# Usar el usuario no privilegiado
USER 10001:10001

# Puerto por defecto de la API
EXPOSE 8080

# Health check para monitorear el estado del contenedor
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Comando de inicio
ENTRYPOINT ["/app/clienthub"]

