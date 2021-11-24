package cmd

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"
)

func Start() {
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"))
	if err != nil {
		panic(err)
	}
	// 10 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = cli.ContainerStart(ctx, "308a0bcdebc3", types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
}