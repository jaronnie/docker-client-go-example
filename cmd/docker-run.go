package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"math/rand"
	"os"
	"time"
)

var DockerYAML = `
name: host-agent
image: gocloudcoder/kube-image-pull:develop
resources:
  request:
    cpu: "100m"
    memory: "100Mi"
  limits:
    cpu: "100m"
    memory: "500Mi"
  env:
    - name: MachineID
      value: machineID
    - name: ENV_TYPE
      value: EnvType
    - name: AgentID
      value: AgentUUID
    - name: CORE_ADDR
      value: coreAddr
    - name: COREWS
      value: coreWSAddr
    - name: Pushgateway
      value: Pushgateway
    - name: Logger
      value: Logger
    - name: LoggerWS
      value: LoggerWS 
ports:
  - containerPort: 9009
    protocol: tcp  
`

type Config struct {
	Resources

	Name string
	Image string

}

type Resources struct {}


func Run() {
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
	_, err = io.Copy(os.Stdout,reader)
	if err != nil {
		panic(err)
	}
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
	}, nil, nil, generateRandomHostname())
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

// machine uuid
func generateRandomHostname() string {
	// 加入随机种子
	rand.Seed(time.Now().UnixNano())
	t := fmt.Sprintf("%d", time.Now().UnixNano())
	f := func() int64 { return rand.Int63n(9999-1000) + 1000 }
	return fmt.Sprintf("docker-%s-%d-%d", t[len(t)-4:], f(), f())
}