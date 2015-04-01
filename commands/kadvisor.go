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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jimmidyson/kadvisor/api"
	"github.com/jimmidyson/kadvisor/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	kube_client "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
)

var KadvisorCmd = &cobra.Command{
	Use:   "kadvisor",
	Short: "KAdvisor is a metrics collector & publisher for Kubernetes",
	Long:  "A configurable metrics collector & publisher for Kubernetes",
}

var kadvisorCmdV *cobra.Command

var (
	verbose             bool
	cfgFile             string
	defaultPollDuration time.Duration
	metricsSources      api.Uris
	metricsSinks        api.Uris
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
	KadvisorCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	KadvisorCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")
	KadvisorCmd.PersistentFlags().DurationVar(&defaultPollDuration, "default-poll", 10*time.Second, "poll duration")
	KadvisorCmd.PersistentFlags().Var(&metricsSources, "source", "sources")
	KadvisorCmd.PersistentFlags().Var(&metricsSinks, "sink", "sinks")

	kadvisorCmdV = KadvisorCmd
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func InitializeConfig() {
	if len(cfgFile) > 0 {
		viper.SetConfigFile(cfgFile)
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
	viper.SetDefault("defaultPoll", 10*time.Second)

	if kadvisorCmdV.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("verbose", verbose)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("default-poll").Changed {
		viper.Set("defaultPoll", defaultPollDuration)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("source").Changed {
		viper.Set("sources", metricsSources)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("sink").Changed {
		viper.Set("sinks", metricsSinks)
	}

	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("config", viper.AllSettings()).Debug("Configured settings")
}
