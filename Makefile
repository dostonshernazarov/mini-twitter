-include .env
export

.PHONY: run
run:
	go run cmd/app/main.go

.PHONY: build
build:
	go build cmd/main.go && ./main

.PHONY: swag-gen
swag-gen:
	swag init -g api/router.go -o api/docs

.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir migrations -seq "$(name)"

.PHONY: migrate-up
migrate-up:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSL_MODE} up

.PHONY: migrate-down
migrate-down:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSL_MODE} down
