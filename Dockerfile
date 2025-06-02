FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/users ./cmd/main/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/users .

ENTRYPOINT ["./users"]
