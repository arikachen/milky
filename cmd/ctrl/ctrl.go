package main

import (
	"flag"
	"os"

	"github.com/golang/glog"

	kclientv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	kclientset "k8s.io/client-go/kubernetes"
	kv1core "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/leaderelectionconfig"

	configapi "github.com/openshift/origin/pkg/cmd/server/api"

	"github.com/arikachen/milky/pkg/controller"
)

var (
	kubeconfig        = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master            = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	plugin            = flag.String("plugin", "redhat/openshift-ovs-multitenant", "openshift network plugin.")
	clusterNetwork    = flag.String("cluster-network", "10.128.0.0/14", "default cluster network CIDR.")
	hostSubenetLength = flag.Uint("host-subnet-length", 9, "host subnet length.")
	serviceNetwork    = flag.String("service-network", "192.168.0.0/24", "k8s service network CIDR.")
)

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}

	kubeClient, err := kclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error new clientset: %v", err)
	}
	leaderElectionClient := kubernetes.NewForConfigOrDie(restclient.AddUserAgent(cfg, "leader-election"))

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

	run := func(stop <-chan struct{}) {
		err = sdn.Run(cfg, stop)
		if err != nil {
			glog.Fatalf("Failed to start sdn controller: %v", err)
		}
		glog.Info("Started sdn controller")
		select {}
	}

	id, err := os.Hostname()
	if err != nil {
		glog.Fatalf("Failed to get election id: %v", err)
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&kv1core.EventSinkImpl{Interface: kv1core.New(kubeClient.Core().RESTClient()).Events("")})
	recorder := eventBroadcaster.NewRecorder(api.Scheme, kclientv1.EventSource{Component: "openshift-sdn-controller"})

	leaderElectionCfg := leaderelectionconfig.DefaultLeaderElectionConfiguration()
	// using configmap lock
	leaderElectionCfg.ResourceLock = resourcelock.ConfigMapsResourceLock

	rl, err := resourcelock.New(leaderElectionCfg.ResourceLock,
		"kube-system",
		"openshift-sdn",
		leaderElectionClient.CoreV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: recorder,
		})
	if err != nil {
		glog.Fatalf("error creating lock: %v", err)
	}

	leaderelection.RunOrDie(leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: leaderElectionCfg.LeaseDuration.Duration,
		RenewDeadline: leaderElectionCfg.RenewDeadline.Duration,
		RetryPeriod:   leaderElectionCfg.RetryPeriod.Duration,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				glog.Fatalf("leaderelection lost")
			},
		},
	})
	panic("unreachable")
}
