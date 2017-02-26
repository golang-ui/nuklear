all:
	c-for-go nk.yml

clean:
	rm -f nk/cgo_helpers.go nk/cgo_helpers.h nk/cgo_helpers.c
	rm -r nk/doc.go nk/types.go nk/const.go
	rm -f nk/nk.go

test:
	cd nk && go build
	