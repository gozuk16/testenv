build:
	go build -o verifyenv

run:
	go run main.go test

test:
	go test -v ./...
