
# Local config
CONTAINER_NAME=home-exporter
CONTAINER_PORT_HTTP=80
CONTAINER_PORT_PROM=2112
NETWORK_NAME=home-exporter-local
RELEASE_VERSION=0.0.3

REPO=vovanms/home_exporter
IMAGE_RELEASE=$(REPO):$(RELEASE_VERSION)
IMAGE_DEV=$(REPO):dev
IMAGE_LATEST=$(REPO):latest

.PHONY: run build tag push start stop deploy

run:
	go run main.go

build:
	docker build . -t $(IMAGE_DEV)

tag:
	docker tag $(IMAGE_DEV) $(IMAGE_RELEASE)
	docker tag $(IMAGE_DEV) $(IMAGE_LATEST)

push:
	docker push $(IMAGE_RELEASE)
	docker push $(IMAGE_LATEST)

deploy:
	kubectl apply -f deploy/deployment.yml

start:
	docker run \
      --rm \
      --name $(CONTAINER_NAME) \
      -p $(CONTAINER_PORT_HTTP):80 \
      -p $(CONTAINER_PORT_PROM):2112 \
      -v $(shell pwd)/config.yml:/app/config.yml \
      $(IMAGE_DEV)

stop:
	docker stop $(CONTAINER_NAME)
