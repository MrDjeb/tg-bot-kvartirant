

build:
	go build -o ./.bin cmd/main.go

run: build
	./.bin

build-image:
	docker build -t mrdjeb/tg-bot-kvartirant:0.1 .

start-container:
	docker run --env-file .env -p 80:80 mrdjeb/tg-bot-kvartirant:0.1