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

FROM centos:7
ENTRYPOINT ["/bin/kadvisor"]
EXPOSE 80

ENV KADVISOR_VERSION 0.1

RUN yum install -y tar && \
    yum clean all && \
    curl -L https://github.com/jimmidyson/kadvisor/releases/download/v${KADVISOR_VERSION}/kadvisor-${KADVISOR_VERSION}-linux-amd64.tar.gz | \
      tar xzv
