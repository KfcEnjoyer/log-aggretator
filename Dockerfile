FROM golang:1.24-alpine AS builder

WORKDIR /log-aggregator

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o auth-service ./cmd/auth
RUN go build -o log-service ./cmd/logservice

FROM alpine:latest

WORKDIR /log-aggregator

COPY --from=builder /log-aggregator/auth-service .
COPY --from=builder /log-aggregator/log-service .
COPY --from=builder /log-aggregator/configs ./configs/

RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["./auth-service"]