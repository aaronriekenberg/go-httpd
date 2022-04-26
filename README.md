# go-httpd

A simple webserver in go based on ideas from [OpenBSD httpd](https://man.openbsd.org/httpd).

## Usage on OpenBSD

* Install go if needed with `pkg_add go`
* Clone this git repo
* `make` - build `go-httpd` executable
* `doas make install`
  * Install executable to `/usr/local/sbin/go-httpd`
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

## Configuration

* Example config files in `configfiles` directory
* Config file sections
  * `dropPrivileges`
    * May be omitted to disable dropping privileges
    * If `chrootEnabled` is true calls `chroot` at startup to change root directory to `chrootDirectory`
    * Calls `setgid` at startup with gid for `groupName`
    * Calls `setuid` at startup with uid for `userName`
  * `requestLogger`
    * May be omitted to disable request logging
    * Write request logs to `requestLogFile` (relative to `chrootDirectory`)
  * `servers` list of server configurations
    * `serverID` string server id used for logging only
    * `networkAndListenAddressList` list of addresses and ports to listen on.
    * `timeouts` read and write timeouts for server sockets
    * `locations` list of location configurations.  Applied in configured order when each request is processed.
      * `httpPathPrefix` url path prefix for matching location
      * Each `location` contains one of the following location types:
      * `blockedLocation`
        * Always return the specified `responseStatus` with no body
      * `directoryLocation`
        * Use go's `http.FileServer` to serve files in the specified `directoryPath` 
        * `directoryPath` is relative to `chrootDirectory`
        * `stripPrefix` may be specified to strip url prefix elements before file serving
        * `cacheControlValue` may be specified to control the `Cache-Control` response header value
      * `compressedDirectoryLocation`
        * Use `github.com/lpar/gzipped/v2` to serve pre-compressed static files ending in `.gz` or `.br` based on `Accept-Encoding` request header
        * Similar to `gzip-static` option in OpenBSD httpd
      * `fastCGILocation`
        * Use `github.com/yookoala/gofast` to connect to a fastcgi application using a unix socket at `unixSocketPath`
      * `redirectLocation`
        * Send a redirect response using the specified `redirectURL` and `responseStatus`
        * `redirectURL` may contain variables `$HTTP_HOST` and `$REQUEST_PATH`
