DEFAULT_PG_URL=postgres://postgres:password@localhost:5432/postgres?sslmode=disable

.PHONY: .migration-up
migration-up:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" up

.PHONY: .migration-down
migration-down:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" down

.PHONY: .migration-status
migration-status:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" status

.PHONY: .db-up
db-up:
	docker compose up -d postgres

.PHONY: .db-down
db-down:
	docker compose down

.PHONY: .generate-mockgen
generate-mockgen: generate-ifacemaker
	find . -name 'mock_*.go' -delete
	go generate -x -run=mockgen ./...

.PHONY: .generate-ifacemaker
generate-ifacemaker:
	ifacemaker -f ./internal/service/order.go -s Order -i orderService -p mock_service -c "DONT EDIT: Auto generated" -o ./internal/service/mocks/order.go
	ifacemaker -f ./internal/storage/wrapper.go -s WrapperStorage -i wrapperStorage -p mock_repository -c "DONT EDIT: Auto generated" -o ./internal/storage/mocks/wrapper.go
	ifacemaker -f ./internal/storage/order.go -s OrderStorage -i orderStorage -p mock_repository -c "DONT EDIT: Auto generated" -o ./internal/storage/mocks/order.go

# tests
DEFAULT_TEST_PG_URL=postgres://postgres:password@localhost:5431/postgres?sslmode=disable
DEFAULT_TEST_KAFKA_BROKER=localhost:9091

.PHONY: .test-migration-status
test-migration-status:
	$(eval PG_URL?=$(DEFAULT_TEST_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" status

.PHONY: .test-migration-down
test-migration-down:
	$(eval PG_URL?=$(DEFAULT_TEST_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" down

.PHONY: .test-migration-up
test-migration-up:
	$(eval PG_URL?=$(DEFAULT_TEST_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" up

.PHONY: .test-db-up
test-db-up:
	docker compose -f docker-compose-test.yml up -d postgres-test

.PHONY: .test-db-down
test-db-down:
	docker compose -f docker-compose-test.yml down

.PHONY: .test-kafka-up
test-kafka-up:
	docker compose up -d zookeeper kafka1 kafka2 kafka3

.PHONY: .test-kafka-down
test-kafka-down:
	docker compose down

.PHONY: .unit-tests
unit-tests:
	ENV=test go test ./...

.PHONY: .integration-tests
integration-tests:
	ENV=test TEST_DATABASE_URL=$(DEFAULT_TEST_PG_URL) TEST_KAFKA_BROKER=$(DEFAULT_TEST_KAFKA_BROKER) go test ./tests/... -tags=integration

# proto
# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

ORDER_PROTO_PATH:=api/proto/order/v1
ORDER_PROTO_PATH_OUT:=api
ORDER_DOCS_PATH:=docs

# Установка всех необходимых зависимостей
.PHONY: .bin-deps
bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest

# Вендоринг внешних proto файлов
vendor-proto: vendor-proto-rm vendor-proto/google/protobuf vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options vendor-proto/validate

vendor-proto-rm:
	rm -fdr 'vendor.proto' || true

# Устанавливаем proto описания protoc-gen-openapiv2/options
.PHONY: vendor-proto/protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.proto/grpc-ecosystem && \
 	cd vendor.proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor.proto/protoc-gen-openapiv2
	mv vendor.proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor.proto/protoc-gen-openapiv2
	rm -rf vendor.proto/grpc-ecosystem



# Устанавливаем proto описания google/protobuf
.PHONY: vendor-proto/google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.proto/protobuf &&\
	cd vendor.proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor.proto/google
	mv vendor.proto/protobuf/src/google/protobuf vendor.proto/google
	rm -rf vendor.proto/protobuf

.PHONY: vendor-proto/google/api
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.proto/googleapis && \
 	cd vendor.proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor.proto/google
	mv vendor.proto/googleapis/google/api vendor.proto/google
	rm -rf vendor.proto/googleapis

.PHONY: vendor-proto/validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.proto/tmp && \
		cd vendor.proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.proto/validate
		mv vendor.proto/tmp/validate vendor.proto/
		rm -rf vendor.proto/tmp

.PHONY: generate-proto
generate-proto: bin-deps vendor-proto
	mkdir -p pkg/$(ORDER_PROTO_PATH_OUT)
	mkdir -p $(ORDER_DOCS_PATH)
	protoc -I api/proto \
		-I vendor.proto \
		${ORDER_PROTO_PATH}/order.proto \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=./pkg/$(ORDER_PROTO_PATH_OUT) --go_opt=paths=source_relative\
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=./pkg/$(ORDER_PROTO_PATH_OUT) --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ./pkg/$(ORDER_PROTO_PATH_OUT)  --grpc-gateway_opt  paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=./$(ORDER_DOCS_PATH) \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:pkg/$(ORDER_PROTO_PATH_OUT)"


