FROM centos:centos7

# Install RPM Build
RUN yum install -y rpm-build epel-release && yum clean all

ARG NAME
ARG VER
ARG REL
ARG SRC

# Setup build environment
RUN mkdir -p /root/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
COPY rpm/${NAME}.spec /root/rpmbuild/SPECS/
COPY ${NAME}-${VER}-${REL}.tar.gz /root/rpmbuild/SOURCES/

# Install build deps
RUN yum-builddep -y /root/rpmbuild/SPECS/${NAME}.spec && yum clean all

# Build RPM
RUN rpmbuild -ba \
      --target="x86_64" \
      --define "_topdir /root/rpmbuild" \
      --define "_ver ${VER}" \
      --define "_rel ${REL}" \
      --define "_src ${SRC}" \
      /root/rpmbuild/SPECS/${NAME}.spec
