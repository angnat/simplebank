postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_PASSWORD=password.1 -d postgres:14-alpine
stopPostgres:
	docker stop postgres14; docker rm postgres14
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres14 dropdb --username=root simple_bank
migrateup:
	.\migrate -verbose -path ./migration -database "postgresql://postgres:root@localhost:5432/simple_bank?sslmode=disable" up
migratedown:
	.\migrate -verbose -path ./migration -database "postgresql://postgres:root@localhost:5432/simple_bank?sslmode=disable" down
sqlc:
	sqlc generate

test:
	go test -v -cover ./...	
.PHONY: postgres stopPostgres createdb dropdb migrateup migratedown sqlc test
