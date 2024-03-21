migrateup:
	migrate -database "postgres://$(shell echo $$DB_USERNAME):$(shell echo $$DB_PASSWORD)@$(shell echo $$DB_HOST):$(shell echo $$DB_PORT)/$(shell echo $$DB_NAME)?$(shell echo $$DB_PARAMS)" -path db/migrations up

migratedown:
	migrate -database "postgres://$(shell echo $$DB_USERNAME):$(shell echo $$DB_PASSWORD)@$(shell echo $$DB_HOST):$(shell echo $$DB_PORT)/$(shell echo $$DB_NAME)?$(shell echo $$DB_PARAMS)" -path db/migrations down

run:
	go run main.go

startprom:
	docker run \
	--rm \
	-p 9090:9090 \
	--name=prometheus \
	-v $(shell pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
	prom/prometheus

startgrafana:
	docker volume create grafana-storage
	docker volume inspect grafana-storage
	docker run --rm -p 3000:3000 --name=grafana grafana/grafana-oss || docker start grafana

build:
	GOARCH=amd64 GOOS=linux go build -o $(shell pwd)/build/main main.go

build-docker:
	docker build -t social-app .

run-docker:
	docker run \
	--rm -d \
	--name social-app-container  \
	-e DB_NAME=$(shell echo $$DB_NAME) \
	-e DB_USERNAME=$(shell echo $$DB_USERNAME) \
	-e DB_PASSWORD=$(shell echo $$DB_PASSWORD) \
	-e DB_HOST=host.docker.internal \
	-e DB_PORT=$(shell echo $$DB_PORT) \
	-e JWT_SECRET=$(shell echo $$JWT_SECRET) \
	-e BCRYPT_SALT=$(shell echo $$BCRYPT_SALT) \
	-e S3_ID=$(shell echo $$S3_ID) \
	-e S3_SECRET_KEY=$(shell echo $$S3_SECRET_KEY) \
	-e S3_BUCKET_NAME=$(shell echo $$S3_BUCKET_NAME) \
	-e S3_REGION=$(shell echo $$S3_REGION) \
	-e DB_PARAMS=$(shell echo $$DB_PARAMS) \
	-p 8000:8000 \
	social-app

clean:
	rm -rf $(shell pwd)/build

.PHONY: migrateup migratedown run startprom startgrafana build clean build-docker run-docker