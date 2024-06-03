include .env
export

DOCKER_IMAGE_NAME=partner

docker:
	docker build -t $(DOCKER_IMAGE_NAME) .
	docker run $(DOCKER_IMAGE_NAME)

run:
	@go get
	@go run main.go
