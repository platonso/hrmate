.PHONY: docker-up docker-clean port-forward port-close migrate-up migrate-down migrate-status help

docker-up:
	@docker-compose up -d --build
docker-clean:
	@docker-compose down --rmi local

port-forward:
	@docker-compose up -d port-forwarder
port-close:
	@docker-compose down port-forwarder

migrate-up:
	@docker-compose run --rm postgres-migrator ./migrator -command=up
migrate-down:
	@docker-compose run --rm postgres-migrator ./migrator -command=down
migrate-status:
	@docker-compose run --rm postgres-migrator ./migrator -command=status

help:
	@echo "Docker Management:"
	@echo "  docker-up         - Start docker containers"
	@echo "  docker-clean      - Clean current docker containers and images"
	@echo ""
	@echo "Port Forwarding:"
	@echo "  make port-forward - Start PostgreSQL proxy"
	@echo "  make port-close   - Stop PostgreSQL proxy"
	@echo ""
	@echo "Database Migrations:"
	@echo "  migrate-up        - Apply all pending migrations"
	@echo "  migrate-down      - Rolling back the last migration"
	@echo "  migrate-status    - Show migration status"
