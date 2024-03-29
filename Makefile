postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password.1 -d postgres:14.5-alpine
stopPostgres:
	docker stop postgres14; docker rm postgres14
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres14 dropdb --username=root simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:password.1@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:password.1@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store

.PHONY: postgres stopPostgres createdb dropdb migrateup migratedown sqlc test server mock
