run-backend:
	go run ./cmd/web/.

run-frontend:
	cd ./ui && npm run dev

run: 
	run-backend run-frontend

build-run:
	cd ./ui && npm run export
	cd ..
	rm ./main
	CGO_ENABLED=0 go build -ldflags "-w" -a -o main ./cmd/web
	./main

build:
	cd ./ui && npm run export
	cd ..
	rm ./main
	CGO_ENABLED=0 go build -ldflags "-w" -a -o main ./cmd/web

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