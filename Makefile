.PHONY: build install tidy clean

INSTALL_DIR := ~/.local/bin

APPS=cacheserver dirserver frontend keyserver storeserver upspin upspin-audit upspinfs upspinserver

build: tidy
	mkdir -p build
	for dir in $(APPS); do\
		go build -o ./build/$$dir ./cmd/$$dir; \
	done

install:
	cp ./build/* $(INSTALL_DIR)

tidy:
	go mod tidy

clean:
	rm -rf build/*


