package module

import (
	"fmt"

	"github.com/mickep76/color"
)

var (
	colorField = color.Cyan.String()
	colorValue = color.Green.String()
	colorList  = color.Yellow.String()
)

func (m *Module) String() string {
	switch m.Type {
	case "SFF-8079":
		return m.string8079()
	case "SFF-8636":
		return m.string8636()
	}
	return ""
}

func (m *Module) string8079() string {
	f1 := fmt.Sprintf("\t\t%s%%-52s%s : %s%%v%s", colorField, color.Reset, colorValue, color.Reset)
	f2 := fmt.Sprintf("\t\t%s%%-52s%s : %s%%v%s", colorField, color.Reset, colorList, color.Reset)

	s := fmt.Sprintf("\n"+f1, "Type", m.Type) +
		fmt.Sprintf("\n"+f1, "Identifier [0]", fmt.Sprintf("0x%02x (%s)", byte(m.Eeprom[0]), m.Identifier)) +
		fmt.Sprintf("\n"+f1, "Extended Identifier [1]", fmt.Sprintf("0x%02x (%s)", byte(m.Eeprom[1]), m.ExtIdentifier)) +
		fmt.Sprintf("\n"+f1, "Connector [2]", fmt.Sprintf("0x%02x (%s)", byte(m.Eeprom[2]), m.Connector)) +
		fmt.Sprintf("\n"+f1, "Transceiver Codes [3-10]", fmt.Sprintf("0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x", m.Eeprom[3], m.Eeprom[4], m.Eeprom[5], m.Eeprom[6], m.Eeprom[7], m.Eeprom[8], m.Eeprom[9], m.Eeprom[10])) +
		fmt.Sprintf("\n"+f2, "Transceiver Type", m.Transceiver[0])

	for _, t := range m.Transceiver[1:] {
		s += fmt.Sprintf("\n"+f2, "", t)
	}

	s += fmt.Sprintf("\n"+f1, "Encoding [11]", fmt.Sprintf("0x%02x (%s)", byte(m.Eeprom[11]), m.Encoding)) +
		fmt.Sprintf("\n"+f1, "BR, Nominal [12]", m.BrNominal) +
		fmt.Sprintf("\n"+f1, "Rate Identifier [13]", fmt.Sprintf("0x%02x", m.RateIdentifier)) +
		fmt.Sprintf("\n"+f1, "Length (SMF) [14]", m.LengthSmfKm) +
		fmt.Sprintf("\n"+f1, "Length (SMF) [15]", m.LengthSmfM) +
		fmt.Sprintf("\n"+f1, "Length (50um) [16]", m.Length50umM) +
		fmt.Sprintf("\n"+f1, "Length (62.5um) [17]", m.Length625umM) +
		fmt.Sprintf("\n"+f1, "Length (Copper) [18]", m.LengthCopper) +
		fmt.Sprintf("\n"+f1, "Length (OM3) [19]", m.LengthOm3) +
		fmt.Sprintf("\n"+f1, "Vendor [20-35]", m.Vendor) +
		fmt.Sprintf("\n"+f1, "Vendor OUI [37-39]", m.VendorOui) +
		fmt.Sprintf("\n"+f1, "Vendor PN [40-55]", m.VendorPn) +
		fmt.Sprintf("\n"+f1, "Vendor Rev [56-59]", m.VendorRev) +
		fmt.Sprintf("\n"+f1, "Option Values [64-65]", fmt.Sprintf("0x%02x 0x%02x", m.Options[0], m.Options[1])) +
		fmt.Sprintf("\n"+f1, "BR Margin, Max [66]", m.BrMax) +
		fmt.Sprintf("\n"+f1, "BR Margin, Min [67]", m.BrMin) +
		fmt.Sprintf("\n"+f1, "Vendor SN [68-83]", m.VendorSn) +
		fmt.Sprintf("\n"+f1, "Date Code [84-91]", m.DateCode)

	/*
		if s.Vendor.String() == "Arista Networks" && strings.HasPrefix(s.VendorPn.String(), "CAB-Q-S-") {
			str += fmt.Sprintf("%-50s : %x\n", "Vendor SA [120]", s.VendorAristaSa)
		}
	*/

	return s
}

func (m *Module) string8636() string {
	f1 := fmt.Sprintf("\t\t%s%%-52s%s : %s%%v%s", colorField, color.Reset, colorValue, color.Reset)
	f2 := fmt.Sprintf("\t\t%s%%-52s%s : %s%%v%s", colorField, color.Reset, colorList, color.Reset)

	s := fmt.Sprintf("\n"+f1, "Type", m.Type) +
		fmt.Sprintf("\n"+f1, "Identifier [128]", fmt.Sprintf("0x%02x (%s)", byte(m.Eeprom[0]), m.Identifier)) +
		fmt.Sprintf("\n"+f1, "Extended Identifier [129]", fmt.Sprintf("0x%02x", byte(m.Eeprom[129]))) +
		fmt.Sprintf("\n"+f2, "Extended Identifier Description", m.ExtIdentifier[0])

	for _, i := range m.ExtIdentifier[1:] {
		s += fmt.Sprintf("\n"+f2, "", i)
	}

	s += fmt.Sprintf("\n"+f1, "Connector [130]", fmt.Sprintf("0x%02x (%s)", byte(m.Eeprom[130]), m.Connector)) +
		fmt.Sprintf("\n"+f1, "Transceiver Codes [131-138]", fmt.Sprintf("0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x", m.Eeprom[131], m.Eeprom[132], m.Eeprom[133], m.Eeprom[134], m.Eeprom[135], m.Eeprom[136], m.Eeprom[137], m.Eeprom[138])) +
		fmt.Sprintf("\n"+f2, "Transceiver Type", m.Transceiver[0])

	for _, t := range m.Transceiver[1:] {
		s += fmt.Sprintf("\n"+f2, "", t)
	}

	s += fmt.Sprintf("\n"+f1, "Encoding [139]", fmt.Sprintf("0x%02x (%s)", byte(m.Eeprom[139]), m.Encoding)) +
		fmt.Sprintf("\n"+f1, "BR, Nominal [140]", m.BrNominal) +
		fmt.Sprintf("\n"+f1, "Rate Identifier [141]", fmt.Sprintf("0x%02x", m.RateIdentifier)) +
		fmt.Sprintf("\n"+f1, "Length (SMF) [142]", m.LengthSmf) +
		fmt.Sprintf("\n"+f1, "Length (OM3 50um) [143]", m.LengthOm3) +
		fmt.Sprintf("\n"+f1, "Length (OM2 50um) [144]", m.LengthOm2) +
		fmt.Sprintf("\n"+f1, "Length (OM1 62.5um) [145]", m.LengthOm1) +
		fmt.Sprintf("\n"+f1, "Length (Copper or Active cable) [146]", m.LengthCopper) +
		fmt.Sprintf("\n"+f1, "Vendor [148-163]", m.Vendor) +
		fmt.Sprintf("\n"+f1, "Vendor OUI [165-167]", m.VendorOui) +
		fmt.Sprintf("\n"+f1, "Vendor PN [168-183]", m.VendorPn) +
		fmt.Sprintf("\n"+f1, "Vendor Rev [184-185]", m.VendorRev) +
		fmt.Sprintf("\n"+f1, "Vendor SN [196-211]", m.VendorSn) +
		fmt.Sprintf("\n"+f1, "Date Code [212-219]", m.DateCode)
	return s
}
