FROM gliderlabs/alpine:3.1
ENTRYPOINT ["/bin/kadvisor"]
VOLUME /mnt/routes
EXPOSE 8000

COPY . /go/src/github.com/fabric8io/kadvisor
RUN apk-install go git mercurial \
  && cd /go/src/github.com/fabric8io/kadvisor \
  && export GOPATH=/go \
  && export PATH=$GOPATH/bin:$PATH \
  && go get github.com/tools/godep \
  && godep go build -ldflags "-X main.Version $(cat VERSION)" -o /bin/kadvisor \
  && rm -rf /go \
  && apk del go git mercurial
