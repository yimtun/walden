# walden
the mirror for kubernetes cluster ip


将 k8s cluster ip  暴露出来 处理南北流量

集群内部 节点 cluster ip 网络路径不受影响


划分 子网 block    外部物理网络 的服务也使用cluster ip 所在的大网段

相当于是cluster IP  的影子




#  pcs


```azure
yum install -y pcs pacemaker corosync fence-agents-all resource-agents
```

```azure
systemctl enable pcsd
systemctl start pcsd
    
    
```


# use lvs full-nat



# walden-agent
walden-agent  can reachable pod network
walden-agent bonind  cluster-ip

















