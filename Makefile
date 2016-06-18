
NAME=mocksmtp
VERSION=0.0.0

.PHONY: build
build: dependencies generate
	go build

.PHONY: build-all
build-all: generate
	GOOS=linux GOARCH=amd64 go build -o ${NAME}-${VERSION}-linux-amd64
	GOOS=linux GOARCH=386 go build -o ${NAME}-${VERSION}-linux-amd64
	GOOS=linux GOARCH=arm go build -o ${NAME}-${VERSION}-linux-arm

.PHONY: clean
clean:
	rm -rfv mocksmtp mocksmtp-* bindata_* vendor/github.com vendor/bitbucket.org

.PHONY: generate
generate: dependencies
	go generate

.PHONY: dependencies
dependencies:
	govendor sync
