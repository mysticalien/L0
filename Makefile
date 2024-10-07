MAIN_APP = main
PRODUCER_APP = producer
DOCKER_COMPOSE_FILE = docker-compose.yml

all: build docker-up wait-for-services run-producer run-main

create-topic:
	@echo "Creating topic orders..."
	docker exec -it kafka kafka-topics --create --topic orders --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

build-main:
	@echo "Building main service..."
	go build -o $(MAIN_APP) ./cmd/main.go

wait-for-services:
	@echo "Waiting for PostgreSQL and Kafka to be ready..."
	chmod +x scripts/wait-for-it.sh
	./scripts/wait-for-it.sh -- echo "PostgreSQL and Kafka is up"

build-producer:
	@echo "Building producer..."
	go build -o $(PRODUCER_APP) ./cmd/producer/producer.go

build: build-producer build-main

docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d

run-main:
	@echo "Running main service..."
	./$(MAIN_APP)

run-producer:
	@echo "Running main service..."
	./$(PRODUCER_APP)


docker-down:
	@echo "Stopping services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

test:
	@echo "Running tests..."
	go test ./...

test_wrk:
	@echo "WRK tests..."
	wrk -t12 -c400 -d30s http://localhost:8080/order/b563feb7b2b84b6test

test_vegeta:
	@echo "Vegeta tests..."
	cat tests/targets.txt | vegeta attack -duration=30s -rate=100 | vegeta report

clean:
	@echo "Cleaning up binaries..."
	rm -f $(MAIN_APP) $(PRODUCER_APP)

redpanda-console:
	@echo "Starting Redpanda Console..."
	docker run -p 8080:8080 -e KAFKA_BROKERS=localhost:9092 docker.redpanda.com/vectorized/redpanda-console:latest
