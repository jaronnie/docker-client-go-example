package cmd

import (
	"context"
	"github.com/docker/docker/client"
	"time"
)

func Stop() {
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"))
	if err != nil {
		panic(err)
	}
	// 10 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	timeout := time.Second*10
	err = cli.ContainerStop(ctx, "308a0bcdebc3", &timeout)
	if err != nil {
		panic(err)
	}
}
