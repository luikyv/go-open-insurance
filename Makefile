start:
	@docker-compose up -d
	@cd gopf && go run .

init-keys:
	@mkdir -p keys

	@echo "Generate the server's key and self signed certificate."
	@openssl req -newkey rsa:4096 -keyout keys/server.key -out keys/req.csr -nodes -subj "/CN=server"
	@openssl x509 -req -in keys/req.csr -signkey keys/server.key -out keys/server.crt

	@echo "Generate the client one's key and self signed certificate."
	@openssl req -newkey rsa:4096 -keyout keys/client_one.key -out keys/req.csr -nodes -subj "/CN=client_one"
	@openssl x509 -req -in keys/req.csr -signkey keys/client_one.key -out keys/client_one.crt

	@echo "Generate the client two's key and self signed certificate."
	@openssl req -newkey rsa:4096 -keyout keys/client_two.key -out keys/req.csr -nodes -subj "/CN=client_two"
	@openssl x509 -req -in keys/req.csr -signkey keys/client_two.key -out keys/client_two.crt

	# The client certificate bundle will be used to validate client certificates during mutual tls.
	@echo "Generate the client certificate bundle."
	@cat keys/client_one.crt keys/client_two.crt > keys/client_bundle.crt

	@rm keys/req.csr
