setup:
	@make keys
	@make setup-cs

setup-dev:
	@go mod download
	@pre-commit install
	@make keys
	@make setup-cs

# Clone and build the Open Insurance Conformance Suite.
# Note: The Dockerfile to build the conformance suite jar is missing, then it is
# being added it manually.
setup-cs:
	@if [ ! -d "conformance-suite" ]; then \
	  echo "Cloning open insurance conformance suite repository..."; \
	  git clone --branch main --single-branch --depth=1 https://gitlab.com/raidiam-conformance/open-insurance/open-insurance-brasil.git conformance-suite; \
	  mkdir conformance-suite/server-dev; \
	  echo 'FROM openjdk:17-jdk-slim\n\nRUN apt-get update && apt-get install redir' > conformance-suite/server-dev/Dockerfile; \
	fi

	@make build-cs

# Run MockIn.
run:
	@docker-compose --profile main up

# Run MockIn alongside the Conformance Suite.
run-with-cs:
	@docker-compose --profile main --profile conformance up

# Run only the MockIn dependencies for local development.
run-dev:
	@docker-compose --profile dev up

# Run only the MockIn dependencies for local development alongside the Conformance Suite.
run-dev-with-cs:
	@docker-compose --profile dev --profile conformance up

# Run the Conformance Suite.
run-cs:
	docker compose --profile conformance up

# Generate certificates, private keys, and JWKS files for both the server and clients.
keys:
	@go run cmd/keymaker/main.go

# Generate API models from the Open Insurance OpenAPI Specification.
models:
	@go generate ./...

# Build the MockIn Docker Image.
build-mockin:
	@docker-compose build mockin

build-cs:
	@docker compose run cs-builder
