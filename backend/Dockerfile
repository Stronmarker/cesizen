# ── Dev : hot-reload via air ───────────────────────────────────────────────
FROM golang:1.25-alpine AS dev
ENV GOFLAGS=-mod=mod
RUN go install github.com/air-verse/air@latest
WORKDIR /app
CMD ["air", "-c", ".air.toml"]

# ── Builder : compile un binaire statique ──────────────────────────────────
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Le contexte est filtré par backend/.dockerignore (.env, secrets, tmp exclus).
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api

# ── Prod : image minimale non-root ─────────────────────────────────────────
FROM alpine:3.20 AS prod
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /app/server /app/server
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/server"]
