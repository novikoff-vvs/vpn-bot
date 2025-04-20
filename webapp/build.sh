#!/bin/bash

# Параметры (можно передавать при запуске)
IMAGE_NAME="${1:-vpn-bot2-front}"
TAG="${2:-latest}"
REGISTRY="${3:-registry.novvs.ru}"

# Собираем Docker-образ
docker build -t registry.novvs.ru/vpn-bot2-front:latest .

# Пушим в registry
docker push "$REGISTRY/$IMAGE_NAME:$TAG"