.PHONY: run stop-docker up-docker restart-docker run-all

run:
	clear
	go run main.go

stop-docker:
	clear
	docker-compose down

up-docker:
	clear
	docker-compose up -d

restart-docker:
	clear
	docker-compose down
	docker-compose up -d

run-all:
	clear
	docker-compose down
	docker-compose up -d
	sleep 3
	go run main.go