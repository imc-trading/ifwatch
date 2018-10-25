package module

import "fmt"

const (
	PwrClassMask = 0xC0
	PwrClass1    = (0 << 6)
	PwrClass2    = (1 << 6)
	PwrClass3    = (2 << 6)
	PwrClass4    = (3 << 6)

	ClieCodeMask = 0x10
	NoClieCode   = (0 << 4)
	ClieCode     = (1 << 4)

	CdrInTxMask = 0x08
	NoCdrInTx   = (0 << 3)
	CdrInTx     = (1 << 3)

	CdrInRxMask = 0x04
	NoCdrInRx   = (0 << 2)
	CdrInRx     = (1 << 2)

	ExtPwrClassMask   = 0x03
	ExtPwrClassUnused = 0
	ExtPwrClass5      = 1
	ExtPwrClass6      = 2
	ExtPwrClass7      = 3
)

var extIdentifiers8079 = map[byte]string{
	0x00: "GBIC not specified / not MOD_DEF compliant",
	0x04: "GBIC/SFP defined by 2-wire interface ID",
	0x07: "GBIC compliant with MOD_DEF",
}

var pwrClassNames8636 = map[byte]string{
	PwrClass1: "1.5 W max. power consumption",
	PwrClass2: "2.0 W max. power consumption",
	PwrClass3: "2.5 W max. power consumption",
	PwrClass4: "3.5 W max. power consumption",
}

var clieCodeNames8636 = map[byte]string{
	NoClieCode: "No CLEI code present",
	ClieCode:   "CLEI code present",
}

var cdrInTxNames8636 = map[byte]string{
	NoCdrInTx: "No CDR in TX",
	CdrInTx:   "CDR in TX",
}

var cdrInRxNames8636 = map[byte]string{
	NoCdrInRx: "No CDR in RX",
	CdrInRx:   "CDR in RX",
}

var extPwrClassNames8636 = map[byte]string{
	ExtPwrClassUnused: "unused (legacy setting)",
	ExtPwrClass5:      "4.0 W max. power consumption",
	ExtPwrClass6:      "4.5 W max. power consumption",
	ExtPwrClass7:      "5.0 W max. power consumption",
}

func getExtIdentifier8079(v byte) []string {
	name, ok := extIdentifiers8079[v]
	if !ok {
		return []string{"Unknown"}
	}
	return []string{name}
}

func getExtIdentifier8636(v byte) []string {
	s := []string{
		pwrClassNames8636[v&PwrClassMask],
		clieCodeNames8636[v&ClieCodeMask],
		fmt.Sprintf("%s, %s", cdrInTxNames8636[v&CdrInTxMask], cdrInRxNames8636[v&CdrInRxMask]),
	}

	if v&ExtPwrClassMask != ExtPwrClassUnused {
		s = append(s, extPwrClassNames8636[v&ExtPwrClassMask])
	}

	return s
}
