postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres domain_mail

dropdb:
	docker exec -it postgres dropdb domain_mail

migrateup:
	migrate -path internal/db/migrations/ -database "postgres://postgres:secret@localhost:5432/domain_mail?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/migrations/ -database "postgres://postgres:secret@localhost:5432/domain_mail?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migrateup_last migratedown migratedown_last sqlc
