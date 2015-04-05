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
	"github.com/jimmidyson/kadvisor/extpoints"
	"github.com/jimmidyson/kadvisor/sinks/api"
)

func init() {
	extpoints.SinkFactories.Register(New, "influxdb")
}

type InfluxdbSink struct {
}

func New(uri string, options map[string][]string) (api.Sink, error) {
	return &InfluxdbSink{}, nil
}

func (k *InfluxdbSink) Start(wg *sync.WaitGroup) chan interface{} {

	metricsChan := make(chan interface{})

	go func() {
		defer wg.Done()
		for msg := range metricsChan {
			log.WithFields(log.Fields{"sink": "influxdb", "msg": msg}).Debug("Received message")
		}
	}()

	return metricsChan
}
