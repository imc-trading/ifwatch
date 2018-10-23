package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/mickep76/color"
	"github.com/mickep76/encoding"
	_ "github.com/mickep76/encoding/json"
	"github.com/mickep76/log"

	"github.com/imc-trading/ifwatch/backend"
	"github.com/imc-trading/ifwatch/backend/kafka"
	"github.com/imc-trading/ifwatch/config"
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

	e := &Event{}
	if err := codec.Decode(in, e); err != nil {
		log.Errorf("decode json: %v", err)
		return
	}

	// Skip host.
	if printHost != "" && printHost != host {
		return
	}

	// Skip interface.
	if len(printInterfaces) > 0 && !inList(e.Interface.Name, printInterfaces) {
		return
	}

	switch e.Action {
	case ActionAdd:
		fmt.Printf("\n%s%s/%s%s %s(added)%s%s\n", color.White, e.Host, e.Name, color.Reset, color.Green, color.Reset, e.Interface)
	case ActionModify:
		fmt.Printf("\n%s%s/%s%s %s(modified)%s%s\n", color.White, e.Host, e.Name, color.Reset, color.Yellow, color.Reset, e.Interface)
	case ActionDelete:
		fmt.Printf("\n%s%s/%s%s %s(deleted)%s%s\n", color.White, e.Host, e.Name, color.Reset, color.Red, color.Reset, e.Interface)
	}
})

var jsonHandler = backend.MessageHandler(func(k string, in []byte) {
	codec, err := encoding.NewCodec("json", encoding.WithIndent("  "))
	if err != nil {
		log.Errorf("new codec: %v", err)
		return
	}

	e := &Event{}
	if err := codec.Decode(in, e); err != nil {
		log.Errorf("decode json: %v", err)
		return
	}

	// Skip host.
	if printHost != "" && printHost != host {
		return
	}

	// Skip interface.
	if !inList(e.Interface.Name, printInterfaces) {
		return
	}

	out, err := codec.Encode(e)
	if err != nil {
		log.Errorf("encode json: %v", err)
		return
	}

	fmt.Print(string(out))
})

func Subscribe(c *config.Config, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)

	host, _ = os.Hostname()

	var printJson bool
	fl.BoolVar(&printJson, "json", false, "Print as JSON.")
	fl.StringVar(&printHost, "host", "", "Host.")
	fl.Parse(args)

	if len(fl.Args()) > 0 {
		printInterfaces = fl.Args()
	}

	var err error
	if sub, err = kafka.NewSubscriber(c.Brokers, c.Topic); err != nil {
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
