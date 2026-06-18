.PHONY: build
build:
	mkdir -p bin
	rm -f bin/app
	go build -o bin/app main.go
	cp meta.json bin/meta.json

.PHONY: run
run:
	docker compose down
	docker build -t <YOUR_IMAGE_NAME_HERE> .
	docker compose up --abort-on-container-exit --exit-code-from app