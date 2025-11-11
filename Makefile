test:
	go test -v -run Test
.PHONY: test

test-cover:
	go test -v -coverprofile=cover.out
	go tool cover -html cover.out
.PHONY: test-cover

example:
	go test -v -run Example
.PHONY: example

bench:
	go test -bench=. -benchmem
.PHONY: bench

bench-cpu:
	go test -bench=. -benchmem -cpuprofile=cpu.out
	go tool pprof -http=:8080 cpu.out
.PHONY: bench-cpu

bench-mem:
	go test -bench=. -benchmem -memprofile=mem.out
	go tool pprof -http=:8080 mem.out
.PHONY: bench-mem

all: test example bench
.PHONY: all
