default: build run

build:
	docker build -t my-golang-app .

run:
	docker run -it --rm --name my-running-app my-golang-app

.PHONY: bin/example
