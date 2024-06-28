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
	go generate -x -run=mockgen ./...

# tests
DEFAULT_TEST_PG_URL=postgres://postgres:password@localhost:5431/test?sslmode=disable

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
	docker compose up -d postgres-test

.PHONY: .test-db-down
test-db-down:
	docker compose down
