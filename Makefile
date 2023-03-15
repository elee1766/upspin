.PHONY: build

INSTALL_DIR := ~/.local/bin

build: tidy
	go build -o ./bin ./cmd/...

install:
	cp ./bin/* $(INSTALL_DIR)

tidy:
	go mod tidy

clean:
	rm -rf bin/*


