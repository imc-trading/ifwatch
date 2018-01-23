package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/mickep76/ethtool"
)

var ErrOpNotSupp = errors.New("operation not supported")

func Module(args map[string]interface{}) {
	if args["<interface>"] == nil {
		return
	}

	e, err := ethtool.NewEthtool()
	if err != nil && err != ErrOpNotSupp {
		log.Fatalf("new ethtool: %v", err)
	}

	eeprom, err := e.ModuleEeprom(args["<interface>"].(string))
	if err != nil && err != ErrOpNotSupp {
		log.Fatalf("ethtool module info for %s: %v", args["<interface>"].(string), err)
	}

	fmt.Println(hex.EncodeToString(eeprom))
}
