%define name %NAME%
%define version %VERSION%
%define release %RELEASE%
%define buildroot %{_topdir}/BUILDROOT
%define sources %{_topdir}/SOURCES

BuildRoot: %{buildroot}
Source: %SOURCE%
Summary: %{name}
Name: %{name}
Version: %{version}
Release: %{release}
License: Apache License, Version 2.0
Group: System
AutoReqProv: no

%description
Listen to netlink events for network interfaces and publish it to etcd

%prep
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/etc/systemd/system
mkdir -p %{buildroot}/etc/init.d
mkdir -p %{buildroot}/etc/sysconfig
cp %{sources}/%{name} %{buildroot}/usr/bin
cp %{sources}/rpm/%{name}.service %{buildroot}/etc/systemd/system/%{name}.service
cp %{sources}/rpm/%{name}.initd %{buildroot}/etc/init.d/%{name}
cp %{sources}/rpm/%{name}.sysconfig %{buildroot}/etc/sysconfig/%{name}

%post
which systemctl &>/dev/null && systemctl daemon-reload

%preun
# Disable and stop on uninstall
if [ "${1}" == "0" ]; then
  if which systemctl &>/dev/null; then
    systemctl stop %{name}
    systemctl disable %{name}
  else
    service %{name} stop
    chkconfig %{name} off
  fi
fi

%postun
# Restart on upgrade
if [ "${1}" == "1" ]; then
  if which systemctl &>/dev/null; then
    systemctl condrestart %{name}
  else
    service %{name} condrestart
  fi
fi

%files
%defattr(-,root,root)
/usr/bin/%{name}
/etc/systemd/system/%{name}.service
%attr(755,-,-) /etc/init.d/%{name}
%config(noreplace) /etc/sysconfig/%{name}
