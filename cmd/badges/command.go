package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/seriousben/badges/internal"
)

type badgesCommand struct {
	fs *flag.FlagSet

	VersionFlag bool
	HelpFlag    bool
	Port        int

	version string
}

func newBadgesCommand(version string) *badgesCommand {
	c := &badgesCommand{
		fs:      flag.NewFlagSet("", flag.ContinueOnError),
		version: version,
	}

	c.fs.Usage = c.Usage

	c.fs.BoolVar(&c.VersionFlag, "version", false, "print version")
	c.fs.BoolVar(&c.HelpFlag, "help", false, "print help")
	c.fs.IntVar(&c.Port, "port", 8123, "TCP port to listen on")

	return c
}

func (c *badgesCommand) Usage() {
	usage := `Usage: badges [--version] [--help] [flags...] serve --port <port>

	Arguments:
		serve
			Start HTTP server.

	Flags:
		--version
			display go-patch-cover version.

		--help
			display this help message.

		-port int
			TCP port to listen on.
	`

	_, _ = fmt.Fprint(os.Stdout, usage)
}

func (c *badgesCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return fmt.Errorf("flag parse error: %v", err)
	}

	if c.HelpFlag {
		c.fs.Usage()
		return nil
	}

	if c.VersionFlag {
		fmt.Println(c.version)
		return nil
	}

	switch cmd := c.fs.Arg(0); cmd {
	case "serve":
		// serve logic
		if c.Port == 0 {
			return fmt.Errorf("-port flag required")
		}
		log.Println("Starting server")
		if err := internal.Serve(c.Port); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return nil
}
