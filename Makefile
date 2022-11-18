test_coverage:
	go test ./... -coverprofile=coverage.out

benchmarks:
	go test ./benchmarks -bench BenchmarkV1URL -memprofile benchmarks\profiles\base1.pprof
	go test ./benchmarks -bench BenchmarkV2URL -memprofile benchmarks\profiles\base2.pprof