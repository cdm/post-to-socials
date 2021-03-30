# Makefile

IMAGE := post-to-socials
CONTAINER := post-to-socials

.PHONY: default
default: help

.PHONY: docker_pull
docker_pull: ## Pull docker image from github registry
	@docker pull "${IMAGE}"

.PHONY: docker_build
docker_build: ## Build local docker image
	@docker build -t "${IMAGE}" .

.PHONY: docker_push
docker_push: docker_build ## Push docker image to github image registry
	@docker push "${IMAGE}"

#.PHONY: docker_run
docker_run: ## Run docker image
	@docker run -d --name=${CONTAINER} -p 8080:8080 -v "$${PWD}/csv/auth.csv:/csv/auth.csv" -v "$${PWD}/config.yaml:/config.yaml" -v "$${PWD}/form.yaml:/form.html" "${IMAGE}"

.PHONY: docker_stop
docker_stop: ## Stop docker container
	@docker ps -q --filter name="${CONTAINER}" | xargs -r docker stop

.PHONY: docker_rm
docker_rm: ## Remove docker container
	@docker ps -qa --filter name="${CONTAINER}" | xargs -r docker rm

.PHONY: docker_stoprm
docker_stoprm: | docker_stop docker_rm ## Stop and remove docker container

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build (dynamic) binary
	@go build -o post-to-socials .

.PHONY: build_static_linux
build_static_linux: # Build static binary
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o post-to-socials .

.PHONY: test
test: ## Run tests
	@echo "No tests available."
