#!/bin/bash

# Сервисы и их пути (можно переопределить через переменную окружения SERVICES)
SERVICES=${SERVICES:-"bot-service=./services/bot-service user-service=./services/user-service payment-service=./services/payment-service"}

# Функция для запуска сервиса через air
start_service() {
    local name=$1
    local path=$2
    echo "🚀 Starting $name (path: $path)"
    cd "$path" || { echo "❌ Error: Cannot cd to $path"; exit 1; }
    air -c .air.toml &  # Запуск в фоне
    cd - > /dev/null || return
}

# Очистка фоновых процессов при завершении скрипта (Ctrl+C)
cleanup() {
    echo "🛑 Stopping all services..."
    pkill -P $$  # Убивает все дочерние процессы
    exit 0
}
trap cleanup SIGINT SIGTERM

# Запуск всех сервисов
for service in $SERVICES; do
    service_name=$(echo "$service" | cut -d'=' -f1)
    service_path=$(echo "$service" | cut -d'=' -f2)
    start_service "$service_name" "$service_path"
done

# Бесконечное ожидание (чтобы скрипт не завершился сразу)
echo "✅ All services started. Press Ctrl+C to stop."
wait