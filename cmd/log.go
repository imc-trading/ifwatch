package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/imc-trading/ifwatch/backend"
	"github.com/imc-trading/ifwatch/backend/etcd"
	"github.com/imc-trading/ifwatch/backend/kafka"

	"github.com/mickep76/go-difflib/difflib"
)

const (
	red     = "\x1b[31m"
	green   = "\x1b[32m"
	yellow  = "\x1b[33m"
	blue    = "\x1b[34m"
	magenta = "\x1b[35m"
	cyan    = "\x1b[36m"
	white   = "\x1b[37m"
	clear   = "\x1b[0m"
)

var eventList []*Event
var eventByHostIf map[string][]*Event

var addToEventListHandler = backend.MessageHandler(func(k string, b []byte) {
	e := &Event{}
	if err := json.Unmarshal(b, e); err != nil {
		log.Fatal("unmarshal: %v", err)
	}

	eventList = append(eventList, e)
})

func Log(args map[string]interface{}) {
	if args["<host/interface>"] != nil {
		a := strings.Split(args["<host/interface>"].(string), "/")
		ifHost = a[0]
		if len(a) > 1 {
			ifName = a[1]
		}
	}

	if ifHost == "" {
		ifHost, _ = os.Hostname()
	}

	var key string
	if args["--backend"].(string) == "etcd" {
		if args["--endpoints"] == nil {
			log.Fatalf("missing option: --endpoints")
		}
		endpoints := strings.Split(args["--endpoints"].(string), ",")

		tmout, err := strconv.Atoi(args["--timeout"].(string))
		if err != nil {
			log.Fatalf("strconv: %v", err)
		}

		prefix = args["--prefix"].(string)
		if prefix != "" {
			key = prefix + "/" + ifHost
		}
		if ifName != "" {
			key += "/" + ifName
		}

		if sub, err = etcd.NewSubscriber(endpoints, tmout, prefix); err != nil {
			log.Fatal(err)
		}
	} else {
		if args["--brokers"] == nil {
			log.Fatalf("missing option: --brokers")
		}
		brokers := strings.Split(args["--brokers"].(string), ",")

		topic := args["--topic"].(string)

		key = ifHost

		var err error
		if sub, err = kafka.NewSubscriber(brokers, topic); err != nil {
			log.Fatal(err)
		}
	}

	sub.Versions(key, addToEventListHandler)

	// Arrange events by host and interface
	eventByHostIf = make(map[string][]*Event)
	for _, e := range eventList {
		if ifName != "" && ifName != e.Name {
			continue
		}

		k := e.Host + "/" + e.Name
		if _, ok := eventByHostIf[k]; !ok {
			eventByHostIf[k] = []*Event{}
		}
		eventByHostIf[k] = append(eventByHostIf[k], e)
	}

	if args["--json"].(bool) {
		b, err := json.MarshalIndent(eventByHostIf, "", "  ")
		if err != nil {
			log.Fatal("encode events: %v", err)
		}

		fmt.Println(string(b))
		return
	}

	// Print per host/interface
	for key, list := range eventByHostIf {
		fmt.Printf("%s%s%s\n│\n", white, key, clear)

		var prev *Event
		for i, e := range list {
			itmPad := "├── "
			txtPad := "│   "
			if i == len(list)-1 {
				itmPad = "└── "
				txtPad = "    "
			}

			switch e.Action {
			case ActionAdd:
				fmt.Printf("%s%s%s%s %s(added)%s\n", itmPad, blue, e.Created.Format("2006-01-02 15:04:05"), clear, green, clear)
				for _, l := range strings.Split(strings.Trim(e.Interface.StringCol(), "\n"), "\n") {
					fmt.Printf("%s%s\n", txtPad, l)
				}
				prev = e
			case ActionRefresh:
				fmt.Printf("%s%s%s%s %s(refresh)%s\n", itmPad, blue, e.Created.Format("2006-01-02 15:04:05"), clear, magenta, clear)
				if prev != nil {
					diff := difflib.UnifiedDiff{
						A:       difflib.SplitLines(prev.Interface.String()),
						B:       difflib.SplitLines(e.Interface.String()),
						Context: 3,
						Colored: true,
					}
					txt, _ := difflib.GetUnifiedDiffString(diff)
					for _, l := range strings.Split(strings.Trim(txt, "\n"), "\n") {
						fmt.Printf("%s%s%s\n", txtPad, l, clear)
					}
				} else {
					for _, l := range strings.Split(strings.Trim(e.Interface.StringCol(), "\n"), "\n") {
						fmt.Printf("%s%s\n", txtPad, l)
					}
				}
				prev = e
			case ActionModify:
				fmt.Printf("%s%s%s%s %s(modified)%s\n", itmPad, blue, e.Created.Format("2006-01-02 15:04:05"), clear, yellow, clear)
				if prev != nil {
					diff := difflib.UnifiedDiff{
						A:       difflib.SplitLines(prev.Interface.String()),
						B:       difflib.SplitLines(e.Interface.String()),
						Context: 3,
						Colored: true,
					}
					txt, _ := difflib.GetUnifiedDiffString(diff)
					for _, l := range strings.Split(strings.Trim(txt, "\n"), "\n") {
						fmt.Printf("%s%s%s\n", txtPad, l, clear)
					}
				} else {
					for _, l := range strings.Split(strings.Trim(e.Interface.StringCol(), "\n"), "\n") {
						fmt.Printf("%s%s\n", txtPad, l)
					}
				}
				prev = e
			case ActionDelete:
				fmt.Printf("%s%s%s%s %s(deleted)%s\n", itmPad, blue, e.Created.Format("2006-01-02 15:04:05"), clear, red, clear)
			}
		}
	}
}
