package netlink

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

const (
	// NetlinkRoute return socket descriptor.
	NetlinkRoute = 0

	// RtmGrpLink Rtnetlink multicast group.
	RtmGrpLink = 0x1
)

// Conn provides an interface for connecting to netlink socket.
type Conn struct {
	Family     int
	Groups     uint32
	FileDescr  int
	SocketAddr *unix.SockaddrNetlink
	Pid        uint32
}

// Dial netlink socket.
func Dial(family int, groups uint32) (*Conn, error) {
	fd, err := unix.Socket(
		unix.AF_NETLINK,
		unix.SOCK_RAW,
		family,
	)
	if err != nil {
		return nil, errors.Wrap(err, "dial netlink socket")
	}

	return &Conn{
		Family:    family,
		Groups:    groups,
		FileDescr: fd,
	}, nil
}

// Bind to netlink socket.
func (c *Conn) Bind() error {
	c.SocketAddr = &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: c.Groups,
	}

	if err := unix.Bind(c.FileDescr, c.SocketAddr); err != nil {
		_ = c.Close()
		return errors.Wrap(err, "bind to netlink socket")
	}

	sa, err := unix.Getsockname(c.FileDescr)
	if err != nil {
		_ = c.Close()
		return errors.Wrap(err, "getsockname")
	}

	c.Pid = sa.(*unix.SockaddrNetlink).Pid

	return nil
}

// Close netlink socket.
func (c *Conn) Close() error {
	return unix.Close(c.FileDescr)
}

// Receive messages from netlink socket.
func (c *Conn) Receive() ([]syscall.NetlinkMessage, error) {
	b := make([]byte, os.Getpagesize())
	for {
		// Get buffer size
		n, _, _, _, err := unix.Recvmsg(c.FileDescr, b, nil, unix.MSG_PEEK)
		if err != nil {
			return nil, errors.Wrap(err, "netlink get buffer size")
		}

		// Break if buffer is large enough
		if n < len(b) {
			break
		}

		// Increase buffer size
		b = make([]byte, len(b)*2)
	}

	// Get all messages
	n, _, _, from, err := unix.Recvmsg(c.FileDescr, b, nil, 0)
	if err != nil {
		return nil, errors.Wrap(err, "netlink get all messages")
	}

	addr, ok := from.(*unix.SockaddrNetlink)
	if !ok {
		return nil, errors.New("expected unix.SockaddrNetlink but received different unix.Sockaddr")
	}
	if addr.Family != unix.AF_NETLINK {
		return nil, errors.New("received invalid netlink family")
	}

	return syscall.ParseNetlinkMessage(b[:n])
}
