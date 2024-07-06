# build

## 手动编译：

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o whisky main.go
```

## 上传文件到 openwrt

```shell
scp -O whisky root@192.168.8.8:/root
```
