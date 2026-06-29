PROG=src

GOFLAGS=-buildvcs=true -ldflags="-s -w" -trimpath -tags="standalone"

.PHONY: $(PROG)

FRONTEND_DEPS_STAMP := .frontend_deps.stamp
FRONTEND_BUILD_STAMP := .frontend_build.stamp

$(PROG): $(FRONTEND_DEPS_STAMP) $(FRONTEND_BUILD_STAMP)
	go build $(GOFLAGS) -o $@ ./cmd/src
	du -sh $@

$(FRONTEND_DEPS_STAMP): package.json package-lock.json
	npm ci
	@touch $@

$(FRONTEND_BUILD_STAMP): $(FRONTEND_DEPS_STAMP) $(shell find public -type f -not -path "public/dist/*")
	npm run build
	@touch $@

test:
	go test ./...

fmt:
	go fmt ./...

clean:
	rm -f $(FRONTEND_DEPS_STAMP)
	rm -rf ./node_modules

	rm -f $(FRONTEND_BUILD_STAMP)
	rm -rf ./public/dist

	rm -f ./src
