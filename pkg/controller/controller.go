package controller

import (
	"fmt"
	"time"

	restclient "k8s.io/client-go/rest"
	kclientsetinternal "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	kinternalinformers "k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion"

	configapi "github.com/openshift/origin/pkg/cmd/server/api"
	netclientset "github.com/openshift/origin/pkg/network/generated/internalclientset"
	sdnmaster "github.com/openshift/origin/pkg/network/master"
)

type SDNController struct {
	NetworkConfig configapi.MasterNetworkConfig
}

func (c *SDNController) Run(clientConfig *restclient.Config, stopCh <-chan struct{}) error {
	networkClient, err := netclientset.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	kClient, err := kclientsetinternal.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	informerFactory := kinternalinformers.NewSharedInformerFactory(kClient, time.Second*30)

	err = sdnmaster.Start(
		c.NetworkConfig,
		networkClient,
		kClient,
		informerFactory,
	)
	if err != nil {
		return fmt.Errorf("failed to start SDN plugin controller: %v", err)
	}
	informerFactory.Start(stopCh)
	informerFactory.WaitForCacheSync(stopCh)

	return nil
}
