migrateup:
	migrate -database "postgres://$(shell echo $$DB_USERNAME):$(shell echo $$DB_PASSWORD)@$(shell echo $$DB_HOST):$(shell echo $$DB_PORT)/$(shell echo $$DB_NAME)?$(shell echo $$DB_PARAMS)" -path db/migrations up

migratedown:
	migrate -database "postgres://$(shell echo $$DB_USERNAME):$(shell echo $$DB_PASSWORD)@$(shell echo $$DB_HOST):$(shell echo $$DB_PORT)/$(shell echo $$DB_NAME)?$(shell echo $$DB_PARAMS)" -path db/migrations down

rundev:
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

.PHONY: migrateup migratedown rundev startprom startgrafana