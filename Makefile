BINARY := ghell

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test ./shell/ -v

lint:
	go vet ./...

clean:
	rm -f $(BINARY)

install: build
	mv $(BINARY) /usr/local/bin/$(BINARY)

.PHONY: build run test lint clean install
