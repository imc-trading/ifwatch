package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/imc-trading/ifwatch/backend"
	"github.com/imc-trading/ifwatch/backend/etcd"
	"github.com/imc-trading/ifwatch/backend/kafka"
	"github.com/mickep76/go-difflib/difflib"
)

var sub backend.Subscriber
var ifHost string
var ifName string

var textHandler = backend.MessageHandler(func(k string, b []byte) {
	e := &Event{}
	if err := json.Unmarshal(b, e); err != nil {
		log.Fatal("unmarshal: %v", err)
	}

	if ifHost != "" && ifHost != e.Host {
		return
	}

	if ifName != "" && ifName != e.Name {
		return
	}

	switch e.Action {
	case ActionAdd:
		fmt.Printf("%s%s/%s%s %s(added)%s", white, e.Host, e.Name, clear, green, clear)
		fmt.Println(e.StringCol())
	case ActionModify:
		if v, ok := eventMap[e.Name]; ok {
			fmt.Printf("%s%s/%s%s %s(modified)%s\n", white, e.Host, e.Name, clear, yellow, clear)
			diff := difflib.UnifiedDiff{
				A:       difflib.SplitLines(v.Interface.String()),
				B:       difflib.SplitLines(e.Interface.String()),
				Context: 3,
				Colored: true,
			}
			txt, _ := difflib.GetUnifiedDiffString(diff)
			for _, l := range strings.Split(strings.Trim(txt, "\n"), "\n") {
				fmt.Printf("%s%s\n", l, clear)
			}
		} else {
			fmt.Printf("%s%s/%s%s %s(modified)%s", white, e.Host, e.Name, clear, yellow, clear)
			fmt.Println(e.StringCol())
		}
	case ActionRefresh:
		fmt.Printf("%s%s/%s%s %s(refresh)%s", white, e.Host, e.Name, clear, magenta, clear)
		fmt.Println(e.StringCol())
	case ActionDelete:
		fmt.Printf("%s%s/%s%s %s(deleted)%s", white, e.Host, e.Name, clear, red, clear)
		fmt.Println(e.StringCol())
	}

	eventMap[e.Name] = e
})

var jsonHandler = backend.MessageHandler(func(k string, b []byte) {
	e := &Event{}
	if err := json.Unmarshal(b, e); err != nil {
		log.Fatal("unmarshal: %v", err)
	}

	if ifHost != "" && ifHost != e.Host {
		return
	}

	if ifName != "" && ifName != e.Name {
		return
	}

	fmt.Println(string(e.JSONPretty()))
	eventMap[e.Name] = e
})

func Watch(args map[string]interface{}) {
	if args["<host/interface>"] != nil {
		a := strings.Split(args["<host/interface>"].(string), "/")
		ifHost = a[0]
		if len(a) > 1 {
			ifName = a[1]
		}
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

	// Populate event map
	eventMap = make(map[string]*Event)
	for _, e := range eventList {
		if v, ok := eventMap[e.Name]; ok {
			if v.Created.Before(e.Created) {
				eventMap[e.Name] = e
			}
		} else {
			eventMap[e.Name] = e
		}
	}

	if args["--json"].(bool) {
		sub.AddHandler(jsonHandler)
	} else {
		sub.AddHandler(textHandler)
	}

	if err := sub.Start(); err != nil {
		log.Fatal(err)
	}
}
