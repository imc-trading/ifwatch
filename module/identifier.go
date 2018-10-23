package module

var identifiers = map[byte]string{
	0x00: "No module present, unknown, or unspecified",
	0x01: "GBIC",
	0x02: "Module soldered to motherboard",
	0x03: "SFP",
	0x04: "300 pin XBI",
	0x05: "XENPAK",
	0x06: "XFP",
	0x07: "XFF",
	0x08: "XFP-E",
	0x09: "XPAK",
	0x0A: "X2",
	0x0B: "DWDM-SFP",
	0x0C: "QSFP",
	0x0D: "QSFP+",
	0x0E: "CXP",
	0x0F: "Shielded Mini Multilane HD 4X",
	0x10: "Shielded Mini Multilane HD 8X",
	0x11: "QSFP28",
	0x12: "CXP2/CXP28",
	0x13: "CDFP Style 1/Style 2",
	0x14: "Shielded Mini Multilane HD 4X Fanout Cable",
	0x15: "Shielded Mini Multilane HD 8X Fanout Cable",
	0x16: "CDFP Style 3",
	0x17: "MicroQSFP",
}

func getIdentifier(v byte) string {
	name, ok := identifiers[v]
	if !ok {
		return "Reserved or unknown"
	}
	return name
}
