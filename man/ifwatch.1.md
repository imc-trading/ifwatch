% IFWATCH(1)
% Michael Persson
% Oct. 2018

# NAME

ifwatch

# SYNOPSIS

**ifwatch** [**-h**] [**--version**] [**-conf** *CONFIGFILE*] [**-print-conf**]

**ifwatch** [**-conf** *CONFIGFILE*] print [**-h**] [**-json**]

**ifwatch** [**-conf** *CONFIGFILE*] publish [**-h**]

**ifwatch** [**-conf** *CONFIGFILE*] subscribe [**-h**] [**-json**] [**-host** *HOST*] [*INTERFACE*...]

**ifwatch** [**-conf** *CONFIGFILE*] history [**-h**] [**-json**] [**-host** *HOST*] [*INTERFACE*...]

# DESCRIPTION

**ifwatch** Listen to netlink events for network interfaces and publish it to Kafka.

# OPTIONS

**-h**, **-help**
: Display a help message.

**-version**
: Print version.

**-conf** *CONFIGFILE*
: Configuration file, defaults to /etc/directord.toml.

**-print-conf**
: Print configuration.

**-json**
: Print as JSON.

**-host**
: Filter by hostname.

# COMMAND

**print**
: Print network interfaces.

**publish**
: Publish events to Kafka.

**subscribe**
: Subscribe to and print events from Kafka.

**history**
: Show events history from Kafka.

# POSITIONAL ARGUMENTS

*INTERFACE*...
: Filter by interface(s).

# FILES

**/etc/ifwatch.toml**
:	Configuration file.

# SEE ALSO

**ifwatch.toml(5)**
