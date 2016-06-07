
NAME=visualsmtp
VERSION=0.0.0

.PHONY: build
build: generate
	go build

.PHONY: build-all
build-all: generate
	GOOS=linux GOARCH=amd64 go build -o ${NAME}-${VERSION}-linux-amd64
	GOOS=linux GOARCH=arm go build -o ${NAME}-${VERSION}-linux-arm

.PHONY: clean
clean:
	rm -fv visualsmtp visualsmtp-* bindata_*

.PHONY: generate
generate: dependencies
	go generate

.PHONY: dependencies
dependencies:
	go get ./...
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...
