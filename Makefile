.PHONY: keys

# Prepares the environment by generating keys and setting up the Open Insurance
# Conformance Suite.
setup:
	@make keys
	@make setup-localstack

# Sets up the development environment by downloading dependencies installing
# pre-commit hooks, generating keys, and setting up the Open Insurance
# Conformance Suite.
setup-dev:
	@go mod download
	@pre-commit install
	@make keys
	@make setup-cs
	@make setup-localstack

# Clone and build the Open Insurance Conformance Suite.
# Also, generate a configuration file for the suite using files in /keys.
# A configuration file for the conformance suite is also generated based on the files inside /keys.
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

	@make cs-config

setup-localstack:
	chmod +x init_aws.sh

# Runs the main MockIn components using Docker Compose.
run:
	@docker-compose --profile main up

# Starts MockIn along with the Open Insurance Conformance Suite.
run-with-cs:
	@docker-compose --profile main --profile conformance up

# Runs only the MockIn dependencies necessary for local development. With this
# command the MockIn server can run and be debugged in the local host.
run-dev:
	@docker-compose --profile dev up

# Runs the local development environment with both MockIn and the Conformance
# Suite. With this command the MockIn server can run and be debugged in the
# local host with the Conformance Suite.
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

# Build the Conformance Suite JAR file.
build-cs:
	@docker compose run cs-builder

# Create a Conformance Suite configuration file using the client keys in /keys.
cs-config:
	@jq --arg clientOneCert "$$(<keys/client_one.crt)" \
	   --arg clientOneKey "$$(<keys/client_one.key)" \
	   --arg clientTwoCert "$$(<keys/client_two.crt)" \
	   --arg clientTwoKey "$$(<keys/client_two.key)" \
	   --argjson clientOneJwks "$$(jq . < keys/client_one.jwks)" \
	   --argjson clientTwoJwks "$$(jq . < keys/client_two.jwks)" \
	   --arg clientOneKey "$$client_one_key" \
	   --arg clientTwoCert "$$client_two_cert" \
	   --arg clientTwoKey "$$client_two_key" \
	   --argjson clientOneJwks "$$client_one_jwks" \
	   --argjson clientTwoJwks "$$client_two_jwks" \
	   '.client.jwks = $$clientOneJwks | \
	    .mtls.cert = $$clientOneCert | \
	    .mtls.key = $$clientOneKey | \
	    .client2.jwks = $$clientTwoJwks | \
	    .mtls2.cert = $$clientTwoCert | \
	    .mtls2.key = $$clientTwoKey' cs_config_base.json > cs_config.json

	@echo "Conformance Suite config successfully written to cs_config.json"
