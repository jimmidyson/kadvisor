# Copyright 2015 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

NAME=kadvisor
VERSION=$(shell cat VERSION)

local: *.go **/*.go
	godep go build -ldflags "-X main.Version dev" -o build/kadvisor

dev:
	@docker build -f Dockerfile.dev -t $(NAME):dev .
	@docker run --rm \
		-v /var/run/docker.sock:/tmp/docker.sock \
		-v $(PWD):/go/src/github.com/fabric8io/$(NAME)\
		-p 8000:8000 \
		$(NAME):dev

build:
	mkdir -p build
	docker build -t $(NAME):$(VERSION) .

release:
	rm -rf release && mkdir release
	go get github.com/progrium/gh-release/...
	cp build/* release
	gh-release create fabric8/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

clean:
	rm -rf build
	docker rmi $(NAME):dev $(NAME):$(VERSION) || true

.PHONY: release clean build
