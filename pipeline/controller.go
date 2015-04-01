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

package pipeline

import "sync"

var pipelineChan chan interface{}
var stopMutex sync.RWMutex

func Start(sinkChannels []chan interface{}, wg *sync.WaitGroup) chan interface{} {
	stopMutex.Lock()
	defer stopMutex.Unlock()
	pipelineChan = make(chan interface{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range pipelineChan {
			for _, sinkChannel := range sinkChannels {
				sinkChannel <- msg
			}
		}
	}()

	return pipelineChan
}

func Stop() {
	stopMutex.Lock()
	defer stopMutex.Unlock()
	close(pipelineChan)
}
