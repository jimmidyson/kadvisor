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

package kubernetes

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/fabric8io/kadvisor/sources"
)

const (
	defaultInsecure       = false
	defaultApiVersion     = "v1beta3"
	defaultClientAuthFile = ""
)

func init() {
	sources.Register("kubernetes", New)
}

type KubernetesMetricsSource struct {
	master         string
	apiVersion     string
	insecure       bool
	clientAuthFile string
}

func New(uri string) (sources.Source, error) {
	parsedUrl, err := url.Parse(os.ExpandEnv(uri))
	if err != nil {
		log.WithFields(log.Fields{"url": uri, "error": err}).Fatal("Could not create Kubernetes source")
	}

	if len(parsedUrl.Scheme) == 0 {
		log.WithField("url", uri).Fatal("Missing scheme in Kubernetes source")
	}
	if len(parsedUrl.Host) == 0 {
		log.WithField("url", uri).Fatal("Missing host in Kubernetes source")
	}

	source := &KubernetesMetricsSource{
		master:         fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host),
		apiVersion:     defaultApiVersion,
		insecure:       defaultInsecure,
		clientAuthFile: defaultClientAuthFile,
	}
	options := parsedUrl.Query()
	if len(options["apiVersion"]) >= 1 {
		source.apiVersion = options["apiVersion"][0]
	}
	if len(options["insecure"]) >= 1 {
		insecure, err := strconv.ParseBool(options["insecure"][0])
		if err != nil {
			log.WithField("url", uri).Fatal("Invalid insecure option for Kubernetes source")
		}
		source.insecure = insecure
	}
	if len(options["auth"]) >= 1 {
		source.clientAuthFile = options["auth"][0]
	}
	return source, nil
}

func (k *KubernetesMetricsSource) Start() {
	kubernetesClient := initializeKubeClient(k.master, k.apiVersion, k.insecure, k.clientAuthFile)
	nodeList, err := kubernetesClient.Nodes().List()
	if err != nil {
		log.Fatal(err)
	}
	for _, kubeNode := range nodeList.Items {
		node := &node{
			Node: &kubeNode,
		}
		log.Debugf("New node: %v (%v)", node.Name, node.UID)
		if !node.isMetricsCollectable() {
			log.WithFields(log.Fields{"node": node.Name, "conditions": node.Status.Conditions, "phase": node.Status.Phase}).Info("Skipping node: node not in collectable state")
			break
		}
		if len(node.getIpAddress()) == 0 {
			log.Debugf("Skipping node %v (%v): cannot find IP address. Node addresses are: %s", node.Name, node.UID, node.Status.Addresses)
			break
		}
		log.Infof("Collecting from node %v (%v)", node.Name, node.UID)
	}
}