postgres:
	docker run --name sampleBankDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=samplebank -d postgres:12.4-alpine

createdb:
	docker exec -it sampleBankDB createdb --username=root --owner=root samplebank

dropdb:
	docker exec -it sampleBankDB dropdb samplebank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/samplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/samplebank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server