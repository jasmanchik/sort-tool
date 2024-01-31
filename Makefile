run: lint start

lint:
	golangci-lint run

start:
	go run ./cmd/sort -i=poems.txt -o=test.txt
