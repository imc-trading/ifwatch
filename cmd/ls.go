package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/imc-trading/ifwatch/backend/etcd"
)

func Ls(args map[string]interface{}) {
	if args["--backend"].(string) == "etcd" {
		if args["--endpoints"] == nil {
			log.Fatalf("missing option: --endpoints")
		}
		endpoints := strings.Split(args["--endpoints"].(string), ",")

		tmout, err := strconv.Atoi(args["--timeout"].(string))
		if err != nil {
			log.Fatalf("strconv: %v", err)
		}

		prefix = args["--prefix"].(string)

		if sub, err = etcd.NewSubscriber(endpoints, tmout, prefix); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("not supported for kafka backend")
	}

	keys, err := sub.Keys()
	if err != nil {
		log.Fatal(err)
	}

	if args["--json"].(bool) {
		res, _ := json.Marshal(keys)
		fmt.Println(string(res))
	} else {
		fmt.Println(strings.Join(keys, "\n"))
	}
}
