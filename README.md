# go-open-insurance
An implementation of the Brazil Open Insurance specifications in Go.

## Open API Specs
The following Open Insurance Open API Specifications are implemented.

### Phase 2
* [API Consents v2.6.0](https://raw.githubusercontent.com/br-openinsurance/areadesenvolvedor/refs/heads/main/documentation/source/files/swagger/consents_v2.yaml)
* [API Resources v2.4.0](https://raw.githubusercontent.com/br-openinsurance/areadesenvolvedor/bf3804bb85d8248a5ea5c45a0a656b732df4975f/documentation/source/files/swagger/resources_v2.yaml)
* [API Customers v1.5.0](https://raw.githubusercontent.com/br-openinsurance/areadesenvolvedor/2e9a2d43d90e6662c2a4dcffc3b95d00d14d41f7/documentation/source/files/swagger/customers.yaml)
* [API Insurance Capitalization Titles v1.4.0](https://raw.githubusercontent.com/br-openinsurance/areadesenvolvedor/e5e54393cafb0988de148ab4c594f86346752cbc/documentation/source/files/swagger/insurance-capitalization-title.yaml)

### Phase 3
* [API Endorsements v1.2.0](https://raw.githubusercontent.com/br-openinsurance/areadesenvolvedor/2f76347b669236ab39c184b68d6e154148f69685/documentation/source/files/swagger/endorsement.yaml)
* [API Quote Auto v1.8.0](https://br-openinsurance.github.io/areadesenvolvedor/files/swagger/quote-auto.yaml)

## Usage and Development Guide

To ensure MockIn works correctly in your local environment, you need to update your system's hosts file (usually located at /etc/hosts on Unix-based systems or C:\Windows\System32\drivers\etc\hosts on Windows). This step allows your machine to resolve the required domains for MockIn.
```bash
127.0.0.1 mockin.local
127.0.0.1 matls-mockin.local
```

If you’re running MockIn directly on your machine instead of in a Docker container, add this additional entry. It ensures MockIn can resolve the mocked directory served by the NGINX container:
```bash
127.0.0.1 directory
```

This project includes a series of Makefile targets to streamline the setup, run, and development process. Down below there is a breakdown of the available commands and their purposes.

If you are developing or modifying this project, start by running `make setup-dev`. For this you will need:
* Docker and Docker Compose installed.
* Go 1.22.x installed and properly configured in the development environment.
* Pre-commit installed for managing Git hooks.

Once the setup is complete, you’ll be able to use all other make commands.

If you only need to run the project without modifying it, you can use the simpler setup with `make setup`. For this you only need Docker and Docker Compose installed.
After this setup, you can start the services using `make run` and `make run-with-cs` which also spins up the Open Insurance Conformance Suite.

### Setup Commands
`make setup` \
Prepares the environment by generating keys and setting up the Open Insurance Conformance Suite.

`make setup-dev` \
Sets up the development environment by downloading dependencies installing pre-commit hooks, generating keys, and setting up the Open Insurance Conformance Suite.

`make setup-cs` \
Clones and prepares the Open Insurance Conformance Suite for use.

### Run Commands
`make run`\
Runs the main MockIn components using Docker Compose.

`make run-with-cs` \
Starts MockIn along with the Open Insurance Conformance Suite.

`make run-dev` \
Runs only the MockIn dependencies necessary for local development. With this command the MockIn server can run and be debugged in the local host.

`make run-dev-with-cs` \
Runs the local development environment with both MockIn and the Conformance Suite. With this command the MockIn server can run and be debugged in the local host with the Conformance Suite.

`make run-cs` \
Starts only the Conformance Suite.

### Utility Commands

`make build-mockin` \
Build the MockIn Docker Image.

`make build-cs` \
Build the Conformance Suite JAR file.

`make keys` \
Generates certificates, private keys, and JWKS files for both the server and clients.

`make models` \
Generates API models from the Open Insurance OpenAPI specification.

## For Developers
This project relies significantly on some Go dependencies that streamline development and reduce boilerplate code.

### oapi-codegen
[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) is used for generating Go code based on the Open Insurance OpenAPI specifications. It simplifies the process of creating schemas, reducing the need to handle HTTP requests directly with the Go standard library.
We recommend reviewing the oapi-codegen documentation, particularly the section on [Strict Server](https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#strict-server), which includes some examples.

The configurations for this module are located in `tools/oapi-config.yml`.

### go-oidc
[go-oidc](https://github.com/luikyv/go-oidc) is a configurable OpenID provider written in Go. It handles OAuth-related functionalities, including authentication, token issuance, and scopes. Familiarity with this library’s concepts is important for understanding the project's implementation of these aspects.

## TODOs
* Env
* Dynamic fields.
* Implement user session.
* Generate cs config file.
* Data generators.
* Business.
