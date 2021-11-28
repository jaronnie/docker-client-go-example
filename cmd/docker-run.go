package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"math/rand"
	"mydocker/util"
	"time"
)

var DockerYAML = `
hostname: host-agent
image: gocloudcoder/kube-image-pull:develop
resources:
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
    - name: agent
      containerPort: 8888
      protocol: tcp  
`

type Config struct {
	Resources `json:"resources"`

	Hostname string `json:"hostname"`
	Image string `json:"image"`
}

type Resources struct {
	Limits 	`json:"limits"`
	Envs []Env `json:"env"`
	Ports []Port `json:"ports"`
}

type Limits struct {
	Cpu string `json:"cpu"`
	Memory string `json:"memory"`
}

type Env struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

type Port struct {
	Name string `json:"name"`
	ContainerPort int `json:"containerPort"`
	Protocol string `json:"protocol"`
}

func Run() {
	// generate client
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"))
	if err != nil {
		panic(err)
	}
	
	// pull image
	// 3 minutes timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	_, err = cli.ImagePull(ctx, "gocloudcoder/kube-image-pull:develop", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	
	// get config from yaml
	content, err := util.Yaml2Json([]byte(DockerYAML))
	if err != nil {
		panic(err)
	}
	
	// unmarshal config
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}

	// generate envs
	envs := SetEnv(config.Envs)

	// generate portBindings
	bindings := SetPortBindings(config.Ports)

	// create container
	createResp, err := cli.ContainerCreate(ctx, &container.Config{
		Hostname: config.Hostname,
		Image:    config.Image,
		// set env
		Env: envs,
	}, &container.HostConfig{
		// bind port
		PortBindings: bindings,
		// limit cpu and memory
		Resources: container.Resources{
			Memory: 104857600, // 100M
		},
	}, nil, nil, generateRandomHostname())
	if err != nil {
		panic(err)
	}

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

func SetEnv(envs []Env) []string {
	setEnvs := make([]string, 0)
	for _, v := range envs {
		setEnvs = append(setEnvs, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}
	return setEnvs
}

func SetPortBindings(ports []Port) map[nat.Port][]nat.PortBinding {
	// portMap
	portMap := make(map[nat.Port][]nat.PortBinding)
	for _, v := range ports {
		portMap[nat.Port(fmt.Sprintf("%d/%s", v.ContainerPort, v.Protocol))] = []nat.PortBinding {
			{
				HostIP:   "0.0.0.0",
				HostPort: "0", // 随机分配
			},
		}
	}
	return portMap
}