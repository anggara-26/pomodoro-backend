include .env

.PHONY: build start restart clean test docker-build docker-run swagger

build:
	go build -o ${BINARY} ./cmd/api

start:
	@env MONGODB_USERNAME=${MONGODB_USERNAME} MONGODB_PASSWORD=${MONGODB_PASSWORD} MONGODB_HOST=${MONGODB_HOST} MONGODB=${MONGODB} PORT=${PORT} ./${BINARY}

restart: build start

clean:
	rm -f ${BINARY}
	rm -f pomodoro.exe

test:
	go test ./... -v

swagger:
	swag init -g ./cmd/api/main.go -o ./docs/v1

docker-build:
	docker build -t pomodoro-api .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

dev: swagger build start