#!/bin/bash

TIMEOUT=30

# Функция для проверки доступности PostgreSQL
wait_for_postgres() {
  local container=$1
  local port=$2

  echo "Waiting for PostgreSQL on $container:$port..."
  for i in $(seq 1 $TIMEOUT); do
    if docker exec "$container" pg_isready -h localhost -p "$port"; then
      echo "PostgreSQL is up!"
      return 0
    fi
    sleep 1
  done
  echo "Timeout: PostgreSQL not available on $container:$port"
  return 1
}

# Функция для проверки доступности Kafka
wait_for_kafka() {
  local container=$1
  local port=$2

  echo "Waiting for Kafka on $container:$port..."
  for i in $(seq 1 $TIMEOUT); do
    # Проверка доступности Kafka с помощью команды для получения списка топиков
    if docker exec "$container" kafka-topics --bootstrap-server localhost:"$port" --list >/dev/null 2>&1; then
      echo "Kafka is up!"
      return 0
    fi
    sleep 1
  done
  echo "Timeout: Kafka not available on $container:$port"
  return 1
}

# Проверка PostgreSQL
wait_for_postgres "postgres_db" 5432 || exit 1

# Проверка Kafka
wait_for_kafka "kafka" 9092 || exit 1

echo "PostgreSQL and Kafka are up!"
