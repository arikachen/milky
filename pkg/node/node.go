package node

import (
	"fmt"
	"os"
	"strings"
	"time"

	kclientv1 "k8s.io/api/core/v1"
	kclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	kv1core "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/apis/componentconfig"
	kclientsetinternal "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	kinternalinformers "k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion"

	configapi "github.com/openshift/origin/pkg/cmd/server/api"
	netclientset "github.com/openshift/origin/pkg/network/generated/internalclientset"
	sdnnode "github.com/openshift/origin/pkg/network/node"

	"github.com/golang/glog"
)

type ProxyConfig struct {
	Mode               string
	IPTablesSyncPeriod time.Duration
}

type SDNAgent struct {
	Options configapi.NodeConfig
	Proxy   ProxyConfig
}

func (a *SDNAgent) getProxyMode() componentconfig.ProxyMode {
	if a.Proxy.Mode == fmt.Sprintf("%s", componentconfig.ProxyModeUserspace) {
		return componentconfig.ProxyModeUserspace
	}
	return componentconfig.ProxyModeIPTables
}

func (a *SDNAgent) getNodeName() string {
	if a.Options.NodeName != "" {
		return a.Options.NodeName
	}
	nodename, err := os.Hostname()
	if err != nil {
		glog.Fatalf("Get hostname failed: %v", err)
	}
	hostname := strings.ToLower(strings.TrimSpace(nodename))
	glog.Infof("Resolved hostname to %q", hostname)
	return hostname
}

func (a *SDNAgent) Run(clientConfig *restclient.Config, stopCh <-chan struct{}) error {
	networkClient, err := netclientset.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	kubeClientset, err := kclientset.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	internalkClient, err := kclientsetinternal.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	hostname := a.getNodeName()

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(&kv1core.EventSinkImpl{Interface: kubeClientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, kclientv1.EventSource{Component: "openshift-sdn", Host: hostname})

	informerFactory := kinternalinformers.NewSharedInformerFactory(internalkClient, time.Second*30)

	node, err := sdnnode.New(&sdnnode.OsdnNodeConfig{
		PluginName:         a.Options.NetworkConfig.NetworkPluginName,
		Hostname:           hostname,
		SelfIP:             a.Options.NodeIP,
		RuntimeEndpoint:    "",
		MTU:                a.Options.NetworkConfig.MTU,
		NetworkClient:      networkClient,
		KClient:            internalkClient,
		KubeInformers:      informerFactory,
		IPTablesSyncPeriod: a.Proxy.IPTablesSyncPeriod,
		ProxyMode:          a.getProxyMode(),
		EnableHostports:    true,
		Recorder:           recorder,
	})
	if err != nil {
		return err
	}

	err = node.Start()
	if err != nil {
		return err
	}

	informerFactory.Start(stopCh)
	informerFactory.WaitForCacheSync(stopCh)

	return nil
}
