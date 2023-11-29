#!make
include .env
CONNECTION="host=${PSQL_HOST} port=${PSQL_PORT} user=${PSQL_USER} password=${PSQL_PASSWORD} dbname=${PSQL_DATABASE} sslmode=${PSQL_SSL_MODE}"

make run:
	go run cmd/server/server.go

up:
	docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d

up-prod:
	docker-compose -f docker-compose.yml -f docker-compose.production.yml up -d

up-build-prod:
	docker-compose -f docker-compose.yml -f docker-compose.production.yml up --build -d

down:
	docker-compose -f docker-compose.yml -f docker-compose.override.yml down

down-prod:
	docker-compose -f docker-compose.yml -f docker-compose.production.yml down

migration-status:
	goose -dir migrations postgres $(CONNECTION) status

.PHONY: migration-create
migration-create: # call with make migration-create name=[name of migration]
	goose -dir migrations create $(name) sql
	goose -dir migrations fix

migrate-up:
	goose -dir migrations postgres $(CONNECTION) up

migrate-down:
	goose -dir migrations postgres $(CONNECTION) down

migrate-redo:
	goose -dir migrations postgres $(CONNECTION) redo

