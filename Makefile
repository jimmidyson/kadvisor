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
	go generate
	godep go build -ldflags "-X main.Version dev" -o build/kadvisor

test: local
	godep go test $(shell godep go list ./...)

deps:
	go get -u github.com/progrium/go-extpoints

dev:
	@docker history $(NAME):dev &> /dev/null \
		|| docker build -f Dockerfile.dev -t $(NAME):dev .
	@docker run --rm \
		-v /var/run/docker.sock:/tmp/docker.sock \
		-v $(PWD):/go/src/github.com/jimmidyson/$(NAME)\
		-p 8000:8000 \
		$(NAME):dev

build:
	mkdir -p build
	docker build -t $(NAME):$(VERSION) .

release:
	rm -rf build release && mkdir build release
	go generate
	for os in linux freebsd darwin ; do \
		GOOS=$$os ARCH=amd64 godep go build -ldflags "-X main.Version $(VERSION)" -o build/kadvisor-$$os-amd64 ; \
		tar --transform 's|^build/||' --transform 's|-.*||' -czvf release/kadvisor-$(VERSION)-$$os-amd64.tar.gz build/kadvisor-$$os-amd64 README.md LICENSE ; \
	done
	GOOS=windows ARCH=amd64 godep go build -ldflags "-X main.Version $(VERSION)" -o build/kadvisor-$(VERSION)-windows-amd64.exe
	zip release/kadvisor-$(VERSION)-windows-amd64.zip build/kadvisor-$(VERSION)-windows-amd64.exe README.md LICENSE && \
		echo -e "@ build/kadvisor-$(VERSION)-windows-amd64.exe\n@=kadvisor.exe"  | zipnote -w release/kadvisor-$(VERSION)-windows-amd64.zip
	go get -u github.com/progrium/gh-release/...
	gh-release create jimmidyson/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

clean:
	rm -rf build relase
	docker rmi $(NAME):dev $(NAME):$(VERSION) || true

.PHONY: release clean build deps test
