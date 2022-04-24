TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=mavimo
NAME=drone
BINARY=terraform-provider-${NAME}
VERSION=0.4.0-dev
OS_ARCH=darwin_amd64

default: install

.PHONY: build release docs install test testacc

build:
	go build -o ${BINARY}

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

docs:
	go generate

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

# Run acceptance tests
testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m