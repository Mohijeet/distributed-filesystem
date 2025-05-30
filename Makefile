build:
	//@go build -o bin/fs
	@go run main.go

run:
	@./bin/fs

test:
	@go test ./... -v

