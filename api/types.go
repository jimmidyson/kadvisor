package api

import (
	"net"

	kube_api "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
)

type Node struct {
	*kube_api.Node
	addresses map[kube_api.NodeAddressType]string
	ipAddress string
}

func (n *Node) IsMetricsCollectable() bool {
	for _, condition := range n.Status.Conditions {
		if condition.Status == kube_api.ConditionTrue {
			return true
		}
	}
	return false
}

func (n *Node) GetIpAddress() string {
	if len(n.ipAddress) == 0 {
		if n.addresses == nil {
			n.addresses = make(map[kube_api.NodeAddressType]string)
			for _, address := range n.Status.Addresses {
				n.addresses[address.Type] = address.Address
			}
		}
		if len(n.addresses[kube_api.NodeLegacyHostIP]) > 0 {
			n.ipAddress = n.addresses[kube_api.NodeLegacyHostIP]
			return n.ipAddress
		}
		if len(n.addresses[kube_api.NodeInternalIP]) > 0 {
			n.ipAddress = n.addresses[kube_api.NodeInternalIP]
			return n.ipAddress
		}
		if len(n.addresses[kube_api.NodeHostName]) > 0 {
			addrs, err := net.LookupIP(n.addresses[kube_api.NodeHostName])
			if err != nil {
				n.ipAddress = addrs[0].String()
			}
			return n.ipAddress
		}
	}
	return n.ipAddress
}
