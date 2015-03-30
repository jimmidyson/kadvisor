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
	kube_client "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/clientauth"
	log "github.com/Sirupsen/logrus"
)

func initializeKubeClient(host string, apiVersion string, insecure bool, clientAuthFile string) *kube_client.Client {
	kubeConfig := kube_client.Config{
		Host:     host,
		Version:  apiVersion,
		Insecure: insecure,
	}

	if len(clientAuthFile) > 0 {
		clientAuth, err := clientauth.LoadFromFile(clientAuthFile)
		if err != nil {
			log.Fatal(err)
		}
		kubeConfig, err = clientAuth.MergeWithConfig(kubeConfig)
		if err != nil {
			log.Fatal(err)
		}
	}

	kubernetesClient := kube_client.NewOrDie(&kubeConfig)
	if _, err := kubernetesClient.ServerVersion(); err != nil {
		log.WithField("error", err).Fatal("Could not validate Kubernetes master")
	}
	return kubernetesClient
}
