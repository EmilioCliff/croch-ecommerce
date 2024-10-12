sqlc:
	cd backend/.envs/configs && sqlc generate

start:
	docker compose up

.PHONY: sqlc start