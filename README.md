# milky
openshift-sdn standalone

## Getting Started

It is working with [kubeadm][network].

```
# Add openshift custom resources  
kubectl apply -f https://raw.githubusercontent.com/arikachen/milky/v0.1/scripts/deploy/openshift-sdn-crd.yaml

#Add RABC
kubectl apply -f https://raw.githubusercontent.com/arikachen/milky/v0.1/scripts/deploy/rabc.yaml

#Create pods
kubectl apply -f https://raw.githubusercontent.com/arikachen/milky/v0.1/scripts/deploy/milky.yaml
```

[network]: https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#pod-network

