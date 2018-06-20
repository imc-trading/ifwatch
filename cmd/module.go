package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/imc-trading/ifwatch/netx"
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

	if args["--hex"].(bool) {
		eeprom, err := e.ModuleEeprom(args["<interface>"].(string))
		if err != nil && err != ErrOpNotSupp {
			log.Fatalf("ethtool module info for %s: %v", args["<interface>"].(string), err)
		}

		fmt.Println(hex.EncodeToString(eeprom))
		return
	}

	all, err := netx.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range all {
		if i.Name == args["<interface>"].(string) {
			host, _ = os.Hostname()
			fmt.Printf("%s%s/%s/module%s", white, host, i.Name, clear)
			for _, line := range strings.Split(i.Module.StringCol(), "\n") {
				fmt.Printf("\n\t\t%s", line)
			}
			fmt.Printf("\n")
			return
		}
	}
}
