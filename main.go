package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/kr/pretty"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/dropprivileges"
	"github.com/aaronriekenberg/go-httpd/logging"
	"github.com/aaronriekenberg/go-httpd/requestlogger"
	"github.com/aaronriekenberg/go-httpd/servers"
)

var gitCommit string

var logger = logging.GetLogger()

func awaitShutdownSignal() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	logger.Fatalf("Signal (%v) received, stopping", s)
}

func main() {
	configFilePath := flag.String("f", "/etc/gohttpd.json", "config file path")
	verboseFlag := flag.Bool("v", false, "enable verbose logging")

	flag.Parse()

	logging.SetVerbose(*verboseFlag)

	logger.Printf("go version = %q gitCommit = %q", runtime.Version(), gitCommit)
	logger.Printf("verboseFlag = %v configFilePath = %q", *verboseFlag, *configFilePath)

	configuration := config.ReadConfiguration(*configFilePath)
	logger.Printf("configuration:\n%# v", pretty.Formatter(configuration))

	servers.CreateServers(configuration.Servers)

	dropprivileges.DropPrivileges(configuration.DropPrivileges)

	requestLogger := requestlogger.CreateRequestLogger(configuration.RequestLogger)

	servers.StartServers(
		configuration.Servers,
		requestLogger,
	)

	awaitShutdownSignal()
}
