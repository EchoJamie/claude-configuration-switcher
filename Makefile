.PHONY: build clean test

build:
	go build -o ccs .

clean:
	rm -f ccs

test:
	go test ./...
