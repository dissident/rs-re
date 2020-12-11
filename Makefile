default: run

run:
	@go run *.go

test:
	@go test -v -cover

.PHONY: fmt test

