package main

import (
	"fmt"
	"github.com/1xyz/jellybeans-workload/workload"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	sizeAvg   = float64(1024 * 32)
	sizeStdev = float64(1024 * 10)
	version   = "0.1"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	usage := `usage: jellybeans-workload [--version] [--help] <command> [<args>...]
options:
   -h, --help
   --verbose      Change the logging level verbosity
The commands are:
   consumer   Run a consumer workload
   producer   Run a producer workload
See 'jellybeans-workload <command> --help' for more information on a specific command.
`
	parser := &docopt.Parser{OptionsFirst: true}
	args, err := parser.ParseArgs(usage, nil, version)
	if err != nil {
		log.Errorf("error = %v", err)
		os.Exit(1)
	}

	cmd := args["<command>"].(string)
	cmdArgs := args["<args>"].([]string)
	fmt.Println("global arguments:", args)
	fmt.Println("command arguments:", cmd, cmdArgs)

	if err := RunCommand(cmd, cmdArgs, version); err != nil {
		log.Errorf("error %v", err)
		os.Exit(1)
	}
}

func RunCommand(c string, args []string, version string) error {
	argv := append([]string{c}, args...)
	switch c {
	case "consumer":
		return workload.CmdConsumer(argv, version)
	case "producer":
		return workload.CmdProducer(argv, version)
	default:
		return fmt.Errorf("runCommand: %s is not a supported command. See 'coolbeans help'", c)
	}
}
