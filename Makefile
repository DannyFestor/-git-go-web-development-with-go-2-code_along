DB_HOST=localhost
DB_PORT=5432
DB_USER=lenslocked
DB_PASSWORD=password
DB_NAME=lenslocked
SSL=disable
CONNECTION="host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(SSL)"

migration-status:
	goose -dir migrations postgres $(CONNECTION) status

.PHONY: migration-create
migration-create: # call with make migration-create name=[name of migration]
	goose -dir migrations create $(name)
	goose -dir migrations fix

migrate-up:
	goose -dir migrations postgres $(CONNECTION) up

migrate-down:
	goose -dir migrations postgres $(CONNECTION) down

migrate-redo:
	goose -dir migrations postgres $(CONNECTION) redo