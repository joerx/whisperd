.PHONY: default
default: clean build

.PHONY: test
test:
	go test ./...

.PHONY: build
build: out/whisperd

out/whisperd:
	go build -o out/whisperd .

.PHONY: clean
clean:
	rm -rf out
