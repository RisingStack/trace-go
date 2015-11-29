default: test

fmt: 
	go fmt ./...

coverage-test: fmt
	go test ./trace -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

test: fmt 
	go vet ./...
	go test ./trace
