package module

var linkCodes8636 = map[byte]string{
	0x00: "Reserved or unknown",
	0x01: "100G LinkCodes: 100G AOC or 25GAUI C2M AOC with worst BER of 5x10^(-5)",
	0x02: "100G LinkCodes: 100G Base-SR4 or 25GBase-SR",
	0x03: "100G LinkCodes: 100G Base-LR4",
	0x04: "100G LinkCodes: 100G Base-ER4",
	0x05: "100G LinkCodes: 100G Base-SR10",
	0x06: "100G LinkCodes: 100G CWDM4 MSA with FEC",
	0x07: "100G LinkCodes: 100G PSM4 Parallel SMF",
	0x08: "100G LinkCodes: 100G ACC or 25GAUI C2M ACC with worst BER of 5x10^(-5)",
	0x09: "100G LinkCodes: 100G CWDM4 MSA without FEC",
	0x0A: "(reserved or unknown)",
	0x0B: "100G LinkCodes: 100G Base-CR4 or 25G Base-CR CA-L",
	0x0C: "25G LinkCodes: 25G Base-CR CA-S",
	0x0D: "25G LinkCodes: 25G Base-CR CA-N",
	0x10: "40G LinkCodes: 40G Base-ER4",
	0x11: "4x10G LinkCodes: 10G Base-SR",
	0x12: "40G LinkCodes: 40G PSM4 Parallel SMF",
	0x13: "LinkCodes: G959.1 profile P1I1-2D1 (10709 MBd, 2km, 1310nm SM)",
	0x14: "LinkCodes: G959.1 profile P1S1-2D2 (10709 MBd, 40km, 1550nm SM)",
	0x15: "LinkCodes: G959.1 profile P1L1-2D2 (10709 MBd, 80km, 1550nm SM)",
	0x16: "10G LinkCodes: 10G Base-T with SFI electrical interface",
	0x17: "100G LinkCodes: 100G CLR4",
	0x18: "100G LinkCodes: 100G AOC or 25GAUI C2M AOC with worst BER of 10^(-12)",
	0x19: "100G LinkCodes: 100G ACC or 25GAUI C2M ACC with worst BER of 10^(-12)",
}

func getLinkCodes8636(v byte) string {
	name, ok := linkCodes8636[v]
	if !ok {
		return "Reserved or unknown"
	}
	return name
}
