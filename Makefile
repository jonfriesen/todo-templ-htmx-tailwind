# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go test -race -vet=off ./...
	go mod verify


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## create-migration: Create a new migration with the specified name. Usage: make create-migration name=<migration_name>
.PHONY: create-migration
create-migration:
	$(eval MIGRATION_NAME=$(filter-out $@,$(MAKECMDGOALS)))
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "Usage: make $@ name=<migration_name>"; \
		exit 1; \
	fi; \
	if command -v goose > /dev/null; then \
		echo "=> creating migration '$(MIGRATION_NAME)'"; \
		goose -dir=./sql/migrations sqlite3 ./data.db create $(MIGRATION_NAME) sql; \
	else \
		echo "=> goose not found"; \
		echo "=> run make install-dependencies"; \
	fi

%:
	@:

	
## generate-db-code: Generate database code using sqlc. Requires sqlc to be installed.
.PHONY: generate-db-code
generate-db-code: install-dependencies
	@if command -v sqlc > /dev/null; then \
		echo "=> generating db code"; \
		sqlc generate; \
	else \
		echo "=> sqlc not found"; \
		echo "=> run dev/install"; \
	fi


# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build: build the cmd/web application
.PHONY: build
build:
	go mod verify
	go build -ldflags='-s' -o=./bin/todo ./cmd/todo
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/todo ./cmd/todo

## run: run the cmd/web application
.PHONY: run
run:
	go run github.com/cosmtrek/air@v1.40.4 --c="./air.toml"

# ==================================================================================== #
# DEPENDENCIES
# ==================================================================================== #

## install-dependencies: Install required tools like sqlc and tern if they are not already installed
.PHONY: install-dependencies
install-dependencies:
	@command -v sqlc > /dev/null || { \
		echo "=> installing sqlc"; \
		go install -mod=readonly github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
	}
	@command -v goose > /dev/null || { \
		echo "=> installing goose"; \
		go install -mod=readonly  github.com/pressly/goose/v3/cmd/goose@latest; \
	}
	@command -v templ > /dev/null || { \
		echo "=> installing templ"; \
		go install -mod=readonly  github.com/a-h/templ/cmd/templ@latest; \
	}


