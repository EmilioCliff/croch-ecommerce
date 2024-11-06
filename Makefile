DB_URL ?=mysql://myuser:mypassword@tcp(localhost:3306)/mydatabase?multiStatements=true&parseTime=true

sqlc:
	cd backend/.envs/configs && sqlc generate

startDocker:
	docker compose up

startServer:
	cd ./backend/cmd/server && go run main.go

swag:
	swag init -dir ./backend/cmd/server/ -output ./backend/docs/swagger/

statik:
	statik -src=./backend/docs/swagger/ -dest=./backend/docs

createMigrate:
	migrate create -ext sql -dir ./backend/internal/mysql/migrations -seq transaction_table

migrateup:
	migrate -path ./backend/internal/mysql/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path ./backend/internal/mysql/migrations -database "$(DB_URL)" -verbose down

.PHONY: sqlc startDocker startServer swag statik createMigrate migrateup migratedown