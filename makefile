build:
	go build verifyenv.go

run:
	go run verifyenv.go test

test:
	go test -v ./...
