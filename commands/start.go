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
	"net/url"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/fabric8io/kadvisor/sinks"
	"github.com/fabric8io/kadvisor/sources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tuxychandru/pubsub"
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

	if !viper.IsSet("sources") {
		log.Fatal("Exiting: No sources specified")
	}
	if !viper.IsSet("sinks") {
		log.Fatal("Exiting: No sinks specified")
	}

	pubSub := pubsub.New(0)
	defer pubSub.Shutdown()

	// Start sinks before sources so we can
	startSinks(pubSub)
	startSources(pubSub)
}

func startSources(pubSub *pubsub.PubSub) {
	log.WithField("uris", reflect.TypeOf(viper.Get("sources"))).Debug("Creating all sources")

	for _, source := range viper.Get("sources").(uris) {
		log.WithField("uri", source).Debug("Creating source")
		u, err := url.Parse(source)
		if err != nil {
			log.WithField("uri", source).Fatal("Unparseable source URL")
		}
		sourceType := u.Scheme
		sourceUrl := source[len(sourceType)+3:]
		log.WithFields(log.Fields{"type": sourceType, "url": sourceUrl}).Debug("Parsed source URL")
		sourceFunc, ok := sources.Lookup(sourceType)
		if !ok {
			log.WithField("type", sourceType).Fatal("Unregistered source type")
		}
		source, err := sourceFunc(sourceUrl)
		if err != nil {
			log.WithFields(log.Fields{"type": sourceType, "url": sourceUrl, "error": err}).Fatal("Could not create source")
		}
		source.Start(pubSub)
	}
}

func startSinks(pubSub *pubsub.PubSub) {
	log.WithField("uris", reflect.TypeOf(viper.Get("sinks"))).Debug("Creating all sinks")

	for _, sink := range viper.Get("sinks").(uris) {
		log.WithField("uri", sink).Debug("Creating sink")
		u, err := url.Parse(sink)
		if err != nil {
			log.WithField("uri", sink).Fatal("Unparseable sink RL")
		}
		sinkType := u.Scheme
		sinkUrl := sink[len(sinkType)+3:]
		log.WithFields(log.Fields{"type": sinkType, "url": sinkUrl}).Debug("Parsed sink URL")
		sinkFunc, ok := sinks.Lookup(sinkType)
		if !ok {
			log.WithField("type", sinkType).Fatal("Unregistered sink type")
		}
		sink, err := sinkFunc(sinkUrl)
		if err != nil {
			log.WithFields(log.Fields{"type": sinkType, "url": sinkUrl, "error": err}).Fatal("Could not create sink")
		}
		sink.Start(pubSub)
	}
}
