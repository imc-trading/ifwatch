package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/imc-trading/ifwatch/netx"
)

func Print(args map[string]interface{}) {
	all, err := netx.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	interfaces := netx.InterfaceList{}
	if args["<interface>"] != nil {
		for _, i := range all {
			if i.Name == args["<interface>"].(string) {
				interfaces = append(interfaces, i)
				break
			}
		}

		if len(interfaces) == 0 {
			log.Fatalf("interface doesn't exist: %s", args["<interface>"].(string))
		}
	} else {
		interfaces = all
	}

	if args["--json"].(bool) {
		if len(interfaces) == 1 {
			fmt.Println(string(interfaces[0].JSONPretty()))
		} else {
			fmt.Println(string(interfaces.JSONPretty()))
		}
	} else {
		host, _ = os.Hostname()
		for _, i := range interfaces {
			fmt.Printf("%s%s/%s%s", white, host, i.Name, clear)
			fmt.Println(i.StringCol())
		}
	}
}
