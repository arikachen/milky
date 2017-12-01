# Dockerfile - CentOS 7

FROM centos:7.4.1708

MAINTAINER Arika Chen <eaglesora@gmail.com>

ADD https://github.com/arikachen/milky/releases/download/v0.1-alpha/milky-0.1.tar.gz /opt

RUN yum install -y \
        ethtool \
        git \
        net-tools \
        iproute \
    && tar -xf /opt/milky-0.1.tar.gz -C /tmp \
    && cp /tmp/milky/milky-ctrl /opt \
    && rm -rf /opt/milky-0.1.tar.gz && rm -rf /tmp/milky \
    && yum clean all

WORKDIR /opt

ENTRYPOINT ["/opt/milky-ctrl"]
