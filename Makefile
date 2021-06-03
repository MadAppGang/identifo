DOCKER_IMAGE_VERSION = 1.0.0


run_boltdb:
	go run main.go --config=file://./cmd/config-boltdb.yaml

docker_image:
	docker build  --tag madappgangd/identifo:latest --tag madappgangd/identifo:$(DOCKER_IMAGE_VERSION) . 


publish: docker_image
	docker push madappgangd/identifo:latest
	docker push madappgangd/identifo:$(DOCKER_IMAGE_VERSION)
