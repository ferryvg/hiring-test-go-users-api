# 'help' target by default
.PHONY: default test
default: help

ifndef VERBOSE
.SILENT:
endif

USER_ID=$(shell id -u)
GROUP_ID=$(shell id -g)
NAME=users-api
DIST_PATH=./deployments/artifacts
COMPOSE_FILE=./deployments/docker-compose.yml

## Remove build/, deployments/artifacts vendor/ dirs and destroy containers
clean:
	rm -rf ./build ./dist ./vendor "${DIST_PATH}"
	docker-compose -f ${COMPOSE_FILE} down

## Compile daemon and cli into deployments/artifacts/
compile: compile.prepare compile.daemon

compile.prepare:
	mkdir -p "${DIST_PATH}"

compile.daemon:
	rm -f "${DIST_PATH}/${NAME}d" "${DIST_PATH}/${NAME}d-debug"
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-s -w' -a -o "${DIST_PATH}/${NAME}d" ./cmd/${NAME}/main.go
	CGO_ENABLED=0 GOOS=linux go build -gcflags "all=-N -l" -a -o "${DIST_PATH}/${NAME}d-debug" ./cmd/${NAME}/main.go

## Run tests
test:
	go test -v ./...

## Run 'compile'
build: compile

## Run 'build' in docker-compose
build.docker:
	docker-compose -f ${COMPOSE_FILE} run --rm -u ${USER_ID} app bash -c "/tmp/scripts/waiter.sh && make build"

## Run 'test' in docker-compose
test.docker:
	docker-compose -f ${COMPOSE_FILE} run --rm -u ${USER_ID} app bash -c "/tmp/scripts/waiter.sh && make test"

## Run 'serve'
serve: compile
	${DIST_PATH}/${NAME}d

## Run 'serve' in docker-compose
serve.docker:
	docker-compose -f ${COMPOSE_FILE} run --rm -u ${USER_ID} -p 6999:6999 app bash -c "/tmp/scripts/waiter.sh && make serve"

## This help screen
help:
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9\.]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "\033[1;32m %-20s \033[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
