.PHONY: build

build: ## Build docker image
	docker-compose build --no-cache

.PHONY: status logs start stop clean image

status: ## Get status of containers
	docker-compose ps

logs: ## Get logs of containers
	docker-compose logs -f

start: ## Start docker containers
	docker-compose up -d

stop: ## Stop docker containers
	docker-compose stop

clean:stop ## Stop docker containers, clean data and workspace
	docker-compose down -v --remove-orphans

image:
	docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64,linux/arm/v6 --tag ghcr.io/adadkins/url_shortener:latest .


.PHONY: test

test: ## Run tests
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down --volumes
	