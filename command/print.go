package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/mickep76/color"
	"github.com/mickep76/encoding"
	_ "github.com/mickep76/encoding/json"

	"github.com/imc-trading/ifwatch/config"
	"github.com/imc-trading/ifwatch/network"
)

func usage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: ifwatch print [OPTIONS] [INTERFACE...]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func Print(c *config.Config, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)

	var printJson bool
	fl.BoolVar(&printJson, "json", false, "Print as JSON.")
	fl.Parse(args)

	var printInterfaces []string
	if len(fl.Args()) > 0 {
		printInterfaces = fl.Args()
	}

	interfaces, err := network.Interfaces()
	if err != nil {
		return err
	}

	codec, err := encoding.NewCodec("json", encoding.WithIndent("  "))
	if err != nil {
		return err
	}

	host, _ := os.Hostname()
	for _, i := range interfaces {
		if len(printInterfaces) != 0 && !inList(i.Name, printInterfaces) {
			continue
		}

		if printJson {
			b, _ := codec.Encode(i)
			fmt.Print(string(b))
		} else {
			fmt.Printf("\n%s%s/%s%s\n%s", color.White, host, i.Name, color.Reset, i)
		}
	}
	return nil
}
