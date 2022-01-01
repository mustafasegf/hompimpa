run:
	go run main.go

watch:
	air -c watcher.conf

build:
	go build -o ./bin/main main.go

up:
	docker-compose up -d
	docker-compose logs -f

down:
	docker-compose down

upd:
	docker-compose -f docker-compose.dev.yml up -d
	docker-compose logs -f

downd:
	docker-compose -f docker-compose.dev.yml down
