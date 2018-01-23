package main

import (
	"log"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/imc-trading/ifwatch/cmd"
)

const version = "0.2.9"

func main() {
	usage := `ifwatch

Usage:
  ifwatch print [<interface>] [--json]
  ifwatch module <interface> [--hex]
  ifwatch info [<interface>]
  ifwatch daemon [--skip-interfaces=<interfaces>] [--skip-drivers=<drivers>] [--backend=<backend>] [--topic=<topic>] [--brokers=<brokers>] [--prefix=<prefix>] [--endpoints=<endpoints>] [--cert=<file>] [--key=<file>] [--ca=<file>] [--refresh=<seconds>] [--ttl=<seconds>] [--user=<user>] [--pass=<file>] [--timeout=<seconds>]
  ifwatch refresh [--skip-interfaces=<interfaces>] [--skip-drivers=<drivers>] [--backend=<backend>] [--topic=<topic>] [--brokers=<brokers>] [--prefix=<prefix>] [--endpoints=<endpoints>] [--cert=<file>] [--key=<file>] [--ca=<file>] [--user=<user>] [--pass=<file>] [--timeout=<seconds>]
  ifwatch watch [<host/interface>] [--json] [--host=<host>] [--interface=<interface>] [--backend=<backend>] [--topic=<topic>] [--brokers=<brokers>] [--prefix=<prefix>] [--endpoints=<endpoints>] [--cert=<file>] [--key=<file>] [--ca=<file>] [--user=<user>] [--pass=<file>] [--timeout=<seconds>]
  ifwatch log [<host/interface>] [--json] [--backend=<backend>] [--topic=<topic>] [--brokers=<brokers>] [--prefix=<prefix>] [--endpoints=<endpoints>] [--cert=<file>] [--key=<file>] [--ca=<file>] [--user=<user>] [--pass=<file>] [--timeout=<seconds>]
  ifwatch ls [--json] [--prefix=<prefix>] [--endpoints=<endpoints>] [--cert=<file>] [--key=<file>] [--ca=<file>] [--user=<user>] [--pass=<file>] [--timeout=<seconds>]
  ifwatch flush [<host/interface>] [--backend=<backend>] [--topic=<topic>] [--brokers=<brokers>] [--prefix=<prefix>] [--endpoints=<endpoints>] [--cert=<file>] [--key=<file>] [--ca=<file>] [--user=<user>] [--pass=<file>] [--timeout=<seconds>]
  ifwatch -h | --help
  ifwatch --version

Arguments:
  print                                 Verbose info about interfaces.
  module                                Verbose info about module.
  info                                  Brief info about interfaces.
  daemon                                Run as a daemon and publish interface events.
  refresh                               Generate refresh event.
  watch                                 Watch all interface events.
  log                                   Log of interface events.
  ls                                    List all interfaces. (etcd only)
  flush                                 ...tbd... Flush events for a host or interface (etcd only).

Options:
  -h --help                             Show this screen.
  --version                             Show version.
  --json                                Print as JSON.
  --hex                                 Dump module as hex.
  --skip-interfaces=<interfaces>        Comma-delimited list of interfaces to skip (env: IFWATCH_SKIP_INTERFACES). [default: lo]
  --skip-drivers=<drivers>              Comma-delimited list of driver to skip (env: IFWATCH_SKIP_DRIVERS). [default: veth,bridge]
  --backend=<backend>                   Backend for events etcd or kafka (env. IFWATCH_BACKEND). [default: etcd]
  --topic=<topic>                       Kafka topic (env: IFWATCH_TOPIC). [default: ifwatch]
  --brokers=<brokers>                   Kafka brokers (env: IFWATCH_BROKERS).
  --prefix=<prefix>                     etcd prefix (env: IFWATCH_PREFIX). [default: /ifwatch]
  --endpoints=<endpoints>               etcd comma-delimited list of hosts in the cluster (env: IFWATCH_ENDPOINTS).
  --refresh=<seconds>                   Send refresh event every number of seconds.
  --ttl=<seconds>                       ...tbd... TTL of event in seconds (etcd only).
  --user=<user>                         ...tbd... etcd user (env: IFWATCH_USER).
  --pass=<file>                         ...tbd... etcd password file (env: IFWATCH_PASS).
  --cert=<file>                         ...tbd... TLS certificate file (env: IFWATCH_CERT).
  --key=<file>                          ...tbd... TLS key file (env: IFWATCH_KEY).
  --ca=<file>                           ...tbd... TLS CA bundle file (env: IFWATCH_CA).
  --timeout=<seconds>                   Connection timeout in seconds. [default: 5]
`

	args, err := docopt.Parse(usage, nil, true, "ifwatch "+version, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	if args["--endpoints"] == nil {
		args["--endpoints"] = os.Getenv("IFWATCH_ENDPOINTS")
	}

	if args["--brokers"] == nil {
		args["--brokers"] = os.Getenv("IFWATCH_BROKERS")
	}

	if args["print"].(bool) {
		cmd.Print(args)
		return
	}

	if args["module"].(bool) {
		cmd.Module(args)
		return
	}

	if args["daemon"].(bool) {
		cmd.Daemon(args)
		return
	}

	if args["refresh"].(bool) {
		cmd.Refresh(args)
		return
	}

	if args["watch"].(bool) {
		cmd.Watch(args)
		return
	}

	if args["log"].(bool) {
		cmd.Log(args)
		return
	}

	if args["ls"].(bool) {
		cmd.Ls(args)
		return
	}

	if args["info"].(bool) {
		cmd.Info(args)
		return
	}
}
