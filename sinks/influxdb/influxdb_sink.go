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

package influxdb

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/fabric8io/kadvisor/sinks"
	"github.com/tuxychandru/pubsub"
)

func init() {
	sinks.Register("influxdb", New)
}

type InfluxdbSink struct {
}

func New(uri string) (sinks.Sink, error) {
	return &InfluxdbSink{}, nil
}

func (k *InfluxdbSink) Start(pubSub *pubsub.PubSub, wg *sync.WaitGroup) {
	metricsSub := pubSub.Sub("metrics")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range metricsSub {
			log.WithFields(log.Fields{"sink": "influxdb", "msg": msg}).Debug("Received message")
		}
	}()
}