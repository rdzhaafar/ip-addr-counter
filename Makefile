BIN=naive sorted-slice sorted-slice-dedup ip-store-concurrent ip-store
TEST=test-naive test-sorted-slice test-sorted-slice-dedup test-ip-store-concurrent-cmd test-ip-store-cmd test-ip-store-pkg

.PHONY: all
all: $(BIN)

.PHONY: clean
clean:
	rm -rf $(BIN)

.PHONY: test
test: $(TEST)

.PHONY: test-ip-store-cmd
test-ip-store-cmd:
	go test github.com/rdzhaafar/ip-addr-counter/cmd/ip-store

.PHONY: test-ip-store-pkg
test-ip-store-pkg:
	go test github.com/rdzhaafar/ip-addr-counter/pkg/ipstore

.PHONY: test-naive
test-naive:
	go test github.com/rdzhaafar/ip-addr-counter/cmd/naive

.PHONY: test-sorted-slice
test-sorted-slice:
	go test github.com/rdzhaafar/ip-addr-counter/cmd/sorted-slice

.PHONY: test-sorted-slice-dedup
test-sorted-slice-dedup:
	go test github.com/rdzhaafar/ip-addr-counter/cmd/sorted-slice-dedup

.PHONY: test-ip-store-concurrent-cmd
test-ip-store-concurrent-cmd:
	go test github.com/rdzhaafar/ip-addr-counter/cmd/ip-store-concurrent

naive: cmd/naive/main.go
	go build -o $@ $?

sorted-slice: cmd/sorted-slice/main.go
	go build -o $@ $?

sorted-slice-dedup: cmd/sorted-slice-dedup/main.go
	go build -o $@ $?

ip-store-concurrent: cmd/ip-store-concurrent/main.go
	go build -o $@ $?

ip-store: cmd/ip-store/main.go
	go build -o $@ $?
