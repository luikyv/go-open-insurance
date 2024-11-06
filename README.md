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

If you're running MockIn directly on your machine instead of in a Docker container, add this additional entry. It ensures MockIn can resolve the mocked directory served by the NGINX container:
```bash
127.0.0.1 directory
```

If you are developing or modifying this project, start by running `make setup-dev`. For this you will need:
* Docker and Docker Compose installed.
* Go 1.22.x installed and properly configured in the development environment.
* Pre-commit installed for managing Git hooks.
* `jq` installed for JSON processing in the Makefile commands.

Once the setup is complete, you'll be able to use all other make commands.

If you only need to run the project without modifying it, you can use the simpler setup with `make setup`. For this you only need Docker and Docker Compose installed. After this setup, you can start the services using `make run`.

## Dependencies
This project relies significantly on some Go dependencies that streamline development and reduce boilerplate code.

### oapi-codegen
[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) is used for generating Go code based on the Open Insurance OpenAPI specifications. It simplifies the process of creating schemas, reducing the need to handle HTTP requests directly with the Go standard library.
We recommend reviewing the oapi-codegen documentation, particularly the section on [Strict Server](https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#strict-server), which includes some examples.

The configurations for this module are located in `tools/oapi-config.yml`.

### go-oidc
[go-oidc](https://github.com/luikyv/go-oidc) is a configurable OpenID provider written in Go. It handles OAuth-related functionalities, including authentication, token issuance, and scopes. Familiarity with this library's concepts is important for understanding the project's implementation of these aspects.

## TODOs
* Make mongo db remove expired records.
* Env. Defaults to DEV a log warning?
* Dynamic fields.
* Implement user session.
* Data generators.
* Business.
* Add more logs.
* Better way to generated the software statement assertion.
