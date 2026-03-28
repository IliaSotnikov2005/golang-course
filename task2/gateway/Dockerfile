FROM golang:1.26.1-alpine AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/main.go --parseInternal --parseDependency -o ./docs

RUN go build -o gateway ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add curl

WORKDIR /app

COPY --from=builder /app/gateway .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/config ./config

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./gateway", "-config=config/config.yaml"]
