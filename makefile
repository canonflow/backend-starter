.PHONY: migrate-create migrate-up migrate-down migrate-version migrate-force migrate-drop key-generate uuid-v7

include .env
export

migrate-create:
ifndef NAME
	$(error NAME is not set. Usage: make migrate-create NAME=<migration_name>)
endif

ifndef DB_DRIVER
	$(error DB_DRIVER is not set. Make sure DB_DRIVER is defined in your .env file)
endif
	@echo "Creating migration: $(NAME)"
	@mkdir -p migrations/$(DB_DRIVER)
	migrate create -ext sql -dir ./migrations/$(DB_DRIVER) -seq $(NAME)
	@echo "Migration created in ./migrations/$(DB_DRIVER)/"

migrate-up:
	go run cmd/migration/main.go up

migrate-down:
	go run cmd/migration/main.go down

migrate-version:
	go run cmd/migration/main.go version

migrate-force:
ifndef VERSION
	$(error VERSION is not set. Usage: make migrate-force VERSION=<version>)
endif
	@echo "WARNING: You are about to force migration to version $(VERSION)."
	@read -p "Are you sure? [y/N]: " confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	@read -p "This may cause data inconsistency. Type 'force' to confirm: " final; \
	if [ "$$final" != "force" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	go run cmd/migration/main.go force $(VERSION)

migrate-drop:
	@echo "WARNING: You are about to DROP ALL migrations. This action is IRREVERSIBLE."
	@read -p "Are you sure? [y/N]: " confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	@read -p "Type 'drop' to confirm: " final; \
	if [ "$$final" != "drop" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	go run cmd/migration/main.go drop

key-generate:
	go run cmd/scripts/key_generator.go

uuid-v7:
	go run cmd/scripts/uuid.go