test:
	go test -v -run Test

test-cover:
	go test -v -coverprofile=cover.out
	go tool cover -html cover.out

example:
	go test -v -run Example

bench:
	go test -bench=. -benchmem

bench-cpu:
	go test -bench=. -benchmem -cpuprofile=cpu.out
	go tool pprof -http=:8080 cpu.out

bench-mem:
	go test -bench=. -benchmem -memprofile=mem.out
	go tool pprof -http=:8080 mem.out
