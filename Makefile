.PHONY: build docker

build:
	CGO_ENABLED=0 go build

docker:
	docker build -t acme-reverseproxy .
