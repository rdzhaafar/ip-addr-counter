BIN=naive sorted-slice sorted-slice-dedup ip-store-concurrent ip-store ip-store-v2
TEST=test-naive test-sorted-slice test-sorted-slice-dedup test-ip-store-concurrent-cmd test-ip-store-cmd test-ip-store-pkg test-ip-store-v2-cmd test-ip-store-v2-pkg

.PHONY: all
all: $(TEST) $(BIN)

.PHONY: clean
clean:
	rm -rf $(BIN)

.PHONY: test
test: $(TEST)

.PHONY: test-ip-store-cmd
test-ip-store-cmd:
	go test github.com/rdzhaafar/ip-addr-counter/cmd/ip-store

.PHONY: test-ip-store-v2
test-ip-store-v2-cmd:
	go test github.com/rdzhaafar/ip-addr-counter/cmd/ip-store-v2

.PHONY: test-ip-store-pkg
test-ip-store-pkg:
	go test github.com/rdzhaafar/ip-addr-counter/pkg/ipstore/v1

.PHONY: test-ip-store-v2-pkg
test-ip-store-v2-pkg:
	go test github.com/rdzhaafar/ip-addr-counter/pkg/ipstore/v2

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

ip-store-v2: cmd/ip-store-v2/main.go
	go build -o $@ $?
