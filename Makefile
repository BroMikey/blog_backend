postgres:
	docker run --name postgres18 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:18-alpine

createdb:
	docker exec -it postgres18 createdb --username=root --owner=root website

dropdb:
	docker exec -it postgres18 dropdb --username=root website  

migratecreate:
	migrate create -ext sql -dir db/migration -seq 

migrateup:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...
.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test migratecreate
