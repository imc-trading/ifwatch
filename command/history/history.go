package history

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mickep76/color"
	"github.com/mickep76/log"
	"gitlab.trading.imc.intra/go/go-difflib/difflib"

	"github.com/imc-trading/ifwatch/backend"
	"github.com/imc-trading/ifwatch/backend/kafka"
	"github.com/imc-trading/ifwatch/config"
	"github.com/imc-trading/ifwatch/event"
)

var eventList []*event.Event

var eventListHandler = backend.MessageHandler(func(k string, b []byte) {
	e := &event.Event{}
	if err := json.Unmarshal(b, e); err != nil {
		log.Errorf("unmarshal: %v", err)
		return
	}
	eventList = append(eventList, e)
})

func inList(a string, l []string) bool {
	for _, b := range l {
		if a == b {
			return true
		}
	}
	return false
}

func usage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: ifwatch history [OPTIONS] [INTERFACE...]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func History(c *config.Config, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)

	host, _ := os.Hostname()

	var printJson bool
	var printHost string
	fl.BoolVar(&printJson, "json", false, "Print as JSON.")
	fl.StringVar(&printHost, "host", host, "Host.")
	fl.Parse(args)

	printInterfaces := fl.Args()

	var err error
	sub, err := kafka.NewSubscriber(c.Brokers, c.Topic, c.ComprAlgo)
	if err != nil {
		return err
	}

	sub.Versions(printHost, eventListHandler)

	events := make(map[string][]*event.Event)
	for _, e := range eventList {
		// Skip host.
		if printHost != "" && printHost != host {
			continue
		}

		// Skip interface.
		if len(printInterfaces) > 0 && !inList(e.Interface.Name, printInterfaces) {
			continue
		}

		k := e.Host + "/" + e.Name
		if _, ok := events[k]; !ok {
			events[k] = []*event.Event{}
		}
		events[k] = append(events[k], e)
	}

	if printJson {
		b, err := json.MarshalIndent(events, "", "  ")
		if err != nil {
			log.Fatal("encode events: %v", err)
		}
		fmt.Println(string(b))
		return nil
	}

	for key, list := range events {
		fmt.Printf("%s%s%s\n│\n", color.White, key, color.Reset)

		var prev *event.Event
		for i, e := range list {
			itmPad := "├── "
			txtPad := "│   "
			if i == len(list)-1 {
				itmPad = "└── "
				txtPad = "    "
			}

			switch e.Action {
			case event.ActionAdd:
				fmt.Printf("%s%s%s%s %s(added)%s\n", itmPad, color.Blue, e.Created.Format("2006-01-02 15:04:05"), color.Reset, color.Green, color.Reset)
				for _, l := range strings.Split(strings.Trim(e.Interface.String(), "\n"), "\n") {
					fmt.Printf("%s%s\n", txtPad, l)
				}
			case event.ActionModify:
				fmt.Printf("%s%s%s%s %s(modified)%s\n", itmPad, color.Blue, e.Created.Format("2006-01-02 15:04:05"), color.Reset, color.Yellow, color.Reset)
				if prev != nil {
					diff := difflib.UnifiedDiff{
						A:       difflib.SplitLines(prev.Interface.String()),
						B:       difflib.SplitLines(e.Interface.String()),
						Context: 3,
						Colored: true,
					}
					txt, _ := difflib.GetUnifiedDiffString(diff)
					for _, l := range strings.Split(strings.Trim(txt, "\n"), "\n") {
						fmt.Printf("%s%s%s\n", txtPad, l, color.Reset)
					}
				} else {
					for _, l := range strings.Split(strings.Trim(e.Interface.String(), "\n"), "\n") {
						fmt.Printf("%s%s\n", txtPad, l)
					}
				}
			case event.ActionDelete:
				fmt.Printf("%s%s%s%s %s(deleted)%s\n", itmPad, color.Blue, e.Created.Format("2006-01-02 15:04:05"), color.Reset, color.Red, color.Reset)
			}
			prev = e
		}

		fmt.Println()
	}
	return nil
}
