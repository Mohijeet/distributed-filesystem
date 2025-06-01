build:
	@go build -o bin/fs
	# run
	@go run main.go

run:
	@go run .

test:
	@go test ./... -v

