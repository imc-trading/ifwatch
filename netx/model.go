package netx

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	lightRed = "\x1b[91m"
	red      = "\x1b[31m"
	green    = "\x1b[32m"
	yellow   = "\x1b[33m"
	blue     = "\x1b[34m"
	magenta  = "\x1b[35m"
	cyan     = "\x1b[36m"
	white    = "\x1b[37m"
	clear    = "\x1b[0m"
)

// InterfaceList list of interfaces
type InterfaceList []*Interface

// JSON return interface list as JSON
func (i *InterfaceList) JSON() []byte {
	b, _ := json.Marshal(i)
	return b
}

// JSONPretty return interface list as indented JSON
func (i *InterfaceList) JSONPretty() []byte {
	b, _ := json.MarshalIndent(i, "", "  ")
	return b
}

// JSON return interface as JSON
func (i *Interface) JSON() []byte {
	b, _ := json.Marshal(i)
	return b
}

// JSONPretty interface as indented JSON
func (i *Interface) JSONPretty() []byte {
	b, _ := json.MarshalIndent(i, "", "  ")
	return b
}

func omitStr(f, k string, v string) string {
	if v == "" {
		return ""
	}
	return fmt.Sprintf(f, k, v)
}

func strCol(k string, v string, omit bool, c1 string, c2 string) string {
	if omit && v == "" {
		return ""
	}
	return fmt.Sprintf("\n\t%s%-58s%s : %s%s%s", c1, k, clear, c2, v, clear)
}

func omitInt(f string, k string, v int) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf(f, k, v)
}

func intCol(k string, v int, omit bool, c1 string, c2 string) string {
	if omit && v == 0 {
		return ""
	}
	return fmt.Sprintf("\n\t%s%-58s%s : %s%d%s", c1, k, clear, c2, v, clear)
}

func omitUint(f string, k string, v uint) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf(f, k, v)
}

func joinSprintf(f string, l []string) string {
	r := ""
	for _, s := range l {
		r += fmt.Sprintf(f, s)
	}
	return r
}

func joinStrCol(k string, l []string, c1 string, c2 string) string {
	if len(l) < 1 {
		return ""
	}

	r := strCol(k, l[0], false, c1, c2)
	for _, s := range l[1:] {
		r += strCol("", s, false, c1, c2)
	}
	return r
}

func (i *Interface) String() string {
	mStr := "\n"
	if i.Module != nil {
		mStr = "\n\tModule:"
		for _, line := range strings.Split(i.Module.String(), "\n") {
			mStr += fmt.Sprintf("\n\t\t%s", line)
		}
	}

	return fmt.Sprintf("\n\t%-58s : %s", "Name", i.Name) +
		fmt.Sprintf("\n\t%-58s : %s", "DNS Names", i.Flags[0]) +
		joinSprintf(fmt.Sprintf("\n\t%-58s : %%s", " "), i.Flags[1:]) +
		fmt.Sprintf("\n\t%-58s : %d", "Index", i.Index) +
		omitStr("\n\t%-58s : %s", "Slot", i.Slot) +
		omitStr("\n\t%-58s : %s", "Driver", i.Driver) +
		omitInt("\n\t%-58s : %d", "MTU", i.MTU) +
		omitStr("\n\t%-58s : %s", "Hardware Address", i.HwAddr) +
		omitStr("\n\t%-58s : %s", "Permanent Hardware Address", i.PermHwAddr) +
		fmt.Sprintf("\n\t%-58s : %s", "Flags", i.Flags[0]) +
		joinSprintf(fmt.Sprintf("\n\t%-58s : %%s", " "), i.Flags[1:]) +
		omitStr("\n\t%-58s : %s", "IPv4 Address", i.IPv4) +
		omitStr("\n\t%-58s : %s", "Netmask", i.Netmask) +
		omitStr("\n\t%-58s : %s", "Network", i.Network) +
		omitStr("\n\t%-58s : %s", "IPv6 Address", i.IPv6) +
		mStr
}

// StringCol interfaces as a string with ascii color
func (i *Interface) StringCol() string {
	mStr := "\n"
	if i.Module != nil {
		mStr = fmt.Sprintf("\n\t%s%s%s", white, "Module", clear)
		for _, line := range strings.Split(i.Module.StringCol(), "\n") {
			mStr += fmt.Sprintf("\n\t\t%s", line)
		}
	}

	hwAddrCol := green
	permHwAddrCol := green
	if i.PermHwAddr != "" && i.HwAddr != i.PermHwAddr {
		hwAddrCol = lightRed
		permHwAddrCol = red
	}

	return strCol("Name", i.Name, false, cyan, green) +
		joinStrCol("DNS Names", i.DNSNames, cyan, yellow) +
		intCol("Index", i.Index, false, cyan, green) +
		strCol("Slot", i.Slot, true, cyan, green) +
		strCol("Driver", i.Driver, true, cyan, green) +
		intCol("MTU", i.MTU, true, cyan, green) +
		strCol("Hardware Address", i.HwAddr, true, cyan, hwAddrCol) +
		strCol("Permanent Hardware Address", i.PermHwAddr, true, cyan, permHwAddrCol) +
		joinStrCol("Flags", i.Flags, cyan, yellow) +
		strCol("IPv4 Address", i.IPv4, true, cyan, green) +
		strCol("Netmask", i.Netmask, true, cyan, green) +
		strCol("Network", i.Network, true, cyan, green) +
		strCol("IPv6 Address", i.IPv6, true, cyan, green) +
		mStr
}
