# Paths
PROTO_DIR = proto
GOOGLEAPIS_DIR = googleapis
OUT_DIR = gen

PROTO_FILE = osu_music.proto

# Protoc plugins (assumed in PATH)
PROTOC ?= protoc
PROTOC_GEN_GO ?= protoc-gen-go
PROTOC_GEN_GO_GRPC ?= protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY ?= protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPIV2 ?= protoc-gen-openapiv2

all: generate

generate:
	$(PROTOC) \
		-I $(PROTO_DIR) \
		-I $(GOOGLEAPIS_DIR) \
		--go_out $(OUT_DIR) --go_opt paths=source_relative \
		--go-grpc_out $(OUT_DIR) --go-grpc_opt paths=source_relative \
		--grpc-gateway_out $(OUT_DIR) --grpc-gateway_opt paths=source_relative \
		--openapiv2_out $(OUT_DIR)/swagger \
		$(PROTO_DIR)/$(PROTO_FILE)

clean:
	rm -rf $(OUT_DIR)

.PHONY: all generate clean