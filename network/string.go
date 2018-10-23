package network

import (
	"fmt"

	"github.com/mickep76/color"
)

var (
	format   = "\n\t%s%-60s%s : %s%v%s"
	labelCol = color.White.String()
	fieldCol = color.Cyan.String()
	valueCol = color.Green.String()
	warnCol  = color.LightRed.String()
	critCol  = color.Red.String()
	listCol  = color.Yellow.String()
)

func inList(a string, l []string) bool {
	for _, b := range l {
		if a == b {
			return true
		}
	}
	return false
}

func strCol(k string, v string, omit bool, c1 string, c2 string) string {
	if omit && v == "" {
		return ""
	}
	return fmt.Sprintf(format, c1, k, color.Reset, c2, v, color.Reset)
}

func intCol(k string, v int, omit bool, c1 string, c2 string) string {
	if omit && v == 0 {
		return ""
	}
	return fmt.Sprintf(format, c1, k, color.Reset, c2, v, color.Reset)
}

func joinStrCol(k string, l []string, c1 string, c2 string) string {
	if len(l) < 1 {
		return ""
	}

	str := strCol(k, l[0], false, c1, c2)
	for _, s := range l[1:] {
		str += strCol("", s, false, c1, c2)
	}
	return str
}

func (i *Interface) String() string {
	hwAddrCol := valueCol
	permHwAddrCol := valueCol
	if i.PermHwAddr != "" && i.HwAddr != i.PermHwAddr {
		hwAddrCol = warnCol
		permHwAddrCol = critCol
	}

	s := strCol("Name", i.Name, false, fieldCol, valueCol) +
		joinStrCol("DNS Names", i.DNSNames, fieldCol, listCol) +
		intCol("Index", i.Index, false, fieldCol, valueCol) +
		strCol("Slot", i.Slot, true, fieldCol, valueCol) +
		strCol("Driver", i.Driver, true, fieldCol, valueCol) +
		intCol("MTU", i.MTU, true, fieldCol, valueCol) +
		strCol("Hardware Address", i.HwAddr, true, fieldCol, hwAddrCol) +
		strCol("Permanent Hardware Address", i.PermHwAddr, true, fieldCol, permHwAddrCol) +
		joinStrCol("Flags", i.Flags, fieldCol, color.Yellow.String()) +
		strCol("IPv4 Address", i.IPv4, true, fieldCol, valueCol) +
		strCol("Netmask", i.Netmask, true, fieldCol, valueCol) +
		strCol("Network", i.Network, true, fieldCol, valueCol) +
		strCol("IPv6 Address", i.IPv6, true, fieldCol, valueCol)

	if i.Module != nil {
		s += fmt.Sprintf("\n\t%sModule%s", color.White, color.Reset)
		s += i.Module.String()
	}

	return s
}
