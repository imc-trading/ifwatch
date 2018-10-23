package network

import (
	"errors"
	"fmt"
	"net"

	"github.com/mickep76/log"

	"github.com/imc-trading/ifwatch/ethtool"
	"github.com/imc-trading/ifwatch/module"
	"github.com/imc-trading/ifwatch/netlink"
)

// Interface provides information about system network interfaces.
type Interface struct {
	Index      int            `json:"index"`
	Slot       string         `json:"slot,omitempty"`
	Driver     string         `json:"driver,omitempty"`
	MTU        int            `json:"mtu"`
	Name       string         `json:"name"`
	DNSNames   []string       `json:"dnsNames"`
	HwAddr     string         `json:"hwaddr,omitempty"`
	PermHwAddr string         `json:"permHwAddr,omitempty"`
	Flags      []string       `json:"flags"`
	IPv4       string         `json:"ipv4,omitempty"`
	Netmask    string         `json:"netmask,omitempty"`
	Network    string         `json:"network,omitempty"`
	IPv6       string         `json:"ipv6,omitempty"`
	Module     *module.Module `json:"module,omitempty"`
}

type InterfaceList []*Interface
type InterfaceMap map[string]*Interface

// ErrOpNotSupp operation not supported.
var ErrOpNotSupp = errors.New("operation not supported")

func uitoa(val uint) string {
	var buf [32]byte // big enough for int64
	i := len(buf) - 1
	for val >= 10 {
		buf[i] = byte(val%10 + '0')
		i--
		val /= 10
	}
	buf[i] = byte(val + '0')
	return string(buf[i:])
}

func maskToDec(m net.IPMask) string {
	if len(m) == net.IPv4len {
		return uitoa(uint(m[0])) + "." +
			uitoa(uint(m[1])) + "." +
			uitoa(uint(m[2])) + "." +
			uitoa(uint(m[3]))
	}
	return ""
}

// ParseInterface parse network interface from netlink.
func ParseInterface(ni *netlink.Interface, f Flag) (*Interface, error) {
	e, err := ethtool.NewEthtool()
	if err != nil && err != ErrOpNotSupp {
		return nil, fmt.Errorf("new ethtool: %v", err)
	}

	i := &Interface{
		Index:  ni.Index,
		MTU:    ni.MTU,
		Name:   ni.Name,
		HwAddr: ni.HwAddr.String(),
		Flags:  ni.Flags.Slice(),
	}

	if ni.Name != "lo" && f != FlagDelete {
		slot, _ := e.BusInfo(ni.Name)
		i.Slot = slot

		driver, _ := e.DriverName(ni.Name)
		i.Driver = driver

		i.PermHwAddr, _ = e.PermAddr(ni.Name)

		eeprom, _ := e.ModuleEeprom(ni.Name)

		if eeprom != nil {
			if m, err := module.Decode(eeprom); err == nil {
				i.Module = m
			} else {
				log.Warn(err)
			}
		}
	}

	addrs, err := ni.NetInterface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("get interface address for %s: %v", ni.Name, err)
	}

	for _, a := range addrs {
		ip, ipNet, err := net.ParseCIDR(a.String())
		if err != nil {
			return nil, fmt.Errorf("parse CIDR for %s: %v", ni.Name, err)
		}

		if ip.To4() != nil {
			i.IPv4 = ip.String()
			s := *ipNet
			i.Netmask = maskToDec(s.Mask)
			i.Network = ipNet.String()
		} else if ip.To16() != nil {
			i.IPv6 = ip.String()
		}
	}

	if i.IPv4 != "" {
		dnsNames, err := net.LookupAddr(i.IPv4)
		if err == nil {
			for _, name := range dnsNames {
				i.DNSNames = append(i.DNSNames, name)
			}
		}
	}

	return i, nil
}

// Interfaces return list of interfaces
func Interfaces() (InterfaceList, error) {
	interfaces, err := netlink.Interfaces()
	if err != nil {
		return nil, err
	}

	r := InterfaceList{}
	for _, ni := range interfaces {
		i, err := ParseInterface(&ni, FlagNew)
		if err != nil {
			return nil, err
		}
		r = append(r, i)
	}

	return r, nil
}
