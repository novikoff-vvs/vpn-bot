# Сервисы для запуска через air (формат: имя_сервиса=путь)
SERVICES ?= bot-service=./services/bot-service/ user-service=./services/bot-service/
export PATH=$PATH:$(go env GOPATH)/bin

REMOTE_USER := root
REMOTE_HOST := 31.128.33.96
REMOTE_PATH := /app/bot-deploy

deploy:
	ssh $(REMOTE_USER)@$(REMOTE_HOST) "cd $(REMOTE_PATH) && docker compose pull && docker compose up -d --build --force-recreate"

build-bot:
	bash ./build.sh

generate-swagger-user:
	$(shell go env GOPATH)/bin/swag init --dir=services/user-service/src/cmd,services/user-service/src/internal/controller --output=services/user-service/src/docs

generate-swagger-vpn:
	$(shell go env GOPATH)/bin/swag init --dir=services/vpn-service/src/cmd,services/vpn-service/src/internal/controller --output=services/vpn-service/src/docs