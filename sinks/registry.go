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

var registry = make(map[string](func(string, map[string][]string) (Sink, error)))

func Register(uriPrefix string, sinkFunc func(string, map[string][]string) (Sink, error)) {
	registry[uriPrefix] = sinkFunc
}

func Lookup(uriPrefix string) (func(string, map[string][]string) (Sink, error), bool) {
	sink, ok := registry[uriPrefix]
	if !ok {
		return nil, ok
	}
	return sink, ok
}
