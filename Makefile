.PHONY:
.SILENT:

build:
	go build -o ./.bin cmd/main.go

run: build
	./.bin