package command

import (
	"flag"
	"os"
	"time"

	"github.com/imc-trading/ifwatch/backend"
	"github.com/imc-trading/ifwatch/backend/kafka"
	"github.com/imc-trading/ifwatch/config"
	"github.com/imc-trading/ifwatch/network"

	"github.com/mickep76/log"
)

var pub backend.Publisher
var published = map[string]bool{}
var conf *config.Config

func eventHandler(i *network.Interface, f network.Flag) {
	if inList(i.Name, conf.SkipInterfaces) {
		log.Warnf("skip interface: %s", i.Name)
		return
	}

	if inList(i.Driver, conf.SkipDrivers) {
		log.Warnf("skip driver: %s for interface: %s", i.Driver, i.Name)
		return
	}

	e := &Event{
		Created:   time.Now(),
		Host:      host,
		Interface: i,
	}

	if f == network.FlagDelete {
		e.Action = ActionDelete
	} else {
		if published[i.Name] {
			e.Action = ActionModify
		} else {
			e.Action = ActionAdd
		}
	}

	log.Infof("send action: %s event for interface: %s", e.Action, e.Name)
	published[i.Name] = true
	pub.Send(host, e)
}

func Publish(c *config.Config, args []string) error {
	conf = c

	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)
	fl.Parse(args)

	host, _ = os.Hostname()

	var err error
	pub, err = kafka.NewPublisher(c.Brokers, c.Topic)
	if err != nil {
		return err
	}

	// Get host interfaces
	interfaces, err := network.Interfaces()
	if err != nil {
		return err
	}

	// Send interfaces
	for _, i := range interfaces {
		eventHandler(i, network.FlagNew)
	}

	// Start watcher for interface events
	w := network.NewWatcher()
	w.AddHandler(eventHandler)

	return w.Start()
}
