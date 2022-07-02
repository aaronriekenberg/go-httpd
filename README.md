# go-httpd

A simple webserver in go based on ideas from [OpenBSD httpd](https://man.openbsd.org/httpd).

## Features

* Simple configuration using JSON
  * See `configfiles` directory for example working configurations.
* Uses go's built-in `net/http` server
  * Supports HTTP/1.1 and HTTP/2.0
  * Multiple servers can be configured with optional TLS.  
  * Easy to use with acme-client, see `configfiles/gohttpd.json` example.
  * Automatic thread creation by go, each request is run in its own goroutine.
* Optional request logging
  * Uses `CombinedLoggingHandler` from `github.com/gorilla/handlers`
  * Uses `gopkg.in/natefinch/lumberjack.v2` to write rotate reuqest log files when they reach a configured size.
  * File I/O for request logging is asynchronous using a go channel.
* Each HTTP server has a configured list of locations that are applied exactly in configured order for each request.
* Configurable response header values at server and server-location levels.
* Blocked locations and HTTP redirect locations.
* Static file and directory servers using standard go `http.FileServer`.
* Pre-compressed file serving using `github.com/lpar/gzipped/v2`
  * Supports brotli and gzip files based on `Accept-Encoding` request header
* Supports FastCGI with unix sockets using `github.com/yookoala/gofast`
* Drops privileges at startup and uses `pledge()`.  Roughly the following happens at startup:
  1. go-httpd daemon is started as root
  2. Read configuration file and TLS certificates as root
  3. Create and bind server sockets (`net.Listener`) as root, allowing use of privileged ports 80 and 443.
  4. Call `chroot` to change root to `/var/www` or other configured directory
  5. Call `setuid` and `setgid` to change to unpriviged `www` user/group or other configured user/group
  6. Call `pledge` to limit system calls to `stdio rpath wpath cpath inet unix`.  
  7. Create request logger if configured.
  8. Create request handlers and start the HTTP servers.
* A noop wrapper for pledge is provided so the app builds and runs on non-OpenBSD OS.

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
      * `locationID` string location id used for logging only
      * `httpPathPrefix` url path prefix for matching location.  If not specified defaults to `""` which matches any URL.
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
        * Use `github.com/yookoala/gofast` to connect to a fastcgi application using a unix socket at `network` (`unix` or `tcp`) and `address`.
        * Optionally specify a `connectionPool` block containing a `size` and `lifetimeMilliseconds`. Defaults to no connection pool if not specified.
      * `redirectLocation`
        * Send a redirect response using the specified `redirectURL` and `responseStatus`
        * `redirectURL` may contain variables `$HTTP_HOST` and `$REQUEST_PATH`
