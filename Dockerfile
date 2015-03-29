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

FROM gliderlabs/alpine:3.1
ENTRYPOINT ["/bin/kadvisor"]
VOLUME /mnt/routes
EXPOSE 8000

COPY . /go/src/github.com/fabric8io/kadvisor
RUN apk-install go git mercurial gcc g++ \
  && cd /go/src/github.com/fabric8io/kadvisor \
  && export GOPATH=/go \
  && export PATH=$GOPATH/bin:$PATH \
  && go get github.com/tools/godep \
  && godep go build -ldflags "-X main.Version $(cat VERSION)" -o /bin/kadvisor \
  && rm -rf /go \
  && apk del go git mercurial gcc g++
