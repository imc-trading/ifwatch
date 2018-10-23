package module

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnknownType = errors.New("unknown type")

func decode8079(eeprom []byte) *Module {
	m := &Module{
		Type: TypeSff8079,
	}
	m.Identifier = getIdentifier(eeprom[0])
	m.ExtIdentifier = getExtIdentifier8079(eeprom[1])
	m.Connector = getConnector(eeprom[2])
	m.Transceiver = getTransceiver8079(Transceiver{eeprom[3], eeprom[4], eeprom[5], eeprom[6], eeprom[7], eeprom[8], eeprom[9], eeprom[10]})
	m.Encoding = getEncoding8079(eeprom[11])
	m.BrNominal = int(eeprom[12]) * 100
	m.RateIdentifier = eeprom[13]
	m.LengthSmfKm = int(eeprom[14])
	m.LengthSmfM = int(eeprom[15])
	m.Length50umM = int(eeprom[16]) * 10
	m.Length625umM = int(eeprom[17]) * 10
	m.LengthCopper = int(eeprom[18])
	m.LengthOm3 = int(eeprom[19]) * 10
	m.Vendor = strings.TrimSpace(string(eeprom[20:36]))
	m.VendorOui = fmt.Sprintf("%x:%x:%x", eeprom[37], eeprom[38], eeprom[39])
	m.VendorPn = strings.TrimSpace(string(eeprom[40:56]))
	m.VendorRev = string(eeprom[56:60])
	m.Options = eeprom[64:66]
	m.BrMax = int(eeprom[66])
	m.BrMin = int(eeprom[67])
	m.VendorSn = strings.TrimSpace(string(eeprom[68:84]))
	m.DateCode = fmt.Sprintf("20%s-%s-%s", string(eeprom[84:86]), string(eeprom[86:88]), string(eeprom[88:90]))
	m.Eeprom = eeprom
	return m
}

func decode8636(eeprom []byte) *Module {
	m := &Module{
		Type: TypeSff8636,
	}
	m.Identifier = getIdentifier(eeprom[128])
	m.ExtIdentifier = getExtIdentifier8636(eeprom[129])
	m.Connector = getConnector(eeprom[130])
	m.Transceiver = getTransceiver8636(Transceiver{eeprom[131], eeprom[132], eeprom[133], eeprom[134], eeprom[135], eeprom[136], eeprom[137], eeprom[138]})
	m.Encoding = getEncoding8636(eeprom[139])
	m.BrNominal = int(eeprom[140]) * 100
	m.RateIdentifier = eeprom[141]
	m.LengthSmf = int(eeprom[142])
	m.LengthOm3 = int(eeprom[143])
	m.LengthOm2 = int(eeprom[144])
	m.LengthOm1 = int(eeprom[145])
	m.LengthCopper = int(eeprom[146])
	m.Vendor = strings.TrimSpace(string(eeprom[148:164]))
	m.VendorOui = fmt.Sprintf("%x:%x:%x", eeprom[165], eeprom[166], eeprom[167])
	m.VendorPn = strings.TrimSpace(string(eeprom[168:184]))
	m.VendorRev = strings.TrimSpace(string(eeprom[184:186]))
	m.LinkCodes = getLinkCodes8636(eeprom[192])
	m.Options = eeprom[193:196]
	m.VendorSn = strings.TrimSpace(string(eeprom[196:212]))
	m.DateCode = fmt.Sprintf("20%s-%s-%s", string(eeprom[212:214]), string(eeprom[214:216]), string(eeprom[216:218]))
	m.Eeprom = eeprom
	return m
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
