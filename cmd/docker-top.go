package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"time"
)

func Top() {
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"))
	if err != nil {
		panic(err)
	}
	for {
		top, err := cli.ContainerTop(context.Background(), "308a0bcdebc3", []string{})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second*3)
		fmt.Println(top)
	}
}
