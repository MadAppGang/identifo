DOCKER_IMAGE_VERSION = 1.0.0

export IDENTIFO_ADMIN_LOGIN = admin@admin.com
export IDENTIFO_ADMIN_PASSWORD = password


run_boltdb:
	go run main.go --config=file://./cmd/config-boltdb.yaml
run_mem:
	go run main.go --config=file://./cmd/config-mem.yaml
run_mongo:
	go run main.go --config=file://./cmd/config-mongodb.yaml
run_dynamodb:
	AWS_ACCESS_KEY_ID=DUMMYIDEXAMPLE \
	AWS_SECRET_ACCESS_KEY=DUMMYEXAMPLEKEY \
	go run main.go --config=file://./cmd/config-dynamodb.yaml


docker_image:
	docker build  --tag madappgangd/identifo:latest --tag madappgangd/identifo:$(DOCKER_IMAGE_VERSION) .

publish: docker_image
	docker push madappgangd/identifo:latest
	docker push madappgangd/identifo:$(DOCKER_IMAGE_VERSION)


test.all: ## run all tests including integration ones, see readme for information how to set up local test environment
	env IDENTIFO_TEST_INGRATION=1 \
		IDENTIFO_TEST_AWS_ENDPOINT="http://localhost:5000" \
		AWS_ACCESS_KEY_ID='testing' \
		AWS_SECRET_ACCESS_KEY='testing' \
		AWS_SECURITY_TOKEN='testing' \
		AWS_SESSION_TOKEN='testing' \
		go test -race ./...

test.local: ## run tests except integration ones
	go test -race ./...
