# build docker image
build:
	docker compose build

up-all:
	docker compose up --force-recreate --build

# run docker image
up-db:
	docker compose up -d postgres

stop-db:
	docker compose stop postgres

start-db:
	docker compose start postgres

down-db:
	docker compose down postgres


up-service:
	docker compose up -d app

stop-service:
	docker compose stop app

start-service:
	docker compose start app

down-service:
	docker compose down app

migration-up:
	goose -dir ./migrations postgres "postgres://admin:MyPWD56821@localhost:5433/priroda?sslmode=disable" up

migration-status:
	goose -dir ./migrations postgres "postgres://admin:MyPWD56821@postgres:5433/priroda?sslmode=disable" status

migration-down:
	goose -dir ./migrations postgres "postgres://admin:MyPWD56821@localhost:5433/priroda?sslmode=disable" down