# Variables
WEB_DIR = ./ui
CMD_DIR = ./cmd/web
BINARY = main

# Commands
run-backend:
	./config/run.sh backend

run-frontend:
	./config/run.sh frontend

run:
	./config/run.sh

build: 
	build-frontend build-backend

build-frontend:
	cd $(WEB_DIR) && npm run export
	cd ..

build-backend:
	rm -f $(BINARY)
	CGO_ENABLED=0 go build -ldflags "-w" -a -o $(BINARY) $(CMD_DIR)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY) $(CMD_DIR)

docker-up:
	docker compose up --build

docker-down:
	docker compose down

web-up:
	docker compose up -d web

db-up:
	docker compose up -d db

docker-rebuild:
	docker compose down --volumes && docker compose up --build

run-log:
	# go run $(CMD_DIR) >>tmp/info.log 2>>tmp/error.log
	go run $(CMD_DIR) -log

run-lint:
	golangci-lint run --enable-all
