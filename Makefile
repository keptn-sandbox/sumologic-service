IMG_TAG ?= dev
IMG_REPO ?= ghcr.io/keptn-sandbox/sumologic-service

build:
	go build -ldflags '-linkmode=external' -v -o sumologic-service

test:
	go test -race -v ./...

docker-build:
	docker build . --network=host -t $(IMG_REPO):$(IMG_TAG)

docker-run:
	docker run --rm -it -p 8080:8080  $(IMG_REPO):$(IMG_TAG)

docker-push:
	docker push $(IMG_REPO):$(IMG_TAG)