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
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fabric8io/kadvisor/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	kube_client "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/clientauth"
)

var KadvisorCmd = &cobra.Command{
	Use:   "kadvisor",
	Short: "KAdvisor is a metrics collector & publisher for Kubernetes",
	Long:  "A configurable metrics collector & publisher for Kubernetes",
}

var kadvisorCmdV *cobra.Command

var (
	Verbose                  bool
	CfgFile                  string
	PollDuration             time.Duration
	KubernetesMaster         string
	KubernetesApiVersion     string
	KubernetesInsecure       bool
	KubernetesClientAuthFile string
	InfluxdbSinkUrl          string
	InfluxdbServiceName      string
	InfluxdbSecure           bool
)

var kubernetesClient *kube_client.Client

func Execute() {
	AddCommands()
	utils.StopOnErr(KadvisorCmd.Execute())
}

func AddCommands() {
	KadvisorCmd.AddCommand(startCmd)
}

//Initializes flags
func init() {
	KadvisorCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "config file")
	KadvisorCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose logging")
	KadvisorCmd.PersistentFlags().DurationVarP(&PollDuration, "poll", "p", 10*time.Second, "poll duration")

	KadvisorCmd.PersistentFlags().StringVarP(&KubernetesMaster, "kubernetes-master", "k", "", "Kubernetes master")
	KadvisorCmd.PersistentFlags().StringVar(&KubernetesApiVersion, "kubernetes-api-version", "v1beta3", "Kubernetes API version")
	KadvisorCmd.PersistentFlags().BoolVar(&KubernetesInsecure, "kubernetes-skip-tls-verify", false, "Skip TLS verify of Kubernetes master certificate")
	KadvisorCmd.PersistentFlags().StringVar(&KubernetesClientAuthFile, "kubernetes-client-auth-file", "", "Kubernetes clien auth file")

	KadvisorCmd.PersistentFlags().StringVarP(&InfluxdbSinkUrl, "influxdb", "i", "", "InfluxDB URL")
	KadvisorCmd.PersistentFlags().StringVar(&InfluxdbServiceName, "influxdb-service", "INFLUXDB", "InfluxDB service name")
	KadvisorCmd.PersistentFlags().BoolVar(&InfluxdbSecure, "influxdb-secure", false, "Use https for InfluxDB service")

	kadvisorCmdV = KadvisorCmd
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func InitializeConfig() {
	if len(CfgFile) > 0 {
		viper.SetConfigFile(CfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("$HOME/.kadvisor")
		viper.AddConfigPath("/etc/kadvisor")
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("kadvisor")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Debug("Unable to locate Config file")
	}

	viper.SetDefault("verbose", false)
	viper.SetDefault("poll", 10*time.Second)
	viper.SetDefault("influxdbSecure", false)
	viper.SetDefault("kubernetesApiVersion", "v1beta2")
	viper.SetDefault("kubernetesInsecure", false)

	if kadvisorCmdV.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("verbose", Verbose)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("poll").Changed {
		viper.Set("poll", PollDuration)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-master").Changed {
		viper.Set("kubernetesMaster", KubernetesMaster)
	} else if len(os.Getenv("KUBERNETES_MASTER")) > 0 {
		viper.Set("kubernetesMaster", "${KUBERNETES_MASTER}")
	} else if len(os.Getenv("KUBERNETES_SERVICE_HOST")) > 0 && len(os.Getenv("KUBERNETES_SERVICE_PORT")) > 0 {
		viper.Set("kubernetesMaster", "https://${KUBERNETES_SERVICE_HOST}:${KUBERNETES_SERVICE_PORT}")
	}
	viper.Set("kubernetesMaster", os.ExpandEnv(viper.GetString("kubernetesMaster")))
	if len(viper.GetString("kubernetesMaster")) == 0 {
		log.Fatal("Kubernetes master is not set")
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-api-version").Changed {
		viper.Set("kubernetesApiVersion", KubernetesApiVersion)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-skip-tls-verify").Changed {
		viper.Set("kubernetesInsecure", KubernetesInsecure)
	}

	if kadvisorCmdV.PersistentFlags().Lookup("influxdb").Changed {
		viper.Set("influxdb", InfluxdbSinkUrl)
	} else {
		viper.Set("influxdb", fmt.Sprintf("${%#[1]s_SERVICE_HOST}:${%#[1]s_SERVICE_PORT}", InfluxdbServiceName))
	}
	viper.Set("influxdb", os.ExpandEnv(viper.GetString("influxdb")))
	if kadvisorCmdV.PersistentFlags().Lookup("influxdb-secure").Changed {
		viper.Set("influxdbSecure", InfluxdbSecure)
	}

	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("config", viper.AllSettings()).Debug("Configured settings")
}

func InitializeKubeClient() *kube_client.Client {
	if kubernetesClient != nil {
		return kubernetesClient
	}

	kubeConfig := kube_client.Config{
		Host:     viper.GetString("kubernetesMaster"),
		Version:  viper.GetString("kubernetesApiVersion"),
		Insecure: viper.GetBool("kubernetesInsecure"),
	}

	clientAuthFile := viper.GetString("kubernetesClientAuthFile")
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
		log.Fatal(err)
	}
	return kubernetesClient
}
