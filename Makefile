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
	cd test && ./test.sh

test.module: ## run tests except integration ones
	go test -race ./...


build:
	go build -o ./identifo

lint:
	golangci-lint run -D deadcode,errcheck,unused,varcheck,govet

build_admin_panel:
	rm -rf static/admin_panel
	web_apps_src/update-admin.sh

build_login_web_app:
	rm -rf static/web
	web_apps_src/update-web.sh

build_web: build_admin_panel build_login_web_app

