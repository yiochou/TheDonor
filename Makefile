IMAGE_NAME=gcr.io/thedonor/main

.PHONY: deploy
deploy:
	terraform validate; \
	TAG=$$(head -n 4096 /dev/urandom | openssl sha1 | cut -c 1-12); \
	docker build -t $(IMAGE_NAME):$$TAG .; \
	docker push $(IMAGE_NAME):$$TAG; \
	export TF_VAR_IMAGE_TAG=$$TAG; \
	terraform plan; \
	terraform apply
	