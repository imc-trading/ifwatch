package module

import (
	"sort"
	"unsafe"
)

var transceivers8079 = map[uint64]string{
	(1 << 7):        "10G Ethernet: 10G Base-ER [SFF-8472 rev10.4 only]",
	(1 << 6):        "10G Ethernet: 10G Base-LRM",
	(1 << 5):        "10G Ethernet: 10G Base-LR",
	(1 << 4):        "10G Ethernet: 10G Base-SR",
	(1 << 3):        "Infiniband: 1X SX",
	(1 << 2):        "Infiniband: 1X LX",
	(1 << 1):        "Infiniband: 1X Copper Active",
	(1 << 0):        "Infiniband: 1X Copper Passive",
	(1 << (7 + 8)):  "ESCON: ESCON MMF, 1310nm LED",
	(1 << (6 + 8)):  "ESCON: ESCON SMF, 1310nm Laser",
	(1 << (5 + 8)):  "SONET: OC-192, short reach",
	(1 << (4 + 8)):  "SONET: SONET reach specifier bit 1",
	(1 << (3 + 8)):  "SONET: SONET reach specifier bit 2",
	(1 << (2 + 8)):  "SONET: OC-48, long reach",
	(1 << (1 + 8)):  "SONET: OC-48, intermediate reach",
	(1 << (0 + 8)):  "SONET: OC-48, short reach",
	(1 << (6 + 16)): "SONET: OC-12, single mode, long reach",
	(1 << (5 + 16)): "SONET: OC-12, single mode, inter. reach",
	(1 << (4 + 16)): "SONET: OC-12, short reach",
	(1 << (2 + 16)): "SONET: OC-3, single mode, long reach",
	(1 << (1 + 16)): "SONET: OC-3, single mode, inter. reach",
	(1 << (0 + 16)): "SONET: OC-3, short reach",
	(1 << (7 + 24)): "Ethernet: BASE-PX",
	(1 << (6 + 24)): "Ethernet: BASE-BX10",
	(1 << (5 + 24)): "Ethernet: 100BASE-FX",
	(1 << (4 + 24)): "Ethernet: 100BASE-LX/LX10",
	(1 << (3 + 24)): "Ethernet: 1000BASE-T",
	(1 << (2 + 24)): "Ethernet: 1000BASE-CX",
	(1 << (1 + 24)): "Ethernet: 1000BASE-LX",
	(1 << (0 + 24)): "Ethernet: 1000BASE-SX",
	(1 << (7 + 32)): "FC: very long distance (V)",
	(1 << (6 + 32)): "FC: short distance (S)",
	(1 << (5 + 32)): "FC: intermediate distance (I)",
	(1 << (4 + 32)): "FC: long distance (L)",
	(1 << (3 + 32)): "FC: medium distance (M)",
	(1 << (2 + 32)): "FC: Shortwave laser, linear Rx (SA)",
	(1 << (1 + 32)): "FC: Longwave laser (LC)",
	(1 << (0 + 32)): "FC: Electrical inter-enclosure (EL)",
	(1 << (7 + 40)): "FC: Electrical intra-enclosure (EL)",
	(1 << (6 + 40)): "FC: Shortwave laser w/o OFC (SN)",
	(1 << (5 + 40)): "FC: Shortwave laser with OFC (SL)",
	(1 << (4 + 40)): "FC: Longwave laser (LL)",
	(1 << (3 + 40)): "Active Cable",
	(1 << (2 + 40)): "Passive Cable",
	(1 << (1 + 40)): "FC: Copper FC-BaseT",
	(1 << (7 + 48)): "FC: Twin Axial Pair (TW)",
	(1 << (6 + 48)): "FC: Twisted Pair (TP)",
	(1 << (5 + 48)): "FC: Miniature Coax (MI)",
	(1 << (4 + 48)): "FC: Video Coax (TV)",
	(1 << (3 + 48)): "FC: Multimode, 62.5um (M6)",
	(1 << (2 + 48)): "FC: Multimode, 50um (M5)",
	(1 << (0 + 48)): "FC: Single Mode (SM)",
	(1 << (7 + 56)): "FC: 1200 MBytes/sec",
	(1 << (6 + 56)): "FC: 800 MBytes/sec",
	(1 << (4 + 56)): "FC: 400 MBytes/sec",
	(1 << (2 + 56)): "FC: 200 MBytes/sec",
	(1 << (0 + 56)): "FC: 100 MBytes/sec",
}

var transceivers8636 = map[uint64]string{
	(1 << 6):        "10G Ethernet: 10G Base-LRM",
	(1 << 5):        "10G Ethernet: 10G Base-LR",
	(1 << 4):        "10G Ethernet: 10G Base-SR",
	(1 << 3):        "40G Ethernet: 40G Base-CR4",
	(1 << 2):        "40G Ethernet: 40G Base-SR4",
	(1 << 1):        "40G Ethernet: 40G Base-LR4",
	(1 << 0):        "40G Ethernet: 40G Active Cable (XLPPI)",
	(1 << (3 + 8)):  "40G OTN (OTU3B/OTU3C)",
	(1 << (2 + 8)):  "SONET: OC-48, long reach",
	(1 << (1 + 8)):  "SONET: OC-48, intermediate reach",
	(1 << (0 + 8)):  "SONET: OC-48, short reach",
	(1 << (5 + 16)): "SAS 6.0G",
	(1 << (4 + 16)): "SAS 3.0G",
	(1 << (3 + 24)): "Ethernet: 1000BASE-T",
	(1 << (2 + 24)): "Ethernet: 1000BASE-CX",
	(1 << (1 + 24)): "Ethernet: 1000BASE-LX",
	(1 << (0 + 24)): "Ethernet: 1000BASE-SX",
	(1 << (7 + 32)): "FC: very long distance (V)",
	(1 << (6 + 32)): "FC: short distance (S)",
	(1 << (5 + 32)): "FC: intermediate distance (I)",
	(1 << (4 + 32)): "FC: long distance (L)",
	(1 << (3 + 32)): "FC: medium distance (M)",
	(1 << (1 + 32)): "FC: Longwave laser (LC)",
	(1 << (0 + 32)): "FC: Electrical inter-enclosure (EL)",
	(1 << (7 + 40)): "FC: Electrical intra-enclosure (EL)",
	(1 << (6 + 40)): "FC: Shortwave laser w/o OFC (SN)",
	(1 << (5 + 40)): "FC: Shortwave laser with OFC (SL)",
	(1 << (4 + 40)): "FC: Longwave laser (LL)",
	(1 << (7 + 48)): "FC: Twin Axial Pair (TW)",
	(1 << (6 + 48)): "FC: Twisted Pair (TP)",
	(1 << (5 + 48)): "FC: Miniature Coax (MI)",
	(1 << (4 + 48)): "FC: Video Coax (TV)",
	(1 << (3 + 48)): "FC: Multimode, 62.5m (M6)",
	(1 << (2 + 48)): "FC: Multimode, 50m (M5)",
	(1 << (1 + 48)): "FC: Multimode, 50um (OM3)",
	(1 << (0 + 48)): "FC: Single Mode (SM)",
	(1 << (7 + 56)): "FC: 1200 Mb/s",
	(1 << (6 + 56)): "FC: 800 Mb/s",
	(1 << (5 + 56)): "FC: 1600 Mb/s",
	(1 << (4 + 56)): "FC: 400 Mb/s",
	(1 << (2 + 56)): "FC: 200 Mb/s",
	(1 << (0 + 56)): "FC: 100 Mb/s",
}

type uint64arr []uint64

func (a uint64arr) Len() int           { return len(a) }
func (a uint64arr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a uint64arr) Less(i, j int) bool { return a[i] < a[j] }

type Transceiver [8]byte

func (t Transceiver) Uint64() uint64 {
	return *(*uint64)(unsafe.Pointer(&t[0]))
}

func getTransceiver8079(t Transceiver) []string {
	r := []string{}

	keys := uint64arr{}
	for k := range transceivers8079 {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	for _, k := range keys {
		if k&t.Uint64() != 0 {
			r = append(r, transceivers8079[k])
		}
	}

	return r
}

func getTransceiver8636(t Transceiver) []string {
	r := []string{}

	keys := uint64arr{}
	for k := range transceivers8636 {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	for _, k := range keys {
		if k&t.Uint64() != 0 {
			r = append(r, transceivers8636[k])
		}
	}

	return r
}
