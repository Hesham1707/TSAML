.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fixtures && CHANNEL_NAME=mychannel 	 IMAGE_TAG=latest docker-compose -f docker-compose.yaml up --force-recreate -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fixtures && CHANNEL_NAME=mychannel IMAGE_TAG=latest docker-compose -f docker-compose.yaml down --volumes
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./tsaml

	##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/tsaml-* tsaml
	@docker rm -f -v `docker ps -a --no-trunc | grep "tsaml" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "tsaml" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"