package commandline

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var gitCommit string

func AppName() string {
	return fmt.Sprintf(
		"%v (go version = %q gitCommit = %q)",
		os.Args[0],
		runtime.Version(),
		gitCommit,
	)
}

type CommandLineFlags struct {
	ConfigFilePath string
	Verbose        bool
}

func ProcessCommandLineFlags() CommandLineFlags {
	commandLineFlags := CommandLineFlags{}

	flag.StringVar(
		&commandLineFlags.ConfigFilePath,
		"f",
		"/etc/gohttpd.json",
		"config file path",
	)

	flag.BoolVar(
		&commandLineFlags.Verbose,
		"v",
		false,
		"enable verbose logging",
	)

	flag.Usage = func() {

		fmt.Fprintf(
			flag.CommandLine.Output(),
			"Usage of %v:\n",
			AppName(),
		)

		flag.PrintDefaults()
	}

	flag.Parse()

	return commandLineFlags
}
