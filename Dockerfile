# syntax=docker/dockerfile:1
#
# Build stage: compile cjk-press for Linux.
FROM golang:1.25-alpine AS builder

WORKDIR /src

# ca-certificates y tzdata son útiles si más adelante se añade la API/web service.
RUN apk add --no-cache git ca-certificates tzdata

# Descargar dependencias primero para aprovechar la caché de Docker.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Binario estático y sin símbolos para reducir el tamaño de la imagen final.
# Cambia GOARCH a arm64 si construyes para Raspberry Pi / servidores ARM.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
    -ldflags "-s -w -X main.version=dev -X main.commit=dev -X main.date=dev" \
    -o /bin/cjk-press ./cmd/cjk-press

# Runtime stage: imagen mínima y segura (distroless).
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /bin/cjk-press /bin/cjk-press
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nonroot:nonroot

ENTRYPOINT ["/bin/cjk-press"]
CMD ["--help"]