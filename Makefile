all:

lint:
	@.ci/go-lint

unit-tests:
	go build -i .
	go test -v . ./logger

integration-tests:
	go build -i .
	go test -v ./examples/...
	go test -v ./tests

.PHONY: integration-tests lint unit-tests
