include ./.env

DBURL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONPATH=db/migrations

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONPATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) up

migrate-down:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) down