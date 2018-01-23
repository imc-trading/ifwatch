package cmd

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/imc-trading/ifwatch/backend"
	"github.com/imc-trading/ifwatch/backend/etcd"
	"github.com/imc-trading/ifwatch/backend/kafka"
	"github.com/imc-trading/ifwatch/netx"

	"github.com/mickep76/go-difflib/difflib"
)

var skipInterfaces []string
var skipDrivers []string
var host string
var prefix string
var eventMap EventMap
var pub backend.Publisher
var refresh int

func inList(v string, l []string) bool {
	for _, i := range l {
		if v == i {
			return true
		}
	}
	return false
}

func interfaceHandler(i *netx.Interface, f netx.Flag) {
	if inList(i.Name, skipInterfaces) {
		return
	}

	if inList(i.Driver, skipDrivers) {
		return
	}

	var k string
	switch pub.Backend() {
	case backend.Kafka:
		k = host
	case backend.Etcd:
		if prefix != "" {
			k += prefix + "/"
		}
		k += host + "/" + i.Name
	}

	e := &Event{
		Created: time.Now(),
		Host:    host,
	}

	var s string
	if f == netx.FlagDelete {
		e.Action = ActionDelete
		e.Interface = &netx.Interface{Name: i.Name}
	} else {
		if v, ok := eventMap[i.Name]; ok {
			// Not modified, ignore
			if i.String() == v.Interface.String() {
				log.Printf("ignore %s event, not modified for interface: %s", ActionModify, i.Name)
				return
			}

			// Modified
			e.Action = ActionModify

			diff := difflib.UnifiedDiff{
				A:       difflib.SplitLines(v.Interface.String()),
				B:       difflib.SplitLines(i.String()),
				Context: 3,
			}
			s, _ = difflib.GetUnifiedDiffString(diff)
			s = "\n" + s
		} else {
			// Add
			e.Action = ActionAdd
			s = i.String()
		}

		e.Interface = i
	}

	log.Printf("send %s event for interface [%s]: %s%s", e.Action, k, e.Name, s)
	eventMap[i.Name] = e
	pub.Send(k, e)
}

func Daemon(args map[string]interface{}) {
	skipInterfaces = strings.Split(args["--skip-interfaces"].(string), ",")
	skipDrivers = strings.Split(args["--skip-drivers"].(string), ",")

	host, _ = os.Hostname()

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
			key += prefix + "/"
		}
		key += host

		if pub, err = etcd.NewPublisher(endpoints, tmout, prefix); err != nil {
			log.Fatal(err)
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
		key = host

		var err error
		if pub, err = kafka.NewPublisher(brokers, topic); err != nil {
			log.Fatal(err)
		}

		if sub, err = kafka.NewSubscriber(brokers, topic); err != nil {
			log.Fatal(err)
		}
	}

	if args["--refresh"] != nil {
		var err error
		refresh, err = strconv.Atoi(args["--refresh"].(string))
		if err != nil {
			log.Fatalf("strconv: %v", err)
		}

		log.Printf("start ticker for refresh event every %d second", refresh)
		ticker := time.NewTicker(time.Duration(refresh) * time.Second)
		go func() {
			for range ticker.C {
				// Get host interfaces
				interfaces, err := netx.Interfaces()
				if err != nil {
					log.Fatal(err)
				}

				// Send refresh interfaces
				for _, i := range interfaces {
					refreshInterface(i)
				}
			}
		}()
	}

	// Get interfaces from backend
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

	// Get host interfaces
	interfaces, err := netx.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	// Send new/modified interfaces
	for _, i := range interfaces {
		interfaceHandler(i, netx.FlagNew)
	}

	// Send delete interfaces
	for _, e := range eventList {
		match := false
		for _, i := range interfaces {
			if e.Name == i.Name {
				match = true
				break
			}
		}
		if !match {
			interfaceHandler(e.Interface, netx.FlagDelete)
		}
	}

	// Start watcher for interface events
	w := netx.NewWatcher()
	w.AddHandler(interfaceHandler)

	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
}
