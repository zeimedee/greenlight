run:
	go run ./cmd/api

migrate-up:
	migrate -path=./migrations -database=postgres://greenlight:password@localhost/greenlight?sslmode=disable up

migrate-down:
	migrate -path=./migrations -database=postgres://greenlight:password@localhost/greenlight?sslmode=disable down

migrate-index:
	migrate create -seq -ext .sql -dir ./migrations add_movies_indexes