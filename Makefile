setup:
	@pre-commit install
	@go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@main
	@make init-keys

run-no-server:
	@docker-compose --profile opin up

run-cs:
	@if [ ! -d "conformance-suite" ]; then \
	  echo "Cloning open insurance conformance suite repository..."; \
	  git clone --branch main --single-branch --depth=1 https://gitlab.com/raidiam-conformance/open-insurance/open-insurance-brasil.git conformance-suite; \
	  docker compose run cs-builder; \
	fi
	docker compose --profile conformance up
	
init-keys:
	@mkdir -p keys

	@echo "Generate the server's key and self signed certificate."
	@openssl req -newkey rsa:4096 -keyout keys/server.key -out keys/req.csr -nodes -subj "/CN=server"
	@openssl x509 -req -in keys/req.csr -signkey keys/server.key -out keys/server.crt -days 365

	@echo "Generate the client one's key and self signed certificate."
	@openssl req -newkey rsa:4096 -keyout keys/client_one.key -out keys/req.csr -nodes -subj "/UID=dad19daf-390b-4a4f-8647-77ca98149cb3/organizationIdentifier=OPIBR-dad19daf-390b-4a4f-8647-77ca98149cb3/CN=client_one"
	@openssl x509 -req -in keys/req.csr -signkey keys/client_one.key -out keys/client_one.crt -days 365

	@echo "Generate the client two's key and self signed certificate."
	@openssl req -newkey rsa:4096 -keyout keys/client_two.key -out keys/req.csr -nodes -subj "/UID=9d2ca524-3c44-400d-b7d6-b76eca469cf0/organizationIdentifier=OPIBR-9d2ca524-3c44-400d-b7d6-b76eca469cf0/CN=client_two"
	@openssl x509 -req -in keys/req.csr -signkey keys/client_two.key -out keys/client_two.crt -days 365

	# The client certificate bundle will be used to validate client certificates during mutual tls.
	@echo "Generate the client certificate bundle."
	@cat keys/client_one.crt keys/client_two.crt > keys/client_bundle.crt

	@rm keys/req.csr

generate-models:
	@oapi-codegen --config=oapi-config.yml oapi-spec.yml