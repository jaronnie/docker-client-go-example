package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Images() {
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"))
	if err != nil {
		panic(err)
	}
	list, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, v := range list {
		fmt.Println(v.RepoTags, v.Created, v.ID)
	}
}