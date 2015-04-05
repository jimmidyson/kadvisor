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
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jimmidyson/kadvisor/api"
	"github.com/jimmidyson/kadvisor/extpoints"
	sources_api "github.com/jimmidyson/kadvisor/sources/api"
	"github.com/spf13/viper"
)

func init() {
	extpoints.SourceFactories.Register(New, "kubernetes")
}

const (
	defaultInsecure       = false
	defaultApiVersion     = "v1beta3"
	defaultClientAuthFile = ""
)

type KubernetesMetricsSource struct {
	master         string
	apiVersion     string
	insecure       bool
	clientAuthFile string
	updateInterval time.Duration
}

func New(uri string, options map[string][]string) (sources_api.Source, error) {
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
		updateInterval: viper.GetDuration("defaultPoll"),
	}
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
	if len(options["interval"]) >= 1 {
		interval, err := time.ParseDuration(options["interval"][0])
		if err != nil {
			log.WithField("interval", options["interval"][0]).Fatal("Invalid poll interval for Kubernetes source")
		}
		source.updateInterval = interval
	}
	return source, nil
}

func (k *KubernetesMetricsSource) Start(pipelineChan chan interface{}, wg *sync.WaitGroup) chan api.Stop {
	kubernetesClient := initializeKubeClient(k.master, k.apiVersion, k.insecure, k.clientAuthFile)
	nodeList, err := kubernetesClient.Nodes().List()
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(1)
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

	stopChan := make(chan api.Stop)
	go func() {
		ticker := time.NewTicker(k.updateInterval)
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				pipelineChan <- "Hello"
			case <-stopChan:
				return
			default:
			}
		}
	}()

	return stopChan
}
