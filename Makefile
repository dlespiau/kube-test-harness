all:

dep:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

lint:
	@.ci/go-lint

unit-tests:
	go test -v . ./logger

integration-tests:
	go test -v ./examples/simple

.PHONY: dep integration-tests lint unit-tests
