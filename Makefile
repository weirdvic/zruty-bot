BOT_NAME := zruty-bot

all: build restart

.PHONY: build
build:
	docker build --tag ${BOT_NAME}:latest .

.PHONY: restart
	docker stop ${BOT_NAME} || true
	docker rm ${BOT_NAME}
	docker run --name=${BOT_NAME} --restart=always ${BOT_NAME}