package netlink

import (
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

// Interfaces connect using rtnetlink and retrieve all network interfaces.
func Interfaces() ([]Interface, error) {
	tab, err := syscall.NetlinkRIB(syscall.RTM_GETLINK, syscall.AF_UNSPEC)
	if err != nil {
		return nil, errors.Wrapf(err, "netlinkrib get interfaces")
	}

	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return nil, errors.Wrap(err, "netlink parse message")
	}

	var ift []Interface
	for _, m := range msgs {
		switch m.Header.Type {
		case syscall.NLMSG_DONE:
			break
		case syscall.RTM_NEWLINK:
			ifim := (*syscall.IfInfomsg)(unsafe.Pointer(&m.Data[0]))
			attrs, err := syscall.ParseNetlinkRouteAttr(&m)
			if err != nil {
				return nil, errors.Wrap(err, "pnetlink parse route attr")
			}
			ift = append(ift, *ParseNewLink(ifim, attrs))
		}
	}

	return ift, nil
}
