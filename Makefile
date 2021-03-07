include Makefile.helper

build:
	@echo "Building ${SERVICE_NAME}:${GIT_SHA}"
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -gcflags='-m -m' -o ${SERVICE_NAME} cmd/${SERVICE_NAME}/main.go
.PHONY: build

start: build
	@echo "Starting ${SERVICE_NAME}:${GIT_SHA}"
	APP_ENV=development \
	go run cmd/${SERVICE_NAME}/main.go
.PHONY: start

