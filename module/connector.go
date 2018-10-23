package module

var connectors = map[byte]string{
	0x00: "Unknown or unspecified",
	0x01: "SC",
	0x02: "Fibre Channel style 1 copper",
	0x03: "Fibre Channel style 2 copper",
	0x04: "BNC/TNC",
	0x05: "Fibre Channel coaxial headers",
	0x06: "FibreJack",
	0x07: "LC",
	0x08: "MT-RJ",
	0x09: "MU",
	0x0A: "SG",
	0x0B: "Optical pigtail",
	0x0C: "MPO Parallel Optic",
	0x0D: "MPO Parallel Optic - 2x16",
	0x20: "HSSDC II",
	0x21: "Copper pigtail",
	0x22: "RJ45",
	0x23: "No separable connector",
	0x24: "MXC 2x16",
}

func getConnector(v byte) string {
	name, ok := connectors[v]
	if !ok {
		return "Reserved or unknown"
	}
	return name
}
