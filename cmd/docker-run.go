package main

import (
	"context"
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"time"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"))
	if err != nil {
		panic(err)
	}
	// pull image
	// 3 minutes timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	reader, err := cli.ImagePull(ctx, "gocloudcoder/kube-image-pull:develop", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	writer, err := io.Copy(os.Stdout,reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(writer)
	// create container
	// portMap
	portMap := make(map[nat.Port][]nat.PortBinding)
	// bind tcp/8888(container) -> 0.0.0.0:28888(host)
	portMap["8888/tcp"] = []nat.PortBinding{
		{
			HostIP:   "0.0.0.0",
			HostPort: "0", // 随机分配
		},
	}
	createResp, err := cli.ContainerCreate(ctx, &container.Config{
		Hostname: "blocface-hostagent",
		Image:    "gocloudcoder/kube-image-pull:develop",
	}, &container.HostConfig{
		PortBindings: portMap,
	}, nil, nil, "ha"+uuid.Generate().String()[:8])
	if err != nil {
		panic(err)
	}
	fmt.Println(createResp.ID)
	// start container
	err = cli.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
}