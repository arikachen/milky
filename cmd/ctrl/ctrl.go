package main

import (
	"flag"

	"github.com/golang/glog"
	"k8s.io/client-go/tools/clientcmd"

	configapi "github.com/openshift/origin/pkg/cmd/server/api"

	"github.com/arikachen/milky/pkg/controller"
	"github.com/arikachen/milky/pkg/signals"
)

var (
	kuberconfig       = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master            = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	plugin            = flag.String("plugin", "redhat/openshift-ovs-multitenant", "openshift network plugin.")
	clusterNetwork    = flag.String("cluster-network", "10.128.0.0/14", "default cluster network CIDR.")
	hostSubenetLength = flag.Uint("host-subnet-length", 9, "host subnet length.")
	serviceNetwork    = flag.String("service-network", "192.168.0.0/24", "k8s service network CIDR.")
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kuberconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}

	sdn := controller.SDNController{
		NetworkConfig: configapi.MasterNetworkConfig{
			NetworkPluginName: *plugin,
			ClusterNetworks: []configapi.ClusterNetworkEntry{
				{
					CIDR:             *clusterNetwork,
					HostSubnetLength: uint32(*hostSubenetLength),
				},
			},
			ServiceNetworkCIDR: *serviceNetwork,
		},
	}
	err = sdn.Run(cfg, stopCh)
	if err != nil {
		glog.Fatalf("Failed to start sdn controller: %v", err)
	}
}
