migrateup:
	migrate -database "postgres://$(shell echo $$DB_USERNAME):$(shell echo $$DB_PASSWORD)@$(shell echo $$DB_HOST):$(shell echo $$DB_PORT)/$(shell echo $$DB_NAME)?sslmode=disable" -path db/migrations up

migratedown:
	migrate -database "postgres://$(shell echo $$DB_USERNAME):$(shell echo $$DB_PASSWORD)@$(shell echo $$DB_HOST):$(shell echo $$DB_PORT)/$(shell echo $$DB_NAME)?sslmode=disable" -path db/migrations down

rundev:
	go run main.go

.PHONY: migrateup migratedown rundev