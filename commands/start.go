package commands

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
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
	log.Info("Running")
	nodeList, err := kubernetesClient.Nodes().List()
	if err != nil {
		log.Fatal(err)
	}
	marshalledNodes, err := json.MarshalIndent(nodeList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(string(marshalledNodes[:]))
}
