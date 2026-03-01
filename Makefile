include .envrc

.PHONY: run/api
run/api:
	@echo "starting server...."
	go run ./cmd/web -dsn=${CONTINUITY_DB_DSN}

.PHONY: mig/new
mig/new:
	@echo "creating migration files for ${name}..."
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: mig/up
mig/up:
	@echo "running up migrations for ${CONTINUITY_DB_DSN}"
	migrate -path ./migrations -database ${CONTINUITY_DB_DSN} up

.PHONY: mig/down
mig/down:
	@echo "running down migrations for ${CONTINUITY_DB_DSN}"
	migrate -path ./migrations -database ${CONTINUITY_DB_DSN} down
