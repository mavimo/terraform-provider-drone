TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=mavimo
NAME=drone
BINARY=terraform-provider-${NAME}
VERSION=0.5.0-dev
OS_ARCH=darwin_amd64

default: install

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: release
release:
	goreleaser release --clean --snapshot --skip=publish,sign

.PHONY: docs
docs:
	go generate

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: create-env
create-env:
	docker compose -f ".github/workflows/acceptance-tests/docker-compose.yaml" down
	docker volume prune --force
	docker compose -f ".github/workflows/acceptance-tests/docker-compose.yaml" up -d

.PHONY: load-fixtures
load-fixtures:
	# Import Gitea data
	docker compose -f ".github/workflows/acceptance-tests/docker-compose.yaml" exec gitea-mysql mysql -uroot -pgitea gitea -e 'SOURCE /tmp/gitea-data.sql'
	# Create repositories
	curl -k -X POST "http://terraform:terraform@localhost:3000/api/v1/user/repos" -H "content-type: application/json" --data '{"name":"repo-test-1"}'
	curl -k -X POST "http://terraform:terraform@localhost:3000/api/v1/user/repos" -H "content-type: application/json" --data '{"name":"repo-test-2"}'
	curl -k -X POST "http://terraform:terraform@localhost:3000/api/v1/user/repos" -H "content-type: application/json" --data '{"name":"repo-test-3"}'
	# Import Drone data
	docker compose -f ".github/workflows/acceptance-tests/docker-compose.yaml" exec drone-mysql mysql -uroot -pdrone drone -e 'SOURCE /tmp/drone-data.sql'
	# Force Drone to sync repositories
	curl -XPOST -H 'Authorization: Bearer 5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL' http://localhost:8000/api/user/repos

.PHONY: test
test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

.PHONY: testacc
testacc:
	TF_ACC=1 DRONE_SERVER=http://localhost:8000/ DRONE_TOKEN=5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL DRONE_USER=terraform go test $(TEST) -v $(TESTARGS) -timeout 120m -coverprofile=coverage.out -race -covermode=atomic