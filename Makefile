BOT_NAME := zruty-bot

all: build restart

.PHONY: build
build:
	docker build --no-cache --tag ${BOT_NAME}:latest .

.PHONY: restart
restart:
	docker stop ${BOT_NAME} || true
	docker rm ${BOT_NAME} || true
	docker run --name=${BOT_NAME} --restart=always --detach ${BOT_NAME}
