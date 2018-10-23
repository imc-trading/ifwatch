% IFWATCH.TOML(5)
% Michael Persson
% Oct. 2018

# NAME

ifwatch.toml - configuration file for ifwatch

# MAIN SECTION

**LogLevel**
: Only print same or higher log level debug, info, warning, error or fatal.

```
LogLevel = "error"
```

**LogNoColor**
: Log with no color.

```
LogNoColor = true
```

**logNoDate**
: Log with no date.

```
logNoDate = true
```

**skipInterfaces**
: Skip interfaces.

```
skipInterfaces = ["lo"]
```

**skipDrivers**
: Skip drivers.

```
skipDrivers = ["veth", "bridge"]
```

**topic**
: Kafka topic.

```
topic = "ifwatch"
```

**brokers**
: Kafka brokers.

```
brokers = ["kafka1", "kafka2", "kafka3"]
```

**timeout**
: Kafka connection timeout in seconds.

```
timeout = 3
```

# EXAMPLE

```
LogLevel = "error"
LogNoColor = true
logNoDate = true
brokers = ["kafka1", "kafka2", "kafka3"]
```

# SEE ALSO

**ifwatch(1)**