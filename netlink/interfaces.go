package netlink

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Interfaces connect using rtnetlink and retrieve all network interfaces.
func Interfaces() ([]Interface, error) {
	tab, err := syscall.NetlinkRIB(syscall.RTM_GETLINK, syscall.AF_UNSPEC)
	if err != nil {
		return nil, fmt.Errorf("netlink rib: %v", err)
	}

	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return nil, fmt.Errorf("parse netlink message: %v", err)
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
				return nil, fmt.Errorf("parse netlink route attr: %v", err)
			}
			ift = append(ift, *ParseNewLink(ifim, attrs))
		}
	}

	return ift, nil
}
