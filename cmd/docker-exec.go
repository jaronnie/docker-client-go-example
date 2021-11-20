package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"
)

/*
reference: https://blog.csdn.net/q1403539144/article/details/115967275
 */

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"))
	if err != nil {
		panic(err)
	}
	// 10 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// create a exec command
	create, err := cli.ContainerExecCreate(ctx, "hae9b59699", types.ExecConfig{
		Tty: true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"ls"},
	})
	if err != nil {
		panic(err)
	}
	// attach a exec
	attach, err := cli.ContainerExecAttach(ctx, create.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(attach.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}