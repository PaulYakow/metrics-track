.PHONY: cover
cover:
	go test -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

benchmarks:
	go test ./benchmarks -bench BenchmarkV1URL -memprofile benchmarks\profiles\base1.pprof

build_linter:
	go build -o ./cmd/staticlint/staticlint.exe ./cmd/staticlint/main.go ./cmd/staticlint/exitchecker.go

run_linter:
	./cmd/staticlint/staticlint.exe ./...