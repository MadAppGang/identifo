---
description: Tips for backend developers
---

# Backend development

Backend is written in GO, so we suggest you have installed [Golang](https://golang.org).

To build project run `make build`.

To run linters for project run `make lint` (you need `golangci-lint` installed).

#### Testing

Identifo's backend has module and integration tests. To run only module tests run `make test.module`. To run all tests (including integration tests) you should have docker and docker-compose installed. Then run `make test.all` this will set up test environment, using test/docker-compose.yml, run tests against that environment and then delete this environment.
