BIN_DIR=${CURDIR}/.local/bin
VENDOR_PROTO_DIR=${CURDIR}/.local/vendor-proto
SHIPMENT_PROTO_PATH=api/shipment/v1

google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf ${VENDOR_PROTO_DIR}/protobuf &&\
	cd ${VENDOR_PROTO_DIR}/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p ${VENDOR_PROTO_DIR}/google
	mv ${VENDOR_PROTO_DIR}/protobuf/src/google/protobuf ${VENDOR_PROTO_DIR}/google
	rm -rf ${VENDOR_PROTO_DIR}/protobuf

.PHONY: vendor-rm
vendor-rm:
	rm -rf ${VENDOR_PROTO_DIR}

.PHONY: vendor-proto
vendor-proto: vendor-rm google/protobuf

.PHONY: bin-deps
bin-deps:
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(BIN_DIR) go install github.com/gojuno/minimock/v3/cmd/minimock@latest

.PHONY: mock
mock:
	${BIN_DIR}/minimock -i ./internal/services/shipment.Repository -o ./internal/services/shipment/repository_mock_test.go -n RepositoryMock -p shipment

.PHONY: protoc-generate
protoc-generate:
	mkdir -p pkg/${SHIPMENT_PROTO_PATH}
	protoc \
	-I ${SHIPMENT_PROTO_PATH} \
	-I ${VENDOR_PROTO_DIR} \
	--plugin=protoc-gen-go=${BIN_DIR}/protoc-gen-go \
	--go_out pkg/${SHIPMENT_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=${BIN_DIR}/protoc-gen-go-grpc \
	--go-grpc_out pkg/${SHIPMENT_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	./${SHIPMENT_PROTO_PATH}/shipment.proto