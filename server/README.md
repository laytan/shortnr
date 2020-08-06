# Migrations

[Migration docs](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

```sh
migrate -database "mysql://{username}:{password}@({host}:{port})/{database}" -path db/migrations up
migrate -database "mysql://{username}:{password}@({host}:{port})/{database}" -path db/migrations down

migrate create -ext sql -dir db/migrations create_users_table
```
