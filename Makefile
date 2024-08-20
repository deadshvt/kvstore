DOCKER_COMPOSE := docker-compose -f docker-compose.yml
GO := go

ENV_FILE := .env

help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

lint: ### Run linting
	golangci-lint run
.PHONY: lint

run: ### Run containers
	@echo "Running containers..."
	$(DOCKER_COMPOSE) up --build -d
.PHONY: run

down: ### Stop and remove containers
	@echo "Stopping and removing containers..."
	$(DOCKER_COMPOSE) down
.PHONY: down

logs: ### Show logs
	@echo "Showing logs..."
	$(DOCKER_COMPOSE) logs -f
.PHONY: logs

test: ### Run tests
	@echo "Running tests..."
	$(GO) test ./... -v
.PHONY: test

clean: ### Clean up
	@echo "Cleaning up..."
	$(DOCKER_COMPOSE) down -v --rmi all --remove-orphans
.PHONY: clean

restart: ### Restart containers
	@echo "Restarting containers..."
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) up --build -d
.PHONY: restart
