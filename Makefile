
BUILD_DIR = build
SERVICES = service agent
CGO_ENABLED ?= 0
GOARCH ?= amd64

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -mod=vendor -ldflags "-s -w" -o ${BUILD_DIR}/mainflux-license-$(1) cmd/$(1)/main.go
endef

all: $(SERVICES)

.PHONY: all $(SERVICES) docker

$(SERVICES):
	$(call compile_service,$(@))

docker:
	docker build \
		--no-cache \
		--build-arg SVC=service \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--tag=mainflux/license-service \
		-f docker/Dockerfile .
