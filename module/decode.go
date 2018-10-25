package module

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnknownType = errors.New("unknown type")

func decode8079(eeprom []byte) *Module {
	m := &Module{
		Type:           TypeSff8079,
		Identifier:     getIdentifier(eeprom[0]),
		ExtIdentifier:  getExtIdentifier8079(eeprom[1]),
		Connector:      getConnector(eeprom[2]),
		Transceiver:    getTransceiver8079(Transceiver{eeprom[3], eeprom[4], eeprom[5], eeprom[6], eeprom[7], eeprom[8], eeprom[9], eeprom[10]}),
		Encoding:       getEncoding8079(eeprom[11]),
		BrNominal:      int(eeprom[12]) * 100,
		RateIdentifier: eeprom[13],
		LengthSmfKm:    int(eeprom[14]),
		LengthSmfM:     int(eeprom[15]),
		Length50umM:    int(eeprom[16]) * 10,
		Length625umM:   int(eeprom[17]) * 10,
		LengthCopper:   int(eeprom[18]),
		LengthOm3:      int(eeprom[19]) * 10,
		Vendor:         strings.TrimSpace(string(eeprom[20:36])),
		VendorOui:      fmt.Sprintf("%x:%x:%x", eeprom[37], eeprom[38], eeprom[39]),
		VendorPn:       strings.TrimSpace(string(eeprom[40:56])),
		VendorRev:      string(eeprom[56:60]),
		Options:        eeprom[64:66],
		BrMax:          int(eeprom[66]),
		BrMin:          int(eeprom[67]),
		VendorSn:       strings.TrimSpace(string(eeprom[68:84])),
		DateCode:       fmt.Sprintf("20%s-%s-%s", string(eeprom[84:86]), string(eeprom[86:88]), string(eeprom[88:90])),
		Eeprom:         eeprom,
	}

	if m.Vendor == "Arista Networks" && strings.HasPrefix(m.VendorPn, "CAB-Q-S-") {
		m.VendorSa = int(eeprom[120])
	}

	return m
}

func decode8636(eeprom []byte) *Module {
	return &Module{
		Type:           TypeSff8636,
		Identifier:     getIdentifier(eeprom[128]),
		ExtIdentifier:  getExtIdentifier8636(eeprom[129]),
		Connector:      getConnector(eeprom[130]),
		Transceiver:    getTransceiver8636(Transceiver{eeprom[131], eeprom[132], eeprom[133], eeprom[134], eeprom[135], eeprom[136], eeprom[137], eeprom[138]}),
		Encoding:       getEncoding8636(eeprom[139]),
		BrNominal:      int(eeprom[140]) * 100,
		RateIdentifier: eeprom[141],
		LengthSmf:      int(eeprom[142]),
		LengthOm3:      int(eeprom[143]),
		LengthOm2:      int(eeprom[144]),
		LengthOm1:      int(eeprom[145]),
		LengthCopper:   int(eeprom[146]),
		Vendor:         strings.TrimSpace(string(eeprom[148:164])),
		VendorOui:      fmt.Sprintf("%x:%x:%x", eeprom[165], eeprom[166], eeprom[167]),
		VendorPn:       strings.TrimSpace(string(eeprom[168:184])),
		VendorRev:      strings.TrimSpace(string(eeprom[184:186])),
		LinkCodes:      getLinkCodes8636(eeprom[192]),
		Options:        eeprom[193:196],
		VendorSn:       strings.TrimSpace(string(eeprom[196:212])),
		DateCode:       fmt.Sprintf("20%s-%s-%s", string(eeprom[212:214]), string(eeprom[214:216]), string(eeprom[216:218])),
		Eeprom:         eeprom,
	}
}

func Decode(eeprom []byte) (*Module, error) {
	if len(eeprom) < 256 {
		return nil, fmt.Errorf("eeprom needs to be 256 bytes or larger got: %d bytes", len(eeprom))
	}

	if (eeprom[0] == 2 || eeprom[0] == 3) && eeprom[1] == 4 {
		return decode8079(eeprom), nil
	}

	if eeprom[128] == 12 || eeprom[128] == 13 || eeprom[128] == 17 {
		return decode8636(eeprom), nil
	}

	return nil, ErrUnknownType
}
