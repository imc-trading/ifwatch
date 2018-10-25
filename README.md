# ifwatch

Print network interface and module info. Listen to netlink events for network interfaces and publish it to kafka

# Usage

```bash
Usage: ifwatch [OPTIONS] COMMAND

Commands:
  print         Print network interfaces.
  publish       Publish events to Kafka.
  subscribe     Subscribe to and print events from Kafka.
  history       Show events history from Kafka.

Options:
  -conf string
    	Configuration file. (default "/etc/ifwatch.toml")
  -print-conf
    	Print configuration.
  -version
    	Print application version.
```
