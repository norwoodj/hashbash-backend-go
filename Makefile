GOLANG_IMAGE=golang:1.12-alpine3.9
DOCKER_REPOSITORY=jnorwood
CONSUMERS_IMAGE=hashbash-consumers-go
WEBAPP_IMAGE=hashbash-webapp-go


.PHONY: all
all: consumers webapp

consumers: version.txt
	docker build --tag $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE):current --file docker/Dockerfile-consumers .
	touch consumers

webapp: version.txt
	docker build --tag $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE):current --file docker/Dockerfile-webapp .
	touch webapp

version.txt:
	echo release-$(shell docker run --rm --entrypoint date $(GOLANG_IMAGE) --utc "+%Y%m%d-%H%M") > version.txt

.PHONY: push
push: all
	docker tag $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE):current $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE):$(shell cat version.txt)
	docker tag $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE):current $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE):$(shell cat version.txt)
	docker push $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE):$(shell cat version.txt)
	docker push $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE):$(shell cat version.txt)

.PHONY: run-deps
run-deps:
	HASHBASH_HOST_IP_ADDRESS=$(shell ./get-wan-ip) docker-compose -f docker/docker-compose-hashbash-deps.yaml up

.PHONY: run
run: all volume
	docker-compose -f docker/docker-compose-hashbash.yaml up

.PHONY: clear-data
clear-data:
	docker-compose -f docker/docker-compose-hashbash.yaml down
	docker volume rm hashbash-data
	docker volume create hashbash-data

volume:
	docker volume create --name=hashbash-data
	touch volume

.PHONY: clean
clean:
	rm -vf version.txt
	rm -vf cmd/hashbash-engine/hashbash-engine
	rm -vf cmd/hashbash-webapp/hashbash-webapp
