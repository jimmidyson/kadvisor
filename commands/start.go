package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fabric8io/kadvisor/api"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start kAdvisor",
	Long:  "kAdvisor can start all-in-one as a single node",
}

func init() {
	startCmd.Run = start
}

func start(cmd *cobra.Command, args []string) {
	InitializeConfig()
	kubernetesClient := InitializeKubeClient()
	nodeList, err := kubernetesClient.Nodes().List()
	if err != nil {
		log.Fatal(err)
	}
	for _, kubeNode := range nodeList.Items {
		node := &api.Node{
			Node: &kubeNode,
		}
		log.Debugf("New node: %v (%v)", node.Name, node.UID)
		if !node.IsMetricsCollectable() {
			log.Debugf("Skipping node %v (%v): node not in collectable state. Node statuses are: %s", node.Name, node.UID, node.Status.Conditions)
			break
		}
		if len(node.GetIpAddress()) == 0 {
			log.Debugf("Skipping node %v (%v): cannot find IP address. Node addresses are: %s", node.Name, node.UID, node.Status.Addresses)
			break
		}
		log.Infof("Collecting from node %v (%v)", node.Name, node.UID)
	}
}
