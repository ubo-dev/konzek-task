build:
	@go build -o bin/konzek-task

run: build
	@./bin/konzek-task

test:
	@go test -v ./..