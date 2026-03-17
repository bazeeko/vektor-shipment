BIN_DIR=${CURDIR}/bin
VENDOR_DIR=${CURDIR}/vendor
SHIPMENT_PROTO_PATH=api/shipment/v1

google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf ${VENDOR_DIR}/protobuf &&\
	cd ${VENDOR_DIR}/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p ${VENDOR_DIR}/google
	mv ${VENDOR_DIR}/protobuf/src/google/protobuf ${VENDOR_DIR}/google
	rm -rf ${VENDOR_DIR}/protobuf

.PHONY: vendor-rm
vendor-rm:
	rm -rf ${VENDOR_DIR}

.PHONY: vendor-proto
vendor-proto: vendor-rm google/protobuf

.PHONY: bin-deps
bin-deps:
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: protoc-generate
protoc-generate:
	mkdir -p pkg/${SHIPMENT_PROTO_PATH}
	protoc \
	-I ${SHIPMENT_PROTO_PATH} \
	-I ${VENDOR_DIR} \
	--plugin=protoc-gen-go=${BIN_DIR}/protoc-gen-go \
	--go_out pkg/${SHIPMENT_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=${BIN_DIR}/protoc-gen-go-grpc \
	--go-grpc_out pkg/${SHIPMENT_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	./${SHIPMENT_PROTO_PATH}/shipment.proto