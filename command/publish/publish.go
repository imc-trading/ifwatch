package publish

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/imc-trading/ifwatch/backend"
	"github.com/imc-trading/ifwatch/backend/kafka"
	"github.com/imc-trading/ifwatch/config"
	"github.com/imc-trading/ifwatch/event"
	"github.com/imc-trading/ifwatch/network"

	"github.com/mickep76/log"
)

type Rate struct {
	Start     time.Time
	Count     int
	LimitStop bool
	Hash      [16]byte
}

var pub backend.Publisher
var conf *config.Config
var host string
var rates = map[string]*Rate{}

func eventHandler(i *network.Interface, f network.Flag) {
	if inList(i.Name, conf.SkipInterfaces) {
		log.Warnf("skip interface: %s", i.Name)
		return
	}

	if inList(i.Driver, conf.SkipDrivers) {
		log.Warnf("skip driver: %s for interface: %s", i.Driver, i.Name)
		return
	}

	if _, ok := rates[i.Name]; !ok {
		rates[i.Name] = &Rate{Start: time.Now(), Count: 0}
	}

	dur := time.Since(rates[i.Name].Start)
	rate := float64(rates[i.Name].Count) / dur.Seconds()
	if rate > float64(conf.RateLimit) {
		log.Debugf("skip event for interface: %s rate limit reached: %d rate: %.02f", i.Name, conf.RateLimit, rate)
		if !rates[i.Name].LimitStop {
			// Send stop event
			log.Warnf("send event: %s for interface: %s", event.ActionLimitStop, i.Name)
			pub.Send(host, &event.Event{
				Created:   time.Now(),
				Action:    event.ActionLimitStop,
				Host:      host,
				Interface: &network.Interface{Name: i.Name},
			})
			rates[i.Name].LimitStop = true
		}
		rates[i.Name].Count++
		return
	}

	// Send resume event
	if rates[i.Name].LimitStop {
		log.Warnf("send event: %s for interface: %s", event.ActionLimitResume, i.Name)
		pub.Send(host, &event.Event{
			Created:   time.Now(),
			Action:    event.ActionLimitResume,
			Host:      host,
			Interface: &network.Interface{Name: i.Name},
		})
		rates[i.Name].LimitStop = false
	}

	b, _ := json.Marshal(i)
	newHash := md5.Sum(b)

	e := &event.Event{
		Created:   time.Now(),
		Host:      host,
		Interface: i,
	}

	if f == network.FlagDelete {
		e.Action = event.ActionDelete
	} else {
		v, _ := rates[i.Name]
		if v.Count == 0 {
			e.Action = event.ActionAdd
		} else if newHash == v.Hash {
			e.Action = event.ActionRefresh
		} else {
			e.Action = event.ActionModify
		}
	}

	log.Debugf("send event: %s for interface: %s rate: %.02f / sec.", e.Action, e.Name, rate)
	rates[i.Name].Count++
	rates[i.Name].Hash = newHash
	pub.Send(host, e)
}

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
		fmt.Fprintf(os.Stderr, "Usage: ifwatch publish\n")
		fl.PrintDefaults()
	}
}

func refresh() error {
	// Get host interfaces
	interfaces, err := network.Interfaces()
	if err != nil {
		return err
	}

	// Send interfaces
	for _, i := range interfaces {
		eventHandler(i, network.FlagNew)
	}

	return nil
}

func Publish(c *config.Config, args []string) error {
	conf = c

	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)
	fl.Parse(args)

	host, _ = os.Hostname()

	var err error
	if pub, err = kafka.NewPublisher(c.Brokers, c.Topic, c.ComprAlgo); err != nil {
		return err
	}

	if err := refresh(); err != nil {
		return err
	}

	if c.Refresh > 0 {
		go func() {
			for {
				time.Sleep(time.Duration(c.Refresh) * time.Hour)
				if err := refresh(); err != nil {
					log.Error(err)
				}
			}
		}()
	}

	// Start watcher for interface events
	w := network.NewWatcher()
	w.AddHandler(eventHandler)

	return w.Start()
}
