FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/app/main.go -o docs

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/app/main.go

# ====================== Финальный образ ======================
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client bash

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations



EXPOSE 8080

CMD ["./main"]