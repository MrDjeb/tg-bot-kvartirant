

build:
	go build -o ./.bin cmd/main.go

run: build
	./.bin

build-image:
	docker build -t tg-bot-kvartirant:0.1 .

start-container:
	docker run --name tg-bot-kvartirant -p 80:80 --env-file .env tg-bot-kvartirant:0.1