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
generate-mockgen:
	find . -name 'mock_*.go' -delete
	make generate-ifacemaker
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
	ENV=test TEST_DATABASE_URL=$(DEFAULT_TEST_PG_URL) TEST_KAFKA_BROKER=$(DEFAULT_TEST_KAFKA_BROKER) go test ./... -tags=integration