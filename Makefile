IMAGE_NAME=gcr.io/thedonor/main

.PHONY: build
build:
	docker build -t $(IMAGE_NAME):latest .

.PHONY: push-image
push-image:
	docker push $(IMAGE_NAME):latest