
GO_ENV := GO111MODULE=on CGO_ENABLED=0


.PHONY: init tidy test build
init: 
	${GO_ENV} go mod init
tidy: 
	${GO_ENV} go mod tidy
test: 
	${GO_ENV} go test -v ./...
build: 
	${GO_ENV} go build -v ./...