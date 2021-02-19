include Makefile.helper

build:
	@echo "Building ${SERVICE_NAME}"
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -gcflags='-m -m' -o ${SERVICE_NAME} cmd/${SERVICE_NAME}/main.go
.PHONY: build

build-mac:
	@echo "Building ${SERVICE_NAME}"
	CGO_ENABLED=0 \
	GOOS=darwin \
	GOARCH=amd64 \
	go build -gcflags='-m -m' -o ${SERVICE_NAME} cmd/${SERVICE_NAME}/main.go
.PHONY: build-mac

start: build
	@echo "Starting ${SERVICE_NAME}"
	APP_ENV=development \
	go run cmd/${SERVICE_NAME}/main.go
.PHONY: start