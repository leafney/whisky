## 设置在 openwrt 系统下开机启动

### 下载文件

将二级制文件通过 `wget` 命令下载到指定目录，并设置可执行权限：

```
wget <下载链接> -O /usr/sbin/whisky

chmod +x /usr/sbin/whisky
```

或者通过 `scp` 命令上传文件：

```
scp -O whisky root@192.168.8.8:/usr/sbin

chmod +x /usr/sbin/whisky
```

### 修改配置

根据需求，自定义修改配置文件。

最终得到：
- 二进制文件：`/usr/sbin/whisky`
- 配置文件：`/etc/whisky/config.yaml`

### 编写自启动服务

新增或编辑文件 `/etc/init.d/whisky` ，内容：

```shell
#!/bin/sh /etc/rc.common

USE_PROCD=1
START=90

start_service(){
    procd_open_instance [whisky]
    procd_set_param command /usr/sbin/whisky
    # procd_append_param command -c /etc/whisky/config.yaml

    procd_set_param respawn ${respawn_threshold:-3600} ${respawn_timeout:-5} ${respawn_retry:-5}
    procd_set_param env LANG=zh_CN.UTF-8
    procd_set_param limits nofile="unlimited"
    procd_set_param stdout 1
    procd_set_param stderr 1
    # procd_set_param user nobody
    procd_close_instance

}
```

注意修改其中二进制文件路径和配置文件路径。

### 启动服务

为启动文件赋予可执行权限：

```
chmod +x /etc/init.d/whisky
```

启动服务：

```
/etc/init.d/whiksy start
```

设置服务开机启动：

```
/etc/init.d/whisky enable
```

#### 管理服务

支持通过命令的方式管理启动项服务：

- `start` 启动服务
- `stop` 停止服务
- `restart` 重启服务
- `relaod` 重载配置文件
- `enable` 设置开机自启动
- `disable` 禁用开机自启动

---
