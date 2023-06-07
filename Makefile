run:
	go run ./cmd/web

build:
	go build -o main ./cmd/web

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/web

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

db-up:
	docker-compose up -d db

docker-rebuild:
	docker-compose down --volumes && docker-compose up --build

run-log:
	# go run ./cmd/web >>tmp/info.log 2>>tmp/error.log
	go run ./cmd/web -log

help:
	go run ./cmd/web -help