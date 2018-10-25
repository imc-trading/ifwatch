package module

import (
	"encoding/hex"
	"encoding/json"
)

type ModuleType string

const (
	TypeSff8079 = ModuleType("SFF-8079")
	TypeSff8636 = ModuleType("SFF-8636")
)

type Bytes []byte

type Module struct {
	Type           ModuleType `json:"type"`
	Identifier     string     `json:"identifier"`
	ExtIdentifier  []string   `json:"extIdentifier"`
	Connector      string     `json:"connector"`
	Transceiver    []string   `json:"transceiver"`
	Encoding       string     `json:"encoding"`
	BrNominal      int        `json:"brNominal"`
	RateIdentifier byte       `json:"rateIdentifier"`
	LengthSmfKm    int        `json:"lengthSmfKm"`
	LengthSmfM     int        `json:"lengthSmfM"`
	LengthSmf      int        `json:"lengthSmf"`
	Length50umM    int        `json:"length50umM"`
	Length625umM   int        `json:"length625umM"`
	LengthCopper   int        `json:"lengthCopper"`
	LengthOm3      int        `json:"lengthOm3"`
	LengthOm2      int        `json:"lengthOm2"`
	LengthOm1      int        `json:"lengthOm1"`
	Vendor         string     `json:"vendor"`
	VendorOui      string     `json:"vendorOui"`
	VendorPn       string     `json:"vendorPn"`
	VendorRev      string     `json:"vendorRev"`
	LinkCodes      string     `json:"linkCodes"`
	Options        Bytes      `json:"options"`
	BrMax          int        `json:"brMax"`
	BrMin          int        `json:"brMin"`
	VendorSn       string     `json:"vendorSn"`
	DateCode       string     `json:"dateCode"`
	VendorSa       int        `json:"vendorSa"`
	Eeprom         Bytes      `json:"eeprom"`
}

func (b *Bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(*b))
}

func (b *Bytes) UnmarshalJSON(in []byte) error {
	var s string
	if err := json.Unmarshal(in, &s); err != nil {
		return err
	}

	v, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	*b = v

	return nil
}
