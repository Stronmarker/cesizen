FROM golang:1.25-alpine AS dev
ENV GOFLAGS=-mod=mod
RUN go install github.com/air-verse/air@latest
WORKDIR /app
CMD ["air", "-c", ".air.toml"]
