run: lint start

lint:
	golangci-lint run

start:
	go run ./cmd/sort poems.txt