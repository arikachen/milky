# Dockerfile - CentOS 7

FROM centos:7.4.1708

MAINTAINER Arika Chen <eaglesora@gmail.com>

COPY init.sh /opt

ADD http://mirror.centos.org/centos/7/cloud/x86_64/openstack-pike/openvswitch-2.7.3-1.1fc27.el7.x86_64.rpm /opt
ADD https://github.com/arikachen/milky/releases/download/v0.1-alpha/milky-0.1.tar.gz /opt

RUN yum install -y \
        ethtool \
        git \
        net-tools \
        iproute \
        wget \
        numactl-libs \
        libpcap \
        openssl \
        conntrack-tools \
        libnetfilter_conntrack \
    && rpm -ivh /opt/openvswitch*.rpm && rm -rf /opt/openvswitch*.rpm \
    && chmod 755 /opt/init.sh \
    && tar -xf /opt/milky-0.1.tar.gz -C /tmp \
    && cp /tmp/milky/milky-agent /opt \
    && cp /tmp/milky/openshift-sdn /opt \
    && rm -rf /opt/milky-0.1.tar.gz && rm -rf /tmp/milky \
    && yum clean all

WORKDIR /opt

ENTRYPOINT ["/opt/milky-agent"]
