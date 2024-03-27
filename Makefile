include .env
export

SHELL:=/usr/bin/bash

migrations_status:
	goose -dir=$(MIGRATIONS_DIR) status

migrations_up:
	goose -dir=$(MIGRATIONS_DIR) up

add_sql_migration:
	goose -dir=${MIGRATIONS_DIR} create ${MIGRATION_NAME} sql

run: 
	go run cmd/skud/main.go --config=config/local.yaml

build:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build ./cmd/skud/main.go 