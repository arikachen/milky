package main

import (
	"flag"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/client-go/tools/clientcmd"

	configapi "github.com/openshift/origin/pkg/cmd/server/api"

	"github.com/arikachen/milky/pkg/node"
	"github.com/arikachen/milky/pkg/signals"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master     = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	plugin     = flag.String("plugin", "redhat/openshift-ovs-multitenant", "openshift network plugin.")
	hostName   = flag.String("hostname-override", "", "If non-empty, will use this string as identification instead of the actual hostname.")
	address    = flag.String("address", "", "node ip.")
	port       = flag.Uint("port", 8998, "node port.")
	mtu        = flag.Uint("mtu", 1450, "default mtu.")
	proxyMode  = flag.String("proxy-mode", "iptables", "proxy mode.")
	syncPeriod = flag.Uint("sync-period", 5, "iptables sync period.")
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}

	agent := node.SDNAgent{
		Options: configapi.NodeConfig{
			NodeName: *hostName,
			NodeIP:   *address,
			NetworkConfig: configapi.NodeNetworkConfig{
				NetworkPluginName: *plugin,
				MTU:               uint32(*mtu),
			},
		},
		Proxy: node.ProxyConfig{
			Mode:               *proxyMode,
			IPTablesSyncPeriod: time.Duration(*syncPeriod) * time.Second,
		},
	}

	err = agent.Run(cfg, stopCh)
	if err != nil {
		glog.Fatalf("Failed to start sdn controller: %v", err)
	}

	mux := http.NewServeMux()
	healthz.InstallHandler(mux)

	mux.Handle("/metrics", prometheus.Handler())

	server := &http.Server{
		Addr:    net.JoinHostPort(*address, strconv.Itoa(int(*port))),
		Handler: mux,
	}
	go wait.Until(func() {
		glog.Fatal(server.ListenAndServe())
	}, 0, wait.NeverStop)

	glog.Info("Started sdn agent")
	<-stopCh
	glog.Info("Shutting down sdn agent")
}
