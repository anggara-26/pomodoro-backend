include .env

build:
	go build -o ${BINARY} ./cmd/api

start:
	@env MONGODB_USERNAME=${MONGODB_USERNAME} MONGODB_PASSWORD=${MONGODB_PASSWORD} MONGODB_HOST=${MONGODB_HOST} PORT=${PORT} ./${BINARY}

restart: build start