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

package sources

import (
	"net/url"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/jimmidyson/kadvisor/api"
	"github.com/spf13/viper"
)

var stopChannels []chan api.Stop
var stopChannelsMutex sync.RWMutex

func Start(wg *sync.WaitGroup, pipelineChan chan interface{}) {
	if !viper.IsSet("sources") {
		log.Fatal("Exiting: No sources specified")
	}

	sourceUris := viper.Get("sources")
	log.WithField("uris", sourceUris).Debug("Creating all sources")

	stopChannelsMutex.Lock()
	defer stopChannelsMutex.Unlock()
	for _, source := range sourceUris.(api.Uris) {
		log.WithField("uri", source).Debug("Creating source")
		u, err := url.Parse(source)
		if err != nil {
			log.WithField("uri", source).Fatal("Unparseable source URL")
		}
		sourceType := u.Scheme
		sourceUrl := u.Opaque
		options := u.Query()
		if len(sourceUrl) == 0 {
			log.WithFields(log.Fields{"url": source}).Fatal("Invalid source configuration")
		}
		log.WithFields(log.Fields{"type": sourceType, "url": sourceUrl, "options": options}).Debug("Parsed source URL")
		sourceFunc, ok := Lookup(sourceType)
		if !ok {
			log.WithField("type", sourceType).Fatal("Unregistered source type")
		}
		source, err := sourceFunc(sourceUrl, options)
		if err != nil {
			log.WithFields(log.Fields{"type": sourceType, "url": sourceUrl, "error": err}).Fatal("Could not create source")
		}
		stopChannel := source.Start(pipelineChan, wg)
		stopChannels = append(stopChannels, stopChannel)
	}
}

func Stop() {
	stopChannelsMutex.Lock()
	defer stopChannelsMutex.Unlock()
	for _, stopChannel := range stopChannels {
		stopChannel <- api.Stop{}
	}
}
