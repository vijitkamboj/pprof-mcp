mock:
	go run cmd/mock/main.go

profiler:
	go run cmd/profiler/main.go

build-profiler:
	go build -o build/profiler cmd/profiler/main.go

tester:
	go run cmd/tester/main.go > build/profiler.json

build-mcp:
	go build -o build/mcp cmd/mcp/main.go