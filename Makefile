SERVICE_NAME ?= store-node

.PHONY: install-protoc
install-protoc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: api
api: install-protoc
	find api -name "*.proto" | xargs protoc \
  		--go_out=./internal/pb/ \
		--go-grpc_out=./internal/pb/

DOCKER_BUILD := DOCKER_BUILDKIT=1 docker buildx build -f build/Dockerfile -t $(SERVICE_NAME)
DOCKER_COMPOSE := docker compose --project-name decvault --file docker-compose.yaml
DOCKER_REGISTRY := "cr.yandex/crpcv4imm48oafilqbal"

.PHONY: docker-local
docker-local:
	  $(DOCKER_BUILD) .

.PHONY: docker-build
docker-build:
	$(DOCKER_BUILD) --platform linux/amd64 .

.PHONY: docker-push
docker-push:
	IMAGE_FULL_NAME = $(DOCKER_REGISTRY)/$(SERVICE_NAME):latest
	docker tag $(SERVICE_NAME):latest IMAGE_FULL_NAME
	docker push $(IMAGE_FULL_NAME)

.PHONY: docker-release
docker-release: docker-build docker-push

.PHONY: compose-up
compose-up:
	$(DOCKER_COMPOSE) up -d

.PHONY: compose-stop
compose-stop:
	$(DOCKER_COMPOSE) stop

.PHONY: compose-down
compose-down:
	$(DOCKER_COMPOSE) down
