MAIN_FILE	=main.go
BINARY_NAME =go_upcycle_connect_backend
BUILD_DIR	=build

serve:
	@go run $(MAIN_FILE) serve

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Compilation de $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Binaire généré : $(BUILD_DIR)/$(BINARY_NAME)"

migrate:
	@go run $(MAIN_FILE) migrate

clean:
	@echo "Nettoyage des fichiers générés..."
	@rm -rf $(BUILD_DIR)
