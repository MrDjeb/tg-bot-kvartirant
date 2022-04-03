FROM golang:alpine AS builder

RUN go version

COPY . /github.com/MrDjeb/tg-bot-kvartirant/
WORKDIR /github.com/MrDjeb/tg-bot-kvartirant/

RUN apk add build-base
RUN go mod download 
RUN GOOS=linux go build -o ./.bin cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/MrDjeb/tg-bot-kvartirant/configs configs/

EXPOSE 80

ENTRYPOINT ["./.bin"]

