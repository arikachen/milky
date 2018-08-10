# milky

Milky is base on Openshift SDN, it is aim to provide pod network for k8s.

[OpenShift SDN][openshift] provides three SDN plug-ins for configuring the pod network:

* redhat/openshift-ovs-subnet
* redhat/openshift-ovs-multitenant
* redhat/openshift-ovs-networkpolicy

[openshift]: https://docs.openshift.com/container-platform/3.5/architecture/additional_concepts/sdn.html#architecture-additional-concepts-sdn

## Getting Started

### Prerequire

* k8s 1.7+
* openvswitch 2.7+
* cni

It is working with [kubeadm][network].

```
# Add openshift custom resources  
kubectl apply -f https://raw.githubusercontent.com/arikachen/milky/master/scripts/deploy/openshift-sdn-crd.yaml

#Add RABC
kubectl apply -f https://raw.githubusercontent.com/arikachen/milky/master/scripts/deploy/rabc.yaml

#Create pods
kubectl apply -f https://raw.githubusercontent.com/arikachen/milky/master/scripts/deploy/milky.yaml

kubectl annotate netnamespace kube-system pod.network.openshift.io/multitenant.change-network="global"
```

[network]: https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#pod-network

## Managing Pod Networks

When using **redhat/openshift-ovs-multitenant** plugin, you can manage the separate pod overlay networks.

### Joining Network

To join two tenant networks, they can connect with each other.

```
kubectl create namespace net1
kubectl create namespace net2
kubectl annotate netnamespace net1 pod.network.openshift.io/multitenant.change-network="join:net2"
```

### Isolating Network

```
kubectl annotate netnamespace net1 pod.network.openshift.io/multitenant.change-network="isolate"
```

### Making Network Global

The tenant can connect to all other tenant's network.

```
kubectl annotate netnamespace net1 pod.network.openshift.io/multitenant.change-network="global"
```

## Limiting Pod BandWidth

Specify the data traffic speed using **kubernetes.io/ingress-bandwidth** and 
**kubernetes.io/egress-bandwidth** annotations when creating pod.
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  annotations:
    kubernetes.io/egress-bandwidth: 10M
    kubernetes.io/ingress-bandwidth: 10M
spec:
  containers:
    - name: nginx
      image: nginx
```
