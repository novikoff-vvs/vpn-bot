#!/bin/bash

# Параметры по умолчанию
BASE_IMAGE_NAME="${1:-vpn-bot}"
TAG="${2:-latest}"
REGISTRY="${3:-registry.novvs.ru}"

# Список сервисов — впиши сюда нужные руками
SERVICES=("bot-service" "user-service" "payment-service" "vpn-service")

# Билдим и пушим каждый сервис
for SERVICE_NAME in "${SERVICES[@]}"; do
    IMAGE_NAME="${BASE_IMAGE_NAME}-${SERVICE_NAME}"

    echo "Собираем образ для сервиса $SERVICE_NAME..."
    docker build \
        --build-arg SERVICE_NAME="$SERVICE_NAME" \
        -t "$REGISTRY/$IMAGE_NAME:$TAG" \
        -f ./Dockerfile-$SERVICE_NAME \
        .

    echo "Пушим образ $REGISTRY/$IMAGE_NAME:$TAG..."
    docker push "$REGISTRY/$IMAGE_NAME:$TAG"

    echo "Готово для $SERVICE_NAME!"
    echo "----------------------------"
done
