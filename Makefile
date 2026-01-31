postgres:
	docker run --name postgres18 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:18-alpine

createdb:
	docker exec -it postgres18 createdb --username=root --owner=root website

dropdb:
	docker exec -it postgres18 dropdb --username=root website  

migrateup:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...
.PHONY: postgres createdb dropdb migrateup migratedown sqlc
