DOCKER_IMAGE_VERSION = 3.0.0

export IDENTIFO_ADMIN_LOGIN = admin@admin.com
export IDENTIFO_ADMIN_PASSWORD = password
# export NODE_OPTIONS=--openssl-legacy-provider

run_boltdb:
	mkdir -p ./data
	go build -o plugins/bin/ github.com/madappgang/identifo/v2/plugins/...
	go run main.go --config=file://./cmd/config-boltdb.yaml
run_mem:
	go run main.go --config=file://./cmd/config-mem.yaml
run_mongo:
	go build -o plugins/bin/ github.com/madappgang/identifo/v2/plugins/...
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
	go build -o plugins/bin/ github.com/madappgang/identifo/v2/plugins/...
	go build -o ./identifo

lint:
	golangci-lint run -D errcheck,unused,govet

build_admin_panel:
	rm -rf static/admin_panel
	web_apps_src/update-admin.sh

build_login_web_app:
	rm -rf static/web/element
	web_apps_src/update-web.sh

build_web: build_admin_panel build_login_web_app


run_ui_tests:
	go build -o plugins/bin/ github.com/madappgang/identifo/v2/plugins/...
	go run main.go --config=file://./cmd/config-boltdb.yaml &
	cd web_apps_src/web-element && npx cypress run
	kill $$(ps  | grep config-boltdb.yaml | awk '{print $1}')

open_ui_tests:
	$$(cd web_apps_src/web-element; npm install; $$(npm bin)/cypress open )

## first install protoc compiler https://grpc.io/docs/protoc-installation/
install_go_protoc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

gen_payload_plugin:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		user_payload_provider/grpc/payload_provider/payload_provider.proto

gen_user_storage_plugin:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		storage/grpc/proto/user.proto

gen_impersonation_plugin:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		impersonation/grpc/impersonation_provider/impersonation.proto
