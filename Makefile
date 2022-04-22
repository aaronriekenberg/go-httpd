GOCMD=go
GOBUILD=$(GOCMD) build

BINARY_NAME=go-httpd

GIT_COMMIT != git rev-parse HEAD

build:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags="-X main.gitCommit=$(GIT_COMMIT)"

install:
	install -o root -g wheel -m 555 go-httpd /usr/local/bin/go-httpd
	install -o root -g wheel -m 555 rc.d/gohttpd /etc/rc.d/gohttpd
