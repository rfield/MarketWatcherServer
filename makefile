
PROTO_SRC_DIR = ./proto
PROTO_OUT_DIR = ./generated

build: build_protos
	go build -o bin/mkt_client client/main.go
	go build -o bin/mkt_server server/main.go
	go build -o bin/proxy_server proxy/main.go

build_protos:
	protoc --go_out=$(PROTO_OUT_DIR) \
	--go-grpc_out=$(PROTO_OUT_DIR) \
	--grpc-gateway_out=$(PROTO_OUT_DIR) \
	--proto_path=$(PROTO_SRC_DIR) \
	$(PROTO_SRC_DIR)/*.proto

.PHONY: clean
clean:
	rm -rf $(PROTO_OUT_DIR)/*
	rm -f bin/mkt_client bin/mkt_server bin/proxy_server

# Note: If you want to generate the protobuf files with source-relative paths, uncomment the following lines:
#       --go_opt=paths=source_relative \
#       --go-grpc_opt=paths=source_relative \