package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/kr/pretty"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/dropprivileges"
	"github.com/aaronriekenberg/go-httpd/logging"
	"github.com/aaronriekenberg/go-httpd/requestlogging"
	"github.com/aaronriekenberg/go-httpd/servers"
)

var gitCommit string

var logger = logging.GetLogger()

func getAppName() string {
	return fmt.Sprintf(
		"%v (go version = %q gitCommit = %q)",
		os.Args[0],
		runtime.Version(),
		gitCommit,
	)
}

type commandLineFlags struct {
	configFilePath string
	verbose        bool
}

func processCommandLineFlags() commandLineFlags {
	commandLineFlags := commandLineFlags{}

	flag.StringVar(
		&commandLineFlags.configFilePath,
		"f",
		"/etc/gohttpd.json",
		"config file path",
	)

	flag.BoolVar(
		&commandLineFlags.verbose,
		"v",
		false,
		"enable verbose logging",
	)

	flag.Usage = func() {

		fmt.Fprintf(
			flag.CommandLine.Output(),
			"Usage of %v:\n",
			getAppName(),
		)

		flag.PrintDefaults()
	}

	flag.Parse()

	return commandLineFlags
}

func main() {

	commandLineFlags := processCommandLineFlags()

	logging.SetVerbose(commandLineFlags.verbose)

	logger.Printf("starting %v", getAppName())
	logger.Printf("commandLineFlags = %+v", commandLineFlags)

	configuration := config.ReadConfiguration(commandLineFlags.configFilePath)
	logger.Printf("configuration:\n%# v", pretty.Formatter(configuration))

	servers.CreateServers(configuration.Servers)

	dropprivileges.DropPrivileges(configuration.DropPrivileges)

	requestLogger := requestlogging.CreateRequestLogger(configuration.RequestLogger)

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
