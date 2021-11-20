package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"time"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"))
	if err != nil {
		panic(err)
	}
	// 3 minutes timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	inspect, err := cli.ContainerInspect(ctx, "0d3ddae80cc9f800acb85570371b604acd1867031d70cd9ae3a8ba0530019fd4")
	if err != nil {
		panic(err)
	}
	marshal, err := json.Marshal(inspect)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}
