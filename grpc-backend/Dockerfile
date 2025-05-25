FROM golang:1.24.3 AS builder

RUN apt-get update && apt-get install -y curl unzip && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip && \
    unzip protoc-21.12-linux-x86_64.zip -d /usr/local && \
    rm protoc-21.12-linux-x86_64.zip

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

RUN protoc \
    -I proto \
    -I googleapis \
    --go_out=gen --go_opt=paths=source_relative \
    --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=gen --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=gen/swagger \
    proto/osu_music.proto

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o osu-server .

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y libsqlite3-0 && rm -rf /var/lib/apt/lists/*

RUN useradd -m appuser

VOLUME ["/data"]

COPY --from=builder /app/osu-server /usr/local/bin/osu-server

WORKDIR /data
USER appuser

CMD ["osu-server"]
