docker-up:
	docker-compose up -d

docker-build:
	docker-compose build

docker-clean:
	docker-compose down --rmi local

help:
	@echo "Available targets:"
	@echo "  docker-up    - Start docker containers"
	@echo "  docker-build - Build docker images"
	@echo "  docker-clean - Clean current docker containers and images"