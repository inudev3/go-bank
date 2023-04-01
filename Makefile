postgres:
	docker run --name postgres -e POSTGRES_PASSWORD=32518458 -e POSTGRES_USER=root -p 5432:5432 -d postgres:12-alpine
createdb:
	docker exec -it postgres createdb --username=root --owner=root go_bank

dropdb:
	docker exec -it postgres dropdb go_bank
migrateup:
	 migrate -path db/migration -database "postgresql://root:32518458@localhost:5432/go_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:32518458@localhost:5432/go_bank?sslmode=disable" -verbose up 1
migratedown:
	 migrate -path db/migration -database "postgresql://root:32518458@localhost:5432/go_bank?sslmode=disable" -verbose down

migratedown1:
	 migrate -path db/migration -database "postgresql://root:32518458@localhost:5432/go_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/inudev5/go-bank/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mock migratedown1
