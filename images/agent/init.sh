#!/bin/bash

CNI_BIN_PATH=/opt/cni/bin
CNI_CONF=/etc/cni/net.d/80-openshift-network.conf

# install openshift cni
if [ ! -f $CNI_BIN_PATH/openshift-sdn ];then
    cp /opt/openshift-sdn $CNI_BIN_PATH
fi

# need to make node ready, so controller may scheduler first.
# TODO
if [ ! -f $CNI_CONF ];then
cat > $CNI_CONF <<EOF
{
  "cniVersion": "0.2.0",
  "name": "openshift-sdn",
  "type": "openshift-sdn"
}
EOF

    sleep 3
fi
