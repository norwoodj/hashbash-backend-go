SRC_FILES := $(shell find . -name "*.go")

DOCKER_REPOSITORY := jnorwood
CONSUMERS_IMAGE := hashbash-consumers-go
WEBAPP_IMAGE := hashbash-webapp-go
VERSION_PLACEHOLDER := _VERSION


##
# Build targets
##
build: hashbash-cli hashbash-engine hashbash-webapp
	:

hashbash-cli: $(SRC_FILES) version.txt
	go build -ldflags "-X main.version=$(shell cat version.txt)" -o hashbash-cli github.com/norwoodj/hashbash-backend-go/cmd/hashbash-cli

hashbash-engine: $(SRC_FILES) version.txt
	go build -ldflags "-X main.version=$(shell cat version.txt)" -o hashbash-engine github.com/norwoodj/hashbash-backend-go/cmd/hashbash-engine

hashbash-webapp: $(SRC_FILES) version.txt
	go build -ldflags "-X main.version=$(shell cat version.txt)" -o hashbash-webapp github.com/norwoodj/hashbash-backend-go/cmd/hashbash-webapp


##
# Versioning targets
##
version.txt:
	date --utc "+%y.%m%d.0" > version.txt

update-deb-version: version.txt
	sed -i "s|$(VERSION_PLACEHOLDER)|$(shell cat version.txt)|g" debian/changelog
	touch update-deb-version


##
# debian packaging
##
.PHONY: deb
deb: update-deb-version
	debuild


##
# Docker images
##
.PHONY: images
images: engine webapp

engine: version.txt
	docker-compose -f docker/docker-compose-hashbash.yaml build hashbash-engine
	touch engine

webapp: version.txt
	docker-compose -f docker/docker-compose-hashbash.yaml build hashbash-webapp
	touch webapp

.PHONY: push
push: images
	docker tag $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE):current $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE):$(shell cat version.txt)
	docker tag $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE):current $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE):$(shell cat version.txt)
	docker push $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE):$(shell cat version.txt)
	docker push $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE):$(shell cat version.txt)


##
# Run application
##
.PHONY: down
down:
	docker-compose -f docker/docker-compose-hashbash.yaml down

.PHONY: schema
schema:
	docker-compose -f docker/docker-compose-hashbash-deps.yaml run migrate

.PHONY: run-deps
run-deps:
	HASHBASH_HOST_IP_ADDRESS=$(shell ./get-wan-ip) docker-compose -f docker/docker-compose-hashbash-deps.yaml up

.PHONY: run
run:
	HASHBASH_HOST_IP_ADDRESS=$(shell ./get-wan-ip) docker-compose -f docker/docker-compose-hashbash.yaml up


##
# Cleanup
##
.PHONY: clean
clean:
	rm -vf version.txt hashbash-cli hashbash-engine hashbash-webapp engine webapp

.PHONY: debclean
debclean: version.txt
	sed -i "s|$(shell cat version.txt)|$(VERSION_PLACEHOLDER)|g" debian/changelog
	rm -rf dist version.txt update-deb-version
