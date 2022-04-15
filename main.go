package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kr/pretty"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/servers"
)

func awaitShutdownSignal() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received, stopping", s)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v <config json file>", os.Args[0])
	}

	configFile := os.Args[1]

	configuration := config.ReadConfiguration(configFile)
	log.Printf("configuration:\n%# v", pretty.Formatter(configuration))

	if configuration.Chroot != nil {
		log.Printf("call Chroot DirectoryPath = %q", configuration.Chroot.DirectoryPath)
		err := syscall.Chroot(configuration.Chroot.DirectoryPath)
		if err != nil {
			log.Fatalf("chroot error = %v", err)
		}
	}

	servers.StartServers(configuration.Servers)

	awaitShutdownSignal()
}
