DOCKER_IMAGE_VERSION = 1.0.0


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
