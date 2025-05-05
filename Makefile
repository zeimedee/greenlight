run:
	go run ./cmd/api

migrate-up:
	migrate -path=./migrations -database=postgres://greenlight:password@localhost/greenlight?sslmode=disable up

migrate-down:
	migrate -path=./migrations -database=postgres://greenlight:password@localhost/greenlight?sslmode=disable down