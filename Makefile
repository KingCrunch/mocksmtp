
NAME=mocksmtp
VERSION=0.0.0

.PHONY: build
build:
	go get github.com/elazarl/go-bindata-assetfs/...
	go get github.com/jteeuwen/go-bindata/...
	go generate
	govendor sync
	go build

.PHONY: test
test: build
	go test -cover

.PHONY: build-all
build-all: generate
	GOOS=linux GOARCH=amd64 go build -o ${NAME}-${VERSION}-linux-amd64
	GOOS=linux GOARCH=386 go build -o ${NAME}-${VERSION}-linux-amd64
	GOOS=linux GOARCH=arm go build -o ${NAME}-${VERSION}-linux-arm

.PHONY: clean
clean:
	rm -rfv mocksmtp mocksmtp-* bindata_* vendor/github.com vendor/bitbucket.org
