# milky
openshift-sdn standalone

## Getting Started

It is working with [kubeadm][network].

```
# Add openshift custom resources  
kubectl apply -f https://raw.githubusercontent.com/arikachen/milky/master/scripts/deploy/openshift-sdn-crd.yaml

#Add RABC
kubectl apply -f https://raw.githubusercontent.com/arikachen/milky/master/scripts/deploy/rabc.yaml

#Create pods
kubectl apply -f https://raw.githubusercontent.com/arikachen/milky/master/scripts/deploy/milky.yaml
```

[network]: https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#pod-network

## TODO

* HA
* metrics