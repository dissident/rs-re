.PHONY: bin/example

default: build run clean

build-base:
	docker build -t rs-re-base -f Dockerfile.base .

build:
	docker build -t rs-re --cache-from rs-re-base .

run:
	docker run -it --rm --name rs-re rs-re

clean-c:
	docker rm $$(docker ps -a -q)

clean:
	docker image rm $$(docker images rs-re:latest -q)
