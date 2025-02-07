postgres:
	docker run --name same-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Kashira -d postgres

createdb:
	docker exec -it same-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it same-postgres dropdb simple_bank

migrateUp:
	migrate -path db/migration -database "postgresql://root:Kashira@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateDown:
	migrate -path db/migration -database "postgresql://root:Kashira@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateUp migrateDown sqlc test
