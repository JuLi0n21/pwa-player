FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM ubuntu:latest
WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates
RUN update-ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/*.env .

EXPOSE 80

CMD ["./main"]
