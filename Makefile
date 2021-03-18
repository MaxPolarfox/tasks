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

docker-build:
	@echo "Building docker image for ${SERVICE_NAME}:${GIT_SHA}"
	APP_ENV=production \
	docker build -t maksimpesetski/${SERVICE_NAME}:${GIT_SHA} .
.PHONY: start

docker-push:
	docker -- push maksimpesetski/${SERVICE_NAME}:${GIT_SHA}
.PHONY: docker-push

k8s-deploy-production:
	kubectl apply -f deployment/k8s/${SERVICE_NAME}/tasks.yaml
.PHONY: production

deploy-production-service:
	make build
	make docker-build
	make docker-push
	make k8s-deploy-production
.PHONY: deploy-production-service