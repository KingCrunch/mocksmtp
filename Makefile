
NAME=mocksmtp
VERSION=0.0.0

.PHONY: test
test: vendor/github.com vendor/bitbucket.org
	go test -cover

.PHONY: build-all
build-all: vendor/github.com vendor/bitbucket.org
	GOOS=linux GOARCH=amd64 go build -o ${NAME}-${VERSION}-linux-amd64
	GOOS=linux GOARCH=386 go build -o ${NAME}-${VERSION}-linux-amd64
	GOOS=linux GOARCH=arm go build -o ${NAME}-${VERSION}-linux-arm

.PHONY: clean
clean:
	rm -rfv mocksmtp mocksmtp-* bindata_* vendor/github.com vendor/bitbucket.org

.PHONY: get-deps
vendor/github.com vendor/bitbucket.org get-deps:
	go get github.com/elazarl/go-bindata-assetfs/...
	go get github.com/jteeuwen/go-bindata/...
	go generate
	govendor sync
