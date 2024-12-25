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

test:
	go test -v -cover ./...

server:
	go run cmd/api/main.go

.PHONY: templ
templ:
	go install github.com/a-h/templ/cmd/templ@latest
	templ generate

.PHONY: deps
deps:
	go mod tidy
	go get github.com/a-h/templ

.PHONY: generate
generate: deps templ

.PHONY: postgres createdb dropdb migrateup migrateup_last migratedown migratedown_last sqlc test server
