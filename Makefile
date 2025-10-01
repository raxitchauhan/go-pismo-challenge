test: lint unit-test

run-migration:
	docker compose --file docker-compose.yml run --build --rm migration

unit-test: 
	@echo "==> Running unit tests..."
	docker-compose build --no-cache unit-test
	docker run --rm go-pismo-unit-test \
		go test -v -parallel=4 -cover -race ./...

lint:
	docker-compose build --no-cache lint
	docker-compose run --rm lint

boot:
	docker compose up --build app

down:
	docker compose down --volumes --remove-orphans

gen:
	go generate -x ./...

dep:
	go mod tidy
	go mod vendor