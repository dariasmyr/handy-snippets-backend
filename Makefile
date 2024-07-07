.PHONY: build clean

APP_NAME := build/server

build:
	go build -ldflags="-s -w" -o $(APP_NAME) cmd/server.go

clean:
	rm -f $(APP_NAME)
