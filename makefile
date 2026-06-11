MAIN_FILE				=main.go
BINARY_NAME 			=go_upcycle_connect_backend
BUILD_DIR				=build
DOCKER_COMPOSE_DEV_FILE	=docker-compose.dev.yml
APP_NAME				=app_upcycle-connect

help:
	@echo "Available commands:"
	@echo ""
	@echo "  serve                 - Start the Go server locally"
	@echo "  build                 - Compile the project"
	@echo "  clean                 - Remove generated files"
	@echo ""
	@echo "  docker-dev            - Start containers in dev mode"
	@echo "  docker-dev-up         - Start containers in dev mode"
	@echo "  docker-dev-down       - Stop containers"
	@echo "  docker-dev-logs       - Show container logs"
	@echo ""
	@echo "  migrate               - Run migrations (local)"
	@echo "  docker-migrate        - Run migrations (inside container)"
	@echo ""
	@echo "  generate-keys         - Generate RSA keys (private_key.pem, public_key.pem)"
	@echo ""
	@echo "  model {name}          - Generate a model"
	@echo "  handler {name}        - Generate a handler"
	@echo "  action {name}         - Generate an action"

# --- Development ---
serve:
	@go run $(MAIN_FILE) serve

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Compiling $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Binary generated: $(BUILD_DIR)/$(BINARY_NAME)"

clean:
	@echo "Cleaning generated files..."
	@rm -rf $(BUILD_DIR)

migrate:
	@go run $(MAIN_FILE) migrate

# --- Docker ---
docker-dev:
	@docker compose -f $(DOCKER_COMPOSE_DEV_FILE) up -d --build

docker-dev-up:
	@docker compose -f $(DOCKER_COMPOSE_DEV_FILE) up -d --build

docker-dev-down:
	@docker compose -f $(DOCKER_COMPOSE_DEV_FILE) down

docker-dev-logs:
	@docker compose -f $(DOCKER_COMPOSE_DEV_FILE) logs -f

docker-migrate:
	@docker exec -it $(APP_NAME) go run $(MAIN_FILE) migrate

# --- Code Generation ---
model:
	@./templates/scripts/model ${name}

handler:
	@./templates/scripts/handler ${name}

action:
	@./templates/scripts/action ${name}
