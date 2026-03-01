include .envrc

.PHONY: run/api
run/api:
	@echo "starting server...."
	air

.PHONY: mig/new
mig/new:
	@echo "creating migration files for ${name}..."
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: mig/up
mig/up:
	@echo "running up migration for ${CONTINUITY_DB_DSN}"
	migrate -path ./migration -database ${CONTINUITY_DB_DSN} up

.PHONY: mig/down
mig/down:
	@echo "running down migration for ${CONTINUITY_DB_DSN}"
	migrate -path ./migration -database ${CONTINUITY_DB_DSN} down
