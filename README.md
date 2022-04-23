# go-httpd

A simple webserver in go based on ideas from [OpenBSD httpd](https://man.openbsd.org/httpd).

## Usage on OpenBSD

* Install go if needed with `pkg_add go`
* `make` - build `go-httpd` executable
* `doas make install`
  * Install executable to `/usr/local/bin/go-httpd`
  * Install rc.d script to `/etc/rc.d/gohttpd`
* Copy example config file to /etc:
  * `doas cp configfiles/gohttpd.json /etc`
  * Modify `/etc/gohttpd.json` as needed
* Set rcctl flag for config file path
  * `doas rcctl set flags gohttpd /etc/gohttpd.json`
* Create directory for request logs if enabled
  * `doas mkdir -p /var/www/gohttpd-logs`
  * `doas chown www:www /var/www/gohttpd-logs`
* Enable and start daemon
  * `doas rcctl enable gohttpd`
  * `doas rcctl start gohttpd`
