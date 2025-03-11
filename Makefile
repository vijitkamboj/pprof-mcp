mock:
	go run cmd/mock-server/main.go

profiler:
	go run cmd/profiler/main.go

tester:
	go run cmd/test-profiler/main.go > build/profiler.json

build-profiler:
	go build -o build/profiler cmd/server/server.go

install-latest-profiler:
	 go install github.com/vijitkamboj/pprof-mcp/cmd/profiler@latest