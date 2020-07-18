.PHONY: build
all:
	go build -o bin/net3 ./cmd/net3

.PHONY: install
install:
	go install ./cmd/net3

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: clean
clean:
	rm -rf bin

