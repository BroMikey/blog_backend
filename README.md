# website_backend

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
