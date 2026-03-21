# website_backend

what this project do is to make a DIY website backend and then learn the frontend. Finally being a full-stack programmer.

## configuration

1. install all the tool and library you need
```bash
# arch linux install
go install 
```
2. install docker and pull the postgres:18-alpine image
```bash
docker pull postgres:18-alpine
```
3. run the following make commands
```bash
make postgres
make createdb
make migrateup
make server
```
4. then you can connect the database app and use postman to test the server function of connection to the database

## make command

## tool use

### gomigrate

**create**
```shell
migrate create -ext sql -seq migrate_name
```

**migrate up**
```shell
migrateup:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose up 1
```

**migrate donw**
```shell
migratedown:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/website?sslmode=disable" -verbose down 1
```


### sqlc

you must write the sql.yaml, to list the query sql folder.

And you can control the go file generate in yaml, to add struct tag, return empty slices instead of nil.

```yaml
version: "2"
sql:
  - schema: "db/migration/"
    queries: "db/query/"
    engine: "postgresql"
    gen:
      go:
        package: "db"
        out: "./db/sqlc/"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_exact_table_names: true
        emit_empty_slices: true
```

**sqc generate**
```shell
sqlc generate
```

### viper
