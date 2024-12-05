.PHONY: build clean tidy

APP_NAME ?= server
BUILD_DIR ?= build
OUTPUT := $(BUILD_DIR)/$(APP_NAME)
MAIN_FILE := cmd/server.go

build:
	mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(OUTPUT) $(MAIN_FILE)

clean:
	rm -f $(OUTPUT)

tidy:
	go mod tidy
