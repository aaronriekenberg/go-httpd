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
* Create directory for request logs if enabled
  * `doas mkdir -p /var/www/gohttpd-logs`
  * `doas chown www:www /var/www/gohttpd-logs`
* Enable and start daemon
  * `doas rcctl enable gohttpd`
  * `doas rcctl start gohttpd`

## Command Line Arguments

* `-h` show usage and exit
* `-f <config file path>` override default config file path `/etc/gohttpd.json`
* `-v` enable verbose logging

## Configuration File

* Example config files in `configfiles` directory
* Config file sections
  * `dropPrivileges`
    * May be omitted to disable dropping privileges
    * If `chrootEnabled` is true calls `chroot` at startup to change root directory to `chrootDirectory`
    * Calls `setgid` at startup with gid for `groupName`
    * Calls `setuid` at startup with uid for `userName`
  * `requestLogger`
    * May be omitted to disable request logging
    * If `logToStdout` is true, write request logs to stdout.  Useful for debugging.
    * Else write request logs to `requestLogFile` (relative to `chrootDirectory`)
  * `servers` list of server configurations
    * `serverID` string server id used for logging only
    * `networkAndListenAddressList` list of addresses and ports to listen on.
    * `timeouts` read and write timeouts for server sockets
    * `responseHeaders` response header keys and values at server level.
    * `locations` list of location configurations.  Applied in configured order when each request is processed.
      * `httpPathPrefix` url path prefix for matching location
      * `responseHeaders` response header keys and values at server-location level.  Can be used to override server level `responseHeaders`.
      * Each `location` contains one of the following location types:
      * `blockedLocation`
        * Always return the specified `responseStatus` with no body
      * `directoryLocation`
        * Use go's `http.FileServer` to serve files in the specified `directoryPath` 
        * `directoryPath` is relative to `chrootDirectory`
        * `stripPrefix` may be specified to strip url prefix elements before file serving
      * `compressedDirectoryLocation`
        * Use `github.com/lpar/gzipped/v2` to serve pre-compressed static files ending in `.gz` or `.br` based on `Accept-Encoding` request header
        * Similar to `gzip-static` option in OpenBSD httpd
        * Configuration fields are the same as `directoryLocation`
      * `fastCGILocation`
        * Use `github.com/yookoala/gofast` to connect to a fastcgi application using a unix socket at `unixSocketPath`
      * `redirectLocation`
        * Send a redirect response using the specified `redirectURL` and `responseStatus`
        * `redirectURL` may contain variables `$HTTP_HOST` and `$REQUEST_PATH`
