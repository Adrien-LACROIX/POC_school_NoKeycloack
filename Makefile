# Nom du module et dossier binaire
APP_NAME := POC_app
BUILD_DIR := build
ENTRY := ./cmd/main.go

# Par défaut, build
default: build

# Build du projet
build windows:
	@echo "📦 Compilation de $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows/$(APP_NAME) $(ENTRY)
	@echo "✅ Binaire créé : $(BUILD_DIR)/$(APP_NAME)"

build linux amd64:
	@echo "📦 Compilation de $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)/linux/amd64
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/amd64/$(APP_NAME)-amd64 $(ENTRY)
	@echo "✅ Binaire créé : $(BUILD_DIR)/$(APP_NAME)"

build linux arm64:
	@echo "📦 Compilation de $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)/linux/arm64
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/linux/arm64/$(APP_NAME)-arm64 $(ENTRY)
	@echo "✅ Binaire créé : $(BUILD_DIR)/$(APP_NAME)"


# Exécution
run: build
	@echo "🚀 Exécution de $(APP_NAME)..."
	@./$(BUILD_DIR)/$(APP_NAME)

# Nettoyage
clean:
	@echo "🧹 Nettoyage..."
	@rm -rf $(BUILD_DIR)
	@echo "🧼 Fait."

.PHONY: build run clean
