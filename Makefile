DEFAULT_PG_URL=postgres://postgres:password@localhost:5432/postgres?sslmode=disable

migration-up:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" up

migration-down:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" down

migration-status:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" status

db-up:
	docker compose up -d postgres

db-down:
	docker compose down