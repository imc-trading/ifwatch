# ifwatch

Print network interface and module info. Listen to netlink events for network interfaces and publish it to etcd or kafka

# TOC

- [Start etcd](#start-etcd)  
- [Start Kafka](#start-kafka)  
- [Install ifwatch using RPM](#install-rpm)  
- [Info interfaces](#info-interfaces)
- [Print interfaces](#print-interfaces)
- [Watch changes](#watch-changes)
- [Log of changes](#log-of-changes)
- [Run daemon](#run-daemon)
- [Refresh event](#refresh-event)
- [Build RPM](#build-rpm)
- [Compile](#compile)
- [Usage](#usage)

<a name="start-etcd"/>

# Start etcd

First install Docker.

```bash
docker run --rm -p 2379:2379 -p 2380:2380 --name etcd gcr.io/etcd-development/etcd:v3.2 /usr/local/bin/etcd \
        --name my-etcd-1 \
        --data-dir /etcd-data \
        --listen-client-urls http://0.0.0.0:2379 \
        --advertise-client-urls http://0.0.0.0:2379 \
        --listen-peer-urls http://0.0.0.0:2380 \
        --initial-advertise-peer-urls http://0.0.0.0:2380 \
        --initial-cluster my-etcd-1=http://0.0.0.0:2380 \
        --initial-cluster-token my-etcd-token \
        --initial-cluster-state new
```

<a name="start-kafka"/>

# Start Kafka

...TBD...

<a name="install-rpm"/>

# Install ifwatch using RPM

```bash
yum install ifwatch
```

Verify install, needs to run as root if you want module info.

```
ifwatch print
```

In order to start modify /etc/sysconfig/ifwatch and add "--endpoints etcd1:2379" to use etcd or "--backend kafka --brokers kafka1:9092" for using Kafka.

```bash
service ifwatch start
```

Or:

```bash
systemctl ifwatch start
```

<a name="info-interfaces"/>

# Info interfaces

Needs to run as root if you want module info.

```bash
ifwatch info
```

<a name="print-interfaces"/>

# Print interfaces

Needs to run as root if you want module info.

```bash
ifwatch print
```

Or as JSON.

```bash
ifwatch print --json | jq -r .
```

<a name="watch-changes"/>

# Watch changes

```bash
ifwatch watch --backend kafka --brokers kafka1:9092,kafka2:9092,kafka3:9092
```

<a name="log-of-changes"/>

# Log of changes

```bash
ifwatch log --backend kafka --brokers kafka1:9092,kafka2:9092,kafka3:9092
```

<a name="run-daemon"/>

# Run daemon

Needs to run as root if you want module info.

```bash
./ifwatch daemon --backend kafka --brokers kafka1:9092,kafka2:9092,kafka3:9092
```

<a name="refresh-event"/>

# Refresh event

Generate refresh event, useful for testing or if use Kafka as backend where events have a TTL.

```bash
./ifwatch refresh --backend kafka --brokers kafka1:9092,kafka2:9092,kafka3:9092
```

<a name="build-rpm"/>

# Build RPM

First install Docker.

```
make build-rpm
```

<a name="compile"/>

# Compile

First install Go.

```bash
export GOPATH="~/go"
mkdir -p $GOPATH/src/github.com/imc-trading
cd $GOPATH/src/github.com/imc-trading
git clone git@github.com/imc-trading/ifwatch.git
cd ifwatch
make deps build
```

<a name="usage"/>

# Usage

```bash
ifwatch

Usage:
  ifwatch print [<interface>] [--json]
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
```
