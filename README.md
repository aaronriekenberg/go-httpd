# go-httpd

A simple webserver in go based on ideas from [OpenBSD httpd](https://man.openbsd.org/httpd).

## Usage on OpenBSD

* Install go if needed with `pkg_add go`
* Clone this git repo
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

## Configuration file sections

* Example config files in `configfiles` directory
* Config file sections
  * dropPrivileges
    * May be null to disable dropping privileges
    * If `chrootEnabled` is true calls `chroot` at startup to change root directory to `chrootDirectory`
    * Calls `setgid` at startup with gid for `groupName`
    * Calls `setuid` at startup with uid for `userName`
  * requestLogger
    * May be null to disable request logging
    * Write request logs to `requestLogFile` (relative to `chrootDirectory`)
  * servers
    * List of server configs
      * `serverID` string server id used for logging only
      * `networkAndListenAddressList` list of addresses and ports to listen on.
      * `timeouts` read and write timeouts for server sockets
      * `locations` list of location configurations.  Applied in configured order when each request is processed.
        * `httpPathPrefix` url path prefix for matching location
        * Each `location` contains one of `blockedLocation`, `directoryLocation`, `redirectLocation`, `fastCGILocation`
