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
	docker save $(NAME):$(VERSION) | gzip -9 > build/$(NAME)_$(VERSION).tgz

release:
	rm -rf release && mkdir release
	go get github.com/progrium/gh-release/...
	cp build/* release
	gh-release create fabric8/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

clean:
	rm -rf build
	docker rmi $(NAME):dev $(NAME):$(VERSION) || true

circleci:
	rm ~/.gitconfig
ifneq ($(CIRCLE_BRANCH), release)
	echo build-$$CIRCLE_BUILD_NUM > VERSION
endif

.PHONY: release clean
