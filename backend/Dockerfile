# ---------- build stage ----------
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Cache deps.
COPY go.mod go.sum ./
RUN go mod download

# Copy source.
COPY . .

# Build static binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/campusx-api ./cmd/api

# ---------- runtime stage ----------
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -S campusx && adduser -S campusx -G campusx

WORKDIR /app

COPY --from=builder /out/campusx-api /app/campusx-api

USER campusx
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/campusx-api"]