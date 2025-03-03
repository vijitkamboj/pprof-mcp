# pprof-mcp

A Go profiling tool that integrates with AI agent through the MCP (Model Control Protocol) interface, allowing you to profile Go servers and get AI-powered analysis of the results.

## Overview

This tool connects to a running Go server's pprof endpoints, collects profiling data (heap, CPU, etc.), and provides the results in a format that the agent can interpret to give you insights about potential memory leaks, CPU usage patterns, and performance bottlenecks.

## Installation

### Option 1: Quick Install (for a system that does not has go compiler)

```bash
curl -fsSl https://raw.githubusercontent.com/vijitkamboj/pprof-mcp/main/run.sh | sh
```

This will:
1. Download the latest release of the profiler binary
2. Make it executable
3. Run it

### Option 2: Go Package Installation

```bash
go install github.com/vijitkamboj/pprof-mcp/cmd/profiler@latest
```

After installation, you can run the tool directly from your terminal:

```bash
profiler
```

## Usage

Once the profiler is running, it will register as a tool. You can then use it by asking the agent to profile golang server:

```
profile_server(host="localhost:6060")
```

Where `localhost:6060` is the address of your Go server with pprof endpoints enabled.


## Requirements

- A running Go server with pprof endpoints enabled
- Access to an agent with MCP capabilities

## Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/vijitkamboj/pprof-mcp.git
cd pprof-mcp

# Build the profiler
make build-profiler
```

### Available Make Commands

- `make mock` - Run the mock server
- `make profiler` - Run the profiler
- `make build-profiler` - Build the profiler binary
- `make tester` - Run tests and output results to a JSON file
- `make build-mcp` - Build the MCP binary

## License

[License information]

## Contributing

[Contribution guidelines]