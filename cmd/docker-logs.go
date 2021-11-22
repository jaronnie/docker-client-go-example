package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
	"time"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"))
	if err != nil {
		panic(err)
	}
	// 10 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	logs, err := cli.ContainerLogs(ctx, "ha219aa123", types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow: true,
	})
	defer logs.Close()
	if err != nil {
		return
	}
	writer, err := io.Copy(os.Stdout,logs)
	if err != nil {
		panic(err)
	}
	fmt.Println(writer)
}