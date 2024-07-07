# build

## 手动编译：

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o whisky main.go
```

## 上传文件到 openwrt

```shell
scp -O whisky root@192.168.8.8:/root
```

## gzip 压缩解压文件

### 压缩文件

```
cd pkgs/cmds

gzip yacd_mode.sh 
```

### 压缩文件并保留源文件

```
cd pkgs/cmds

gzip -k yacd_mode.sh
```

### 解压文件

```
cd pkgs/cmds

gzip -d yacd_mode.sh.gz
```

## 关于 shell 脚本执行

### 对于纯脚本文件

直接定义为 `string` 类型，调用执行，示例代码：

```
command:="echo 'hello'"
cmd := exec.Command("bash", "-c", command)
```

### 对于脚本文件需要传入参数

将脚本文件定义为 `[]byte` 类型，同时使用 `gzip` 压缩。在执行时先将文件解压后保存到本地的临时文件目录下，然后通过命令方式调用执行并传入参数，示例代码：

```
command := fmt.Sprintf("%s %s", shellFilePath, strings.Join(args, " "))
cmd := exec.Command("/bin/sh", "-c", command)
```

---
