SHELL = /bin/bash

.PHONY: run, init-db, init-kafka

run:
	go run cmd/main.go

init-db:
	docker-compose up -d

init-kafka:
	docker run --rm -d -it -p 2183:2181  -p 3030:3030 -p 8082:8082 -p 8083:8083  -p 8084:8084 -p 9092:9092 -e ADV_HOST=127.0.0.1 landoop/fast-data-dev