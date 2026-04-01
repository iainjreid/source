PROG=source

.PHONY: $(PROG)

$(PROG):
	go build -ldflags="-s -w" -o $@
	du -sh $@
