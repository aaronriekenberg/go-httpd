GOCMD=go
GOBUILD=$(GOCMD) build

BINARY_NAME=go-httpd

VERSION_TAG != git describe --tags

build:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags="-X github.com/aaronriekenberg/go-httpd/commandline.version=$(VERSION_TAG)"

install:
	install -o root -g wheel -m 555 go-httpd /usr/local/sbin/go-httpd
	install -o root -g wheel -m 555 rc.d/gohttpd /etc/rc.d/gohttpd
