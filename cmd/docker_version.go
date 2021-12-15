package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"time"
)

func Version() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	// 10 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	clientVersion := cli.ClientVersion()
	fmt.Println(clientVersion)
	serverVersion, err := cli.ServerVersion(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(serverVersion)
}
