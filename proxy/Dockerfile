FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM ubuntu:latest
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/*.env .

EXPOSE 5002

CMD ["./main"]
