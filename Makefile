PROG=src

GOFLAGS=-buildvcs=true -ldflags="-s -w" -trimpath -tags="standalone"

.PHONY: $(PROG)

$(PROG):
	go build $(GOFLAGS) -o $@ ./cmd/src
	du -sh $@
