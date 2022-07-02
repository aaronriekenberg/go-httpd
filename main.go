package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kr/pretty"

	"github.com/aaronriekenberg/go-httpd/commandline"
	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/dropprivileges"
	"github.com/aaronriekenberg/go-httpd/logging"
	"github.com/aaronriekenberg/go-httpd/pledge"
	"github.com/aaronriekenberg/go-httpd/requestlogging"
	"github.com/aaronriekenberg/go-httpd/servers"
)

var logger = logging.GetLogger()

func main() {

	commandLineFlags := commandline.ProcessCommandLineFlags()

	logger.SetVerboseEnabled(commandLineFlags.Verbose)

	logger.Printf("starting %v", commandline.AppName())
	logger.Printf("commandLineFlags = %+v", commandLineFlags)

	configuration := config.ReadConfiguration(commandLineFlags.ConfigFilePath)
	logger.Printf("configuration:\n%# v", pretty.Formatter(configuration))

	servers.CreateServers(configuration.Servers)

	dropprivileges.DropPrivileges(configuration.DropPrivileges)

	pledge.Pledge()

	requestLogger := requestlogging.NewRequestLogger(configuration.RequestLogger)

	servers.StartServers(
		configuration.Servers,
		requestLogger,
	)

	awaitShutdownSignal()
}

func awaitShutdownSignal() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	logger.Fatalf("Signal (%v) received, stopping", s)
}
