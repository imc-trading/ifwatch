package network

import (
	"syscall"
	"unsafe"

	"github.com/pkg/errors"

	"github.com/imc-trading/ifwatch/netlink"
)

// Flag for event type
type Flag uint

const (
	// FlagNew new interface event
	FlagNew Flag = iota

	// FlagDelete delete interface event
	FlagDelete
)

// Handler for netlink events
type Handler func(*Interface, Flag)

// Watcher for netlink events
type Watcher struct {
	conn     *netlink.Conn
	Handlers []Handler
}

// NewWatcher constructor
func NewWatcher() *Watcher {
	return &Watcher{}
}

// AddHandler for netlink events
func (w *Watcher) AddHandler(h func(*Interface, Flag)) {
	w.Handlers = append(w.Handlers, Handler(h))
}

// Start watcher for netlink events
func (w *Watcher) Start() error {
	conn, err := netlink.Dial(netlink.NetlinkRoute, netlink.RtmGrpLink)
	if err != nil {
		return err
	}

	if err := conn.Bind(); err != nil {
		return err
	}

	defer conn.Close()

	for {
		msgs, err := conn.Receive()
		if err != nil {
			//			log.Fatal(err)
			return errors.Wrap(err, "conn receive")
		}

		for _, m := range msgs {
			switch m.Header.Type {
			case syscall.NLMSG_DONE:
				break
			case syscall.RTM_NEWLINK:
				ifim := (*syscall.IfInfomsg)(unsafe.Pointer(&m.Data[0]))
				attrs, err := syscall.ParseNetlinkRouteAttr(&m)
				if err != nil {
					return errors.Wrap(err, "parse netlink route attr")
				}

				ni := netlink.ParseNewLink(ifim, attrs)
				i, err := ParseInterface(ni, FlagNew)
				if err != nil {
					return err
				}

				for _, handler := range w.Handlers {
					handler(i, FlagNew)
				}
			case syscall.RTM_DELLINK:
				ifim := (*syscall.IfInfomsg)(unsafe.Pointer(&m.Data[0]))
				attrs, err := syscall.ParseNetlinkRouteAttr(&m)
				if err != nil {
					return errors.Wrap(err, "parse netlink route attr")
				}

				ni := netlink.ParseNewLink(ifim, attrs)
				i, err := ParseInterface(ni, FlagDelete)
				if err != nil {
					return err
				}

				for _, handler := range w.Handlers {
					handler(i, FlagDelete)
				}
			}
		}
	}
}
