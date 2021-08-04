run:
	go run ./main.go

watch:
	air -c watcher.conf

build:
	go build -o main.go main