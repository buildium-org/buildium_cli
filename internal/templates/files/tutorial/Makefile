.PHONY: build
build:
	go build -o bin/app main.go

.PHONY: docker-build
docker-build:
	docker build -t <YOUR_IMAGE_NAME_HERE>_harness .

.PHONY: docker-publish
docker-publish:
	docker tag <YOUR_IMAGE_NAME_HERE>_harness:latest <YOUR_DOCKER_USERNAME_HERE>/<YOUR_IMAGE_NAME_HERE>_harness:latest
	docker push <YOUR_DOCKER_USERNAME_HERE>/<YOUR_IMAGE_NAME_HERE>_harness:latest