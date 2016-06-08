APPNAME = rocks
VERSION = 0.1.0-dev

setup:
	glide install

build-all: build-mac build-linux

build:
	go build -o ${APPNAME} .

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -X main.Version=${VERSION}" -v -o ${APPNAME}-linux-amd64 .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -X main.Version=${VERSION}" -v -o ${APPNAME}-darwin-amd64 .

build-ci:
	bin/ci-bootstrap.sh
	CGO_CFLAGS="-I${BASE}/build/include" CGO_LDFLAGS="-L${BASE}/build -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy" go build -o ${APPNAME} .

clean:
	rm -f ${APPNAME}
	rm -f ${APPNAME}-linux-amd64
	rm -f ${APPNAME}-darwin-amd64

all:
	setup
	build
	install

test:
	go test -v github.com/ind9/rocks
	go test -v github.com/ind9/rocks/ops

test-only:
	go test -v github.com/ind9/rocks/${name}

install: build
	sudo install -d /usr/local/bin
	sudo install -c ${APPNAME} /usr/local/bin/${APPNAME}

uninstall:
	sudo rm /usr/local/bin/${APPNAME}
