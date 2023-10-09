SRC_FILES := $(shell find . -name "*.go")

DOCKER_REPOSITORY := jnorwood
CONSUMERS_IMAGE := hashbash-consumers-go
MIGRATE_IMAGE := hashbash-migrate
WEBAPP_IMAGE := hashbash-webapp-go

##
# Build targets
##
build: hashbash-cli hashbash-engine hashbash-webapp
	:

hashbash-cli: $(SRC_FILES)
	go build -ldflags "-X main.version=$(shell git describe)" -o hashbash-cli github.com/norwoodj/hashbash-backend-go/cmd/hashbash-cli

hashbash-engine: $(SRC_FILES)
	go build -ldflags "-X main.version=$(shell git describe)" -o hashbash-engine github.com/norwoodj/hashbash-backend-go/cmd/hashbash-engine

hashbash-webapp: $(SRC_FILES)
	go build -ldflags "-X main.buildTimestamp=$(shell date --utc --iso-8601=seconds) -X main.gitRevision=$(shell git rev-parse HEAD) -X main.version=$(shell git describe)" -o hashbash-webapp github.com/norwoodj/hashbash-backend-go/cmd/hashbash-webapp

deb:
	debuild

images: engine webapp migrate

engine:
	docker-compose -f docker/docker-compose-hashbash-deps.yaml -f docker/docker-compose-hashbash.yaml build hashbash-engine
	touch engine

webapp:
	docker-compose -f docker/docker-compose-hashbash-deps.yaml -f docker/docker-compose-hashbash.yaml build hashbash-webapp
	touch webapp

migrate:
	docker-compose -f docker/docker-compose-hashbash-deps.yaml build migrate
	touch migrate

release:
	./scripts/release.sh

push: images
	docker tag $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE) $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE):$(shell git tag -l | tail -n 1)
	docker tag $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE) $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE):$(shell git tag -l | tail -n 1)
	docker tag $(DOCKER_REPOSITORY)/$(MIGRATE_IMAGE) $(DOCKER_REPOSITORY)/$(MIGRATE_IMAGE):$(shell git tag -l | tail -n 1)
	docker push $(DOCKER_REPOSITORY)/$(CONSUMERS_IMAGE):$(shell git tag -l | tail -n 1)
	docker push $(DOCKER_REPOSITORY)/$(WEBAPP_IMAGE):$(shell git tag -l | tail -n 1)
	docker push $(DOCKER_REPOSITORY)/$(MIGRATE_IMAGE):$(shell git tag -l | tail -n 1)

down:
	docker-compose -f docker/docker-compose-hashbash-deps.yaml -f docker/docker-compose-hashbash.yaml down --volumes

schema:
	docker-compose -f docker/docker-compose-hashbash-deps.yaml run --rm migrate

new-schema:
	docker-compose -f docker/docker-compose-hashbash-deps.yaml run --entrypoint /migrate --rm migrate create -dir versions -ext sql $(SCHEMA_NAME)

run-deps:
	docker-compose -f docker/docker-compose-hashbash-deps.yaml up

run:
	docker-compose -f docker/docker-compose-hashbash-deps.yaml -f docker/docker-compose-hashbash.yaml up

clean:
	rm -vf hashbash-cli hashbash-engine hashbash-webapp engine webapp migrate
