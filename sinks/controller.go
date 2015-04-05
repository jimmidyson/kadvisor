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

package sinks

import (
	"net/url"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/jimmidyson/kadvisor/api"
	"github.com/jimmidyson/kadvisor/extpoints"
	"github.com/spf13/viper"
)

var metricsChannels []chan interface{}
var stopChannelsMutex sync.RWMutex

func Start(wg *sync.WaitGroup) []chan interface{} {
	if !viper.IsSet("sinks") {
		log.Fatal("Exiting: No sinks specified")
	}

	sinkUris := viper.Get("sinks")
	log.WithField("uris", sinkUris).Debug("Creating all sinks")

	stopChannelsMutex.Lock()
	defer stopChannelsMutex.Unlock()
	for _, sink := range sinkUris.(api.Uris) {
		log.WithField("uri", sink).Debug("Creating sink")
		u, err := url.Parse(sink)
		if err != nil {
			log.WithField("uri", sink).Fatal("Unparseable sink RL")
		}
		sinkType := u.Scheme
		sinkUrl := u.Opaque
		options := u.Query()
		if len(sinkUrl) == 0 {
			log.WithFields(log.Fields{"url": sink}).Fatal("Invalid sink configuration")
		}
		log.WithFields(log.Fields{"type": sinkType, "url": sinkUrl}).Debug("Parsed sink URL")
		sinkFunc := extpoints.SinkFactories.Lookup(sinkType)
		if sinkFunc == nil {
			log.WithField("type", sinkType).Fatal("Unregistered sink type")
		}
		sink, err := sinkFunc(sinkUrl, options)
		if err != nil {
			log.WithFields(log.Fields{"type": sinkType, "url": sinkUrl, "error": err}).Fatal("Could not create sink")
		}
		wg.Add(1)
		metricsChannels = append(metricsChannels, sink.Start(wg))
	}

	return metricsChannels
}

func Stop() {
	stopChannelsMutex.Lock()
	defer stopChannelsMutex.Unlock()
	for _, metricsChannels := range metricsChannels {
		close(metricsChannels)
	}
}
