################################
# Variables
################################
DB_URL=postgresql://root:root@localhost:5432/backpago?sslmode=disable

################################
# Database and migrations
################################
.PHONY: new-migration
new-migration:
	migrate create -dir db/migrations -ext sql -seq ${name}

.PHONY: postgres
postgres:
	docker run --name postgres-backpago -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres

.PHONY: run-container
run-container:
	docker start postgres-backpago

.PHONY: createdb
createdb:
	docker exec -it postgres-backpago createdb --username=root --owner=root backpago

.PHONY: dropdb
dropdb:
	docker exec -it postgres-backpago dropdb backpago

.PHONY: migrate-up
migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

.PHONY: create-rabbitmq 
create-rabbitmq: 
	docker run -d --hostname backpago --name backpago-rabbitmq -p 5672:5672 rabbitmq:3

################################
# Tests
################################
.PHONY: test
test:
	go test ./...