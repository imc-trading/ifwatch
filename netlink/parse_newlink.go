package netlink

import (
	"encoding/json"
	"net"
	"strings"
	"syscall"
	"unsafe"
)

const (
	iffUp           = 0x1
	iffBroadcast    = 0x2
	iffDebug        = 0x4
	iffLoopback     = 0x8
	iffPointToPoint = 0x10
	iffNoTrailers   = 0x20
	iffRunning      = 0x40
	iffNoArp        = 0x80
	iffPromisc      = 0x100
	iffAllMulti     = 0x200
	iffMaster       = 0x400
	iffSlave        = 0x800
	iffMulticast    = 0x1000
	iffPortSel      = 0x2000
	iffAutoMedia    = 0x4000
	iffDynamic      = 0x8000
	iffLowerUp      = 0x10000
	iffDormant      = 0x20000
	iffEcho         = 0x40000
)

// Flags type for network interface state.
type Flags uint

const (
	// FlagUp interface is up (administratively).
	FlagUp Flags = 1 << iota

	// FlagBroadcast broadcast address valid.
	FlagBroadcast

	// FlagDebug turn on debugging.
	FlagDebug

	// FlagLoopback is a loopback net.
	FlagLoopback

	// FlagPointToPoint interface is has p-p link.
	FlagPointToPoint

	// FlagNoTrailers avoid use of trailers.
	FlagNoTrailers

	// FlagRunning interface RFC2863 OPER_UP.
	FlagRunning

	// FlagNoArp no ARP protocol.
	FlagNoArp

	// FlagPromisc receive all packets.
	FlagPromisc

	// FlagAllMulti receive all multicast packets.
	FlagAllMulti

	// FlagMaster master of a load balancer.
	FlagMaster

	// FlagSlave slave of a load balancer.
	FlagSlave

	// FlagMulticast supports multicast.
	FlagMulticast

	// FlagPortSel can set media type.
	FlagPortSel

	// FlagAutoMedia auto media select active.
	FlagAutoMedia

	// FlagDynamic dialup device with changing addresses.
	FlagDynamic

	// FlagLowerUp driver signals L1 up.
	FlagLowerUp

	// FlagDormant driver signals dormant.
	FlagDormant

	// FlagEcho echo sent packets.
	FlagEcho
)

var flagNames = []string{
	"up",
	"broadcast",
	"debug",
	"loopback",
	"pointopoint",
	"notrailers",
	"running",
	"noarp",
	"promisc",
	"allmulti",
	"master",
	"slave",
	"multicast",
	"portsel",
	"automedia",
	"dynamic",
	"lower_up",
	"dormant",
	"echo",
}

// HwAddr hardware address type.
type HwAddr []byte

const hexDigit = "0123456789abcdef"

const (
	// See linux/if_arp.h.
	// Note that Linux doesn't support IPv4 over IPv6 tunneling.
	sysARPHardwareIPv4IPv4 = 768 // IPv4 over IPv4 tunneling.
	sysARPHardwareIPv6IPv6 = 769 // IPv6 over IPv6 tunneling.
	sysARPHardwareIPv6IPv4 = 776 // IPv6 over IPv4 tunneling.
	sysARPHardwareGREIPv4  = 778 // Any over GRE over IPv4 tunneling.
	sysARPHardwareGREIPv6  = 823 // Any over GRE over IPv6 tunneling.
)

// Interface provides information about a network interface.
type Interface struct {
	Index        int            `json:"index"`
	MTU          int            `json:"mtu"`
	Name         string         `json:"name"`
	HwAddr       HwAddr         `json:"hwaddr,omitempty"`
	Flags        Flags          `json:"flags"`
	NetInterface *net.Interface `json:"-"`
}

// String return a string of all flags.
func (f Flags) String() string {
	return strings.Join(f.Slice(), "|")
}

// Slice return a list of all flags.
func (f Flags) Slice() []string {
	var l []string
	for i, name := range flagNames {
		if f&(1<<uint(i)) != 0 {
			l = append(l, name)
		}
	}
	return l
}

// MarshalJSON marshal flags into JSON.
func (f Flags) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Slice())
}

func parseFlags(flags uint32) Flags {
	var f Flags
	if flags&iffUp != 0 {
		f |= FlagUp
	}
	if flags&iffBroadcast != 0 {
		f |= FlagBroadcast
	}
	if flags&iffDebug != 0 {
		f |= FlagDebug
	}
	if flags&iffLoopback != 0 {
		f |= FlagLoopback
	}
	if flags&iffPointToPoint != 0 {
		f |= FlagPointToPoint
	}
	if flags&iffNoTrailers != 0 {
		f |= FlagNoTrailers
	}
	if flags&iffRunning != 0 {
		f |= FlagRunning
	}
	if flags&iffNoArp != 0 {
		f |= FlagNoArp
	}
	if flags&iffPromisc != 0 {
		f |= FlagPromisc
	}
	if flags&iffAllMulti != 0 {
		f |= FlagAllMulti
	}
	if flags&iffMaster != 0 {
		f |= FlagMaster
	}
	if flags&iffSlave != 0 {
		f |= FlagSlave
	}
	if flags&iffMulticast != 0 {
		f |= FlagMulticast
	}
	if flags&iffPortSel != 0 {
		f |= FlagPortSel
	}
	if flags&iffAutoMedia != 0 {
		f |= FlagAutoMedia
	}
	if flags&iffDynamic != 0 {
		f |= FlagDynamic
	}
	if flags&iffLowerUp != 0 {
		f |= FlagLowerUp
	}
	if flags&iffDormant != 0 {
		f |= FlagDormant
	}
	if flags&iffEcho != 0 {
		f |= FlagEcho
	}
	return f
}

func parseNetFlags(rawFlags uint32) net.Flags {
	var f net.Flags
	if rawFlags&iffUp != 0 {
		f |= net.FlagUp
	}
	if rawFlags&iffBroadcast != 0 {
		f |= net.FlagBroadcast
	}
	if rawFlags&iffLoopback != 0 {
		f |= net.FlagLoopback
	}
	if rawFlags&iffPointToPoint != 0 {
		f |= net.FlagPointToPoint
	}
	if rawFlags&iffMulticast != 0 {
		f |= net.FlagMulticast
	}
	return f
}

func (a HwAddr) String() string {
	if len(a) == 0 {
		return ""
	}
	buf := make([]byte, 0, len(a)*3-1)
	for i, b := range a {
		if i > 0 {
			buf = append(buf, ':')
		}
		buf = append(buf, hexDigit[b>>4])
		buf = append(buf, hexDigit[b&0xF])
	}
	return string(buf)
}

// MarshalJSON marshal hardware address into JSON.
func (a HwAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// ParseNewLink parse interface info message.
func ParseNewLink(ifim *syscall.IfInfomsg, attrs []syscall.NetlinkRouteAttr) *Interface {
	i := Interface{
		Index:        int(ifim.Index),
		Flags:        parseFlags(ifim.Flags),
		NetInterface: &net.Interface{Index: int(ifim.Index), Flags: parseNetFlags(ifim.Flags)},
	}

	for _, a := range attrs {
		switch a.Attr.Type {
		case syscall.IFLA_ADDRESS:
			// We never return any /32 or /128 IP address
			// prefix on any IP tunnel interface as the
			// hardware address.
			switch len(a.Value) {
			case net.IPv4len:
				switch ifim.Type {
				case sysARPHardwareIPv4IPv4, sysARPHardwareGREIPv4, sysARPHardwareIPv6IPv4:
					continue
				}
			case net.IPv6len:
				switch ifim.Type {
				case sysARPHardwareIPv6IPv6, sysARPHardwareGREIPv6:
					continue
				}
			}
			var nonzero bool
			for _, b := range a.Value {
				if b != 0 {
					nonzero = true
					break
				}
			}
			if nonzero {
				i.HwAddr = a.Value[:]
				i.NetInterface.HardwareAddr = a.Value[:]
			}
		case syscall.IFLA_IFNAME:
			i.Name = string(a.Value[:len(a.Value)-1])
			i.NetInterface.Name = i.Name
		case syscall.IFLA_MTU:
			i.MTU = int(*(*uint32)(unsafe.Pointer(&a.Value[:4][0])))
			i.NetInterface.MTU = i.MTU
		}
	}
	return &i
}
