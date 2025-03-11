package main

import (
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/vijitkamboj/pprof-mcp/profiler"
)

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport(), mcp_golang.WithVersion("0.1.0"))
	err := server.RegisterTool("profile_server", "Get the performance profile of a golang server [ 'always take host as input from user first']", profileHandler)
	if err != nil {
		panic(err)
	}

	// err = server.RegisterResource("test://resource", "resource_test", "This is a test resource", "application/json", func() (*mcp_golang.ResourceResponse, error) {
	// 	return mcp_golang.NewResourceResponse(mcp_golang.NewTextEmbeddedResource("test://resource", "This is a test resource", "application/json")), nil
	// })
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("Profiler server started")
	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}

type ProfileRequest struct {
	Host string `json:"host" jsonschema:"required,description="`
}

func profileHandler(ctx context.Context, request ProfileRequest) (*mcp_golang.ToolResponse, error) {

	pf := profiler.NewProfiler("localhost:6060", profiler.WithSamplers(
		profiler.NewHeapSampler(map[string]string{}),
		profiler.NewCPUSampler(map[string]string{
			"seconds": "5",
		}),
	))
	data := pf.RunAll()
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(data)), nil

}
