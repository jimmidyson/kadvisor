/*
 * Copyright 2015 Red Hat, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
