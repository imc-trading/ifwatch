%undefine _missing_build_ids_terminate_build

%define sources %{_topdir}/SOURCES

BuildRoot: %{_topdir}/BUILDROOT
Source: ifwatch-%{_ver}-%{_rel}.tar.gz
Summary: ifwatch - Listen to netlink events for network interfaces and publish it to Kafka
Name: ifwatch
Version: %{_ver}
Release: %{_rel}
License: Apache License, Version 2.0
Group: System
AutoReqProv: no
BuildRequires: go git pandoc

%description
Listen to netlink events for network interfaces and publish it to Kafka

%prep
# Setup Go environment
GOPATH=$PWD/go
PATH=$GOPATH/bin:$PATH
export GOPATH PATH

# Install Glide
go get github.com/Masterminds/glide/...

# Setup src
mkdir -p $GOPATH/src/%{_src}
tar zvxf %{sources}/ifwatch-%{_ver}-%{_rel}.tar.gz -C $GOPATH/src/%{_src}

%build
# Setup Go environment
GOPATH=$PWD/go
PATH=$GOPATH/bin:$PATH
export GOPATH PATH

cd $GOPATH/src/%{_src}

# Download packages
glide up

# Build
go build -ldflags "-X main.version=%{_ver}"

# Convert markdown to manpage
pandoc man/%{name}.1.md -s -t man -o %{name}.1
pandoc man/%{name}.toml.5.md -s -t man -o %{name}.toml.5
gzip %{name}.1
gzip %{name}.toml.5

%install
cd go/src/%{_src}

mkdir -p %{buildroot}/usr/sbin
mkdir -p %{buildroot}/usr/share/man/{man1,man5}
mkdir -p %{buildroot}/usr/share/%{name}
mkdir -p %{buildroot}/etc/systemd/system
mkdir -p %{buildroot}/etc/sysconfig

cp %{name} %{buildroot}/usr/sbin/
cp %{name}.toml.5.gz %{buildroot}/usr/share/man/man5/
cp %{name}.1.gz %{buildroot}/usr/share/man/man1/
cp rpm/%{name}.service %{buildroot}/etc/systemd/system/%{name}.service
cp rpm/%{name}.sysconfig %{buildroot}/etc/sysconfig/%{name}
cp rpm/%{name}.toml %{buildroot}/etc/
cp LICENSE %{buildroot}/usr/share/%{name}/

%post
which systemctl &>/dev/null && systemctl daemon-reload

%preun
# Disable and stop on uninstall
if [ "${1}" == "0" ]; then
  if which systemctl &>/dev/null; then
    systemctl stop %{name}
    systemctl disable %{name}
  fi
fi

%postun
# Restart on upgrade
if [ "${1}" == "1" ]; then
  if which systemctl &>/dev/null; then
    systemctl condrestart %{name}
  fi
fi

%files
%defattr(-,root,root)
/usr/sbin/%{name}
/etc/systemd/system/%{name}.service
%config(noreplace) /etc/sysconfig/%{name}
%config(noreplace) /etc/%{name}.toml
/usr/share/man/man1/%{name}.1.gz
/usr/share/man/man5/%{name}.toml.5.gz
/usr/share/%{name}/LICENSE
