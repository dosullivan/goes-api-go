# Define variables
IMAGE_NAME=goes-api
TAG=latest

.PHONY: build run clean

build:
	docker build -t $(IMAGE_NAME):$(TAG) .

run:
	docker-compose up --build

clean:
	docker-compose down
	docker rmi $(IMAGE_NAME):$(TAG)
