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
	"os"
	"os/signal"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/jimmidyson/kadvisor/pipeline"
	"github.com/jimmidyson/kadvisor/sinks"
	"github.com/jimmidyson/kadvisor/sources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	if !viper.IsSet("sinks") {
		log.Fatal("Exiting: No sinks specified")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		sources.Stop()
		pipeline.Stop()
		sinks.Stop()
	}()

	var sinkWg sync.WaitGroup
	sinksMetricsChannels := sinks.Start(&sinkWg)

	var pipelineWg sync.WaitGroup
	pipelineChan := pipeline.Start(sinksMetricsChannels, &pipelineWg)

	var sourceWg sync.WaitGroup
	sources.Start(&sourceWg, pipelineChan)

	sourceWg.Wait()
	log.Debug("Sources completed")
	pipelineWg.Wait()
	log.Debug("Pipelines completed")
	sinkWg.Wait()
	log.Debug("Sinks completed")
}
