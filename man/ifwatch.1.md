% IFWATCH(1)
% Michael Persson
% Oct. 2018

# NAME

ifwatch

# SYNOPSIS

**ifwatch** [**-h**] [**-version**] [**-conf** *CONFIGFILE*] [**-print-conf**]

# DESCRIPTION

**ifwatch** Listen to netlink events for network interfaces and publish it to Kafka.

# OPTIONS

**-h**, **-help**
:	Display a help message.

**-version**
:       Print version.

**-conf**
:	Configuration file, defaults to /etc/directord.toml.

**-print-conf**
:	Print configuration.

# FILES

**/etc/ifwatch.toml**
:	Configuration file.

# SEE ALSO

**ifwatch.toml(5)**
