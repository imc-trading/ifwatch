package subscribe

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mickep76/color"
	"github.com/mickep76/encoding"
	_ "github.com/mickep76/encoding/json"
	"github.com/mickep76/log"

	"github.com/imc-trading/ifwatch/backend"
	"github.com/imc-trading/ifwatch/backend/kafka"
	"github.com/imc-trading/ifwatch/config"
	"github.com/imc-trading/ifwatch/event"
)

var sub backend.Subscriber
var host string
var printHost string
var printInterfaces []string

var textHandler = backend.MessageHandler(func(k string, in []byte) {
	codec, err := encoding.NewCodec("json", encoding.WithIndent("  "))
	if err != nil {
		log.Errorf("new codec: %v", err)
		return
	}

	e := &event.Event{}
	if err := codec.Decode(in, e); err != nil {
		log.Errorf("decode json: %v", err)
		return
	}

	// Skip host.
	if printHost != "" && printHost != host {
		log.Infof("skip event for: %s/%s (%s) not matching host: %s", e.Host, e.Interface.Name, e.Action, printHost)
		return
	}

	// Skip interface.
	if len(printInterfaces) > 0 && !inList(e.Interface.Name, printInterfaces) {
		log.Infof("skip event for: %s/%s (%s) not matching interface(s): %s", e.Host, e.Interface.Name, e.Action, strings.Join(printInterfaces, ","))
		return
	}

	switch e.Action {
	case event.ActionAdd:
		fmt.Printf("%s%s/%s%s %s(added)%s%s\n\n", color.White, e.Host, e.Name, color.Reset, color.Green, color.Reset, e.Interface)
	case event.ActionModify:
		fmt.Printf("%s%s/%s%s %s(modified)%s%s\n\n", color.White, e.Host, e.Name, color.Reset, color.Yellow, color.Reset, e.Interface)
	case event.ActionDelete:
		fmt.Printf("%s%s/%s%s %s(deleted)%s%s\n\n", color.White, e.Host, e.Name, color.Reset, color.Red, color.Reset, e.Interface)
	}
})

var jsonHandler = backend.MessageHandler(func(k string, in []byte) {
	codec, err := encoding.NewCodec("json", encoding.WithIndent("  "))
	if err != nil {
		log.Errorf("new codec: %v", err)
		return
	}

	e := &event.Event{}
	if err := codec.Decode(in, e); err != nil {
		log.Errorf("decode json: %v", err)
		return
	}

	// Skip host.
	if printHost != "" && printHost != host {
		log.Infof("skip event for: %s/%s (%s) not matching host: %s", e.Host, e.Interface.Name, e.Action, printHost)
		return
	}

	// Skip interface.
	if len(printInterfaces) > 0 && !inList(e.Interface.Name, printInterfaces) {
		log.Infof("skip event for: %s/%s (%s) not matching interface(s): %s", e.Host, e.Interface.Name, e.Action, strings.Join(printInterfaces, ","))
		return
	}

	out, err := codec.Encode(e)
	if err != nil {
		log.Errorf("encode json: %v", err)
		return
	}

	fmt.Print(string(out))
})

func usage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: ifwatch subscribe [OPTIONS] [INTERFACE...]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func inList(a string, l []string) bool {
	for _, b := range l {
		if a == b {
			return true
		}
	}
	return false
}

func Subscribe(c *config.Config, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)

	host, _ = os.Hostname()

	var printJson bool
	fl.BoolVar(&printJson, "json", false, "Print as JSON.")
	fl.StringVar(&printHost, "host", "", "Host.")
	fl.Parse(args)

	printInterfaces = fl.Args()

	var err error
	if sub, err = kafka.NewSubscriber(c.Brokers, c.Topic, c.ComprAlgo); err != nil {
		log.Fatal(err)
	}

	if printJson {
		sub.AddHandler(jsonHandler)
	} else {
		sub.AddHandler(textHandler)
	}

	if err := sub.Start(); err != nil {
		return err
	}

	return nil
}
