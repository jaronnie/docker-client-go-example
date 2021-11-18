# docker-client-go

## linux-centos7

**/usr/lib/systemd/system/docker.service**，配置远程访问。主要是在[Service]这个部分，加上下面两个参数

```shell
[Service] 

ExecStart= 

ExecStart=/usr/bin/dockerd -H tcp://0.0.0.0:2375 -H unix://var/run/docker.sock
```

```shell
systemctl daemon-reload
systemctl restart docker
```

## MacOs

```shell
brew install socat
socat TCP-LISTEN:2375,reuseaddr,fork UNIX-CONNECT:/var/run/docker.sock &
```


