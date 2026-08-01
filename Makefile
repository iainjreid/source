# Installation paths (configurable).
PREFIX  ?= /usr/local
DESTDIR ?=

# Go build flags.
GOFLAGS := -buildvcs=true -ldflags="-s -w" -gcflags="-m=2" -trimpath -tags=standalone

# Frontend build markers.
FRONTEND_DEPS_STAMP  := .frontend_deps.stamp
FRONTEND_BUILD_STAMP := .frontend_build.stamp

# Sources that are used to determine when the frontend should be rebuilt.
FRONTEND_SOURCES := $(shell find public -type f -not -path "public/dist/*")

.PHONY: src clean install

# Build the Source binary.
src: $(FRONTEND_DEPS_STAMP) $(FRONTEND_BUILD_STAMP)
	go build $(GOFLAGS) -o $@ ./cmd/src
	@echo "built: $$(du -sh $@)"

# Install the Source binary.
install:
	install -Dm755 src $(DESTDIR)$(PREFIX)/bin/src

# Install the frontend dependencies.
$(FRONTEND_DEPS_STAMP): package.json package-lock.json
	npm ci --ignore-scripts
	@touch $@

# Build the frontend assets.
$(FRONTEND_BUILD_STAMP): $(FRONTEND_DEPS_STAMP) $(FRONTEND_SOURCES)
	npm run build
	@touch $@

# Remove all artefacts.
clean:
	rm -f $(FRONTEND_DEPS_STAMP)
	rm -rf ./node_modules

	rm -f $(FRONTEND_BUILD_STAMP)
	rm -rf ./public/dist

	rm -f ./src
