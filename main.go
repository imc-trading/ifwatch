package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mickep76/log"

	"github.com/imc-trading/ifwatch/command"
	"github.com/imc-trading/ifwatch/config"
)

var version string

func main() {
	// Define usage.
	flag.Usage = func() {
		fmt.Print(`Usage: ifwatch [OPTIONS] COMMAND

Commands:
  print         Print network interfaces.
  publish       Publish events to Kafka.
  subscribe     Subscribe to and print events from Kafka.
  history       Show events history from Kafka.

Options:
`)
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Define flags.
	var printVers, printConf bool
	var conf string
	flag.BoolVar(&printVers, "version", false, "Print application version.")
	flag.StringVar(&conf, "conf", "/etc/ifwatch.toml", "Configuration file.")
	flag.BoolVar(&printConf, "print-conf", false, "Print configuration.")
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
	}

	// Print version.
	if printVers {
		fmt.Println(version)
		os.Exit(0)
	}

	// Expand home folder prefix.
	if strings.HasPrefix(conf, "~") {
		conf = filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(conf, "~"))
	}

	// New config.
	c := config.NewConfig()

	// Read config file.
	if err := c.Load(conf); err != nil {
		log.Fatal(err)
	}

	// Print config.
	if printConf {
		c.Print()
		os.Exit(0)
	}

	// Log no color.
	if c.LogNoColor {
		log.NoColor()
	}

	// Log no date/time.
	if c.LogNoDate {
		log.SetFlags(0)
	}

	// Set log level.
	if c.LogLevel != "" {
		if err := log.SetLogLevelString(c.LogLevel); err != nil {
			log.Fatal(err)
		}
	}

	// Execute command.
	var err error
	switch flag.Args()[0] {
	case "print":
		err = command.Print(c, flag.Args()[1:])
	case "publish":
		err = command.Publish(c, flag.Args()[1:])
	case "subscribe":
		command.Subscribe(c, flag.Args()[1:])
	case "history":
		command.History(c, flag.Args()[1:])
	}
	if err != nil {
		log.Fatal(err)
	}
}
