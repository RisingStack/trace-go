default: test

fmt: 
	go fmt ./...

run: fmt test
	go test ./...
	go run example/main.go

coverage-test: fmt
	go test ./trace -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

test: fmt 
	go vet ./...
	go test ./trace
