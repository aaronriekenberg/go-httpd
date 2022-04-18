package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/kr/pretty"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/dropprivileges"
	"github.com/aaronriekenberg/go-httpd/servers"
)

var gitCommit string

func awaitShutdownSignal() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received, stopping", s)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Printf("go version = %q gitCommit = %q", runtime.Version(), gitCommit)

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v <config json file>", os.Args[0])
	}

	configFile := os.Args[1]

	configuration := config.ReadConfiguration(configFile)
	log.Printf("configuration:\n%# v", pretty.Formatter(configuration))

	dropprivileges.DropPrivileges(configuration.DropPrivileges)

	servers.StartServers(configuration.Servers)

	awaitShutdownSignal()
}
