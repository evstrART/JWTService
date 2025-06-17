migrateup:
	go run cmd/migrations/main.go -up

migratedown:
	go run cmd/migrations/main.go -down

build:
	docker-compose up --build

run:
	docker-compose up
createTable:
	@if [ -z "$(name)" ]; then \
		echo "Enter table name: make createTable name=table_name"; \
		exit 1; \
	fi && \
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrateup migratedown createTable build run