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
)

func refreshInterface(i *netx.Interface) {
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
		Created:   time.Now(),
		Host:      host,
		Action:    ActionRefresh,
		Interface: i,
	}

	log.Printf("send %s event for interface [%s]: %s%s", e.Action, k, e.Name, i.String())
	pub.Send(k, e)
}

func Refresh(args map[string]interface{}) {
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
