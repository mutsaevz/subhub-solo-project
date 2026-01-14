# ---------- build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# system deps (на будущее, gorm / ssl)
RUN apk add --no-cache git ca-certificates

# deps
COPY go.mod go.sum ./
RUN go mod download

# source
COPY . .

# build main app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/app/main.go

# build seed (отдельный бинарь)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o seed ./cmd/seed/main.go


# ---------- runtime stage ----------
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/seed .

EXPOSE 8080

CMD ["./app"]