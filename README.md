
# Trace Go collector

## Usage
For an example usage please see the example service in the example folder.

## Example project
Running the example project:
```bash
go run example/main.go
```
To initiate a request:
```bash
curl localhost:9876/test
```
The sample application will call another endpoint.

## Development
The project uses Makefile for building.
Testing:
```bash
make test
```
The test goal will run go fmt, vet on the code. </br>
To see test coverage:
```bash
make coverage
```
