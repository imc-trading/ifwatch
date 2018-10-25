package module

var encodings8079 = map[byte]string{
	0x00: "Unspecified",
	0x01: "8B/10B",
	0x02: "4B/5B",
	0x03: "NRZ",
	0x04: "Manchester",
	0x05: "SONET Scrambled",
	0x06: "64B/66B",
	0x07: "256B/257B (transcoded FEC-enabled data)",
	0x08: "PAM4",
}

var encodings8636 = map[byte]string{
	0x00: "Unspecified",
	0x01: "8B/10B",
	0x02: "4B/5B",
	0x03: "NRZ",
	0x04: "SONET Scrambled",
	0x05: "64B/66B",
	0x06: "Manchester",
	0x07: "256B/257B (transcoded FEC-enabled data)",
	0x08: "PAM4",
}

func getEncoding8079(v byte) string {
	name, ok := encodings8079[v]
	if !ok {
		return "Reserved or unknown"
	}
	return name
}

func getEncoding8636(v byte) string {
	name, ok := encodings8636[v]
	if !ok {
		return "Reserved or unknown"
	}
	return name
}
