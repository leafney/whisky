# Makefile 命令说明文档

本文档详细说明了 whisky 项目中 Makefile 的所有可用命令及其用法。

## 概述

whisky 项目使用统一的 Makefile 管理前后端开发、构建和部署流程。所有命令都在项目根目录下执行。

## 命令分类

### 📋 帮助命令

#### `make help`
显示所有可用命令的帮助信息。

```bash
make help
```

**输出示例：**
```
whisky 项目构建工具

开发命令：
  install       安装所有依赖
  server        启动后端程序
  wire          运行wire生成依赖文件
  web           启动前端项目
  dev           同时启动前后端

构建命令：
  build-front   构建前端项目
  build         构建完整项目(本地平台)
  build-windows 构建Windows版本
  build-linux   构建Linux版本
  build-arm64   构建ARM64版本
  build-darwin  构建macOS版本
  release       使用goreleaser构建

工具命令：
  clean         清理构建产物
```

---

## 🛠️ 开发命令

### `make install`
安装项目的所有依赖包。

```bash
make install
```

**功能：**
- 安装后端 Go 模块依赖 (`go mod tidy`)
- 检测并安装前端依赖 (`bun install`)

**适用场景：**
- 首次克隆项目后
- 依赖包更新后
- 环境重置后

---

### `make server`
启动后端开发服务器。

```bash
make server
```

**功能：**
- 在开发模式下运行后端程序
- 自动监听文件变化并重新加载
- 默认监听端口：8080

**等效命令：**
```bash
cd backend && go run main.go
```

**注意事项：**
- 确保已安装所有依赖
- 可通过命令行参数指定端口：`go run main.go -p 9000`

---

### `make wire`
运行 Google Wire 生成依赖注入代码。

```bash
make wire
```

**功能：**
- 扫描 `backend/cmd/` 目录下的 Wire 配置
- 生成依赖注入代码到 `wire_gen.go` 文件

**适用场景：**
- 修改依赖注入配置后
- 添加新的服务或组件后

**等效命令：**
```bash
cd backend/cmd && wire
```

---

### `make web`
启动前端开发服务器。

```bash
make web
```

**功能：**
- 启动 Vite 开发服务器
- 支持热模块替换 (HMR)
- 默认监听端口：5173

**前提条件：**
- 前端项目已创建 (`frontend/package.json` 存在)
- 已安装前端依赖

**等效命令：**
```bash
cd frontend && bun run dev
```

**如果前端项目未创建：**
```bash
cd frontend && bun create vite . --template react-ts
```

---

### `make dev`
同时启动前后端开发环境。

```bash
make dev
```

**功能：**
- 提供同时启动前后端的指导信息
- 建议在不同终端窗口中分别运行前后端服务

**推荐用法：**
```bash
# 终端1：启动后端
make server

# 终端2：启动前端
make web
```

**或后台启动：**
```bash
make server &
make web
```

---

## 🔨 构建命令

### `make build-front`
构建前端项目。

```bash
make build-front
```

**功能：**
- 使用 Vite 构建生产版本前端代码
- 将构建产物复制到 `backend/web/dist/` 目录
- 为后端嵌入静态文件做准备

**构建流程：**
1. 执行 `bun run build` 构建前端
2. 创建 `backend/web/dist/` 目录
3. 复制 `frontend/dist/*` 到 `backend/web/dist/`

**输出产物：**
- `backend/web/dist/index.html`
- `backend/web/dist/assets/` (CSS, JS, 图片等)

---

### `make build`
构建完整项目（本地平台）。

```bash
make build
```

**功能：**
- 先构建前端项目 (`make build-front`)
- 构建后端可执行文件
- 嵌入前端静态文件到后端二进制文件中

**构建信息：**
- 自动获取 Git 版本信息（分支、提交哈希、构建时间）
- 生成优化的二进制文件（去除调试信息）

**输出文件：**
- `bin/whisky` - 包含前后端的完整可执行文件

**等效命令：**
```bash
make build-front
cd backend && go build -o ../bin/whisky -ldflags="..." ./main.go
```

---

### `make build-windows`
构建 Windows 平台版本。

```bash
make build-windows
```

**功能：**
- 交叉编译生成 Windows 可执行文件
- 包含完整的前后端功能

**输出文件：**
- `bin/whisky_windows_amd64.exe`

**技术细节：**
- `GOOS=windows GOARCH=amd64`
- `CGO_ENABLED=0` (静态链接)

---

### `make build-linux`
构建 Linux 平台版本。

```bash
make build-linux
```

**功能：**
- 交叉编译生成 Linux 可执行文件
- 适用于服务器部署

**输出文件：**
- `bin/whisky_linux_amd64`

**技术细节：**
- `GOOS=linux GOARCH=amd64`
- `CGO_ENABLED=0` (静态链接)

---

### `make build-arm64`
构建 ARM64 平台版本。

```bash
make build-arm64
```

**功能：**
- 交叉编译生成 ARM64 可执行文件
- 适用于 ARM 架构服务器（如树莓派、Apple Silicon）

**输出文件：**
- `bin/whisky_linux_arm64`

**技术细节：**
- `GOOS=linux GOARCH=arm64`
- `CGO_ENABLED=0` (静态链接)

---

### `make build-darwin`
构建 macOS 平台版本。

```bash
make build-darwin
```

**功能：**
- 交叉编译生成 macOS 可执行文件

**输出文件：**
- `bin/whisky_darwin_amd64`

**技术细节：**
- `GOOS=darwin GOARCH=amd64`
- `CGO_ENABLED=0` (静态链接)

---

### `make release`
使用 GoReleaser 构建发布版本。

```bash
make release
```

**功能：**
- 先构建前端项目
- 使用 GoReleaser 构建多平台版本
- 生成发布包和变更日志

**前提条件：**
- 已安装 GoReleaser
- 配置了 `.goreleaser.yaml`

**等效命令：**
```bash
make build-front
goreleaser build --snapshot --clean
```

**输出目录：**
- `dist/` - GoReleaser 构建产物

---

## 🧹 工具命令

### `make clean`
清理所有构建产物。

```bash
make clean
```

**清理内容：**
- `bin/` - 本地构建的可执行文件
- `dist/` - GoReleaser 构建产物
- `backend/web/dist/` - 前端构建产物
- `frontend/dist/` - 前端构建缓存
- Go 构建缓存 (`go clean`)

**适用场景：**
- 构建出现问题时
- 切换分支前
- 释放磁盘空间

---

## 📝 使用示例

### 开发工作流

```bash
# 1. 安装依赖
make install

# 2. 启动开发环境（两个终端）
# 终端1
make server

# 终端2  
make web
```

### 构建工作流

```bash
# 构建本地版本
make build

# 构建多平台版本
make build-linux
make build-windows
make build-arm64

# 或使用 GoReleaser
make release
```

### 部署工作流

```bash
# 清理环境
make clean

# 构建生产版本
make build-linux

# 部署到服务器
scp bin/whisky_linux_amd64 user@server:/usr/sbin/whisky
```

---

## ⚙️ 环境变量

### 版本信息
Makefile 自动获取以下版本信息并嵌入到构建中：

- `VERSION` - Git 提交哈希 (可通过环境变量覆盖)
- `GIT_BRANCH` - 当前 Git 分支
- `GIT_COMMIT` - Git 提交哈希
- `BUILD_TIME` - 构建时间

### 自定义版本号

```bash
# 使用自定义版本号构建
VERSION=v1.0.0 make build
```

---

## 🔍 故障排除

### 常见问题

1. **前端项目未创建**
   ```
   前端项目尚未创建，请先创建前端项目
   运行: cd frontend && bun create vite . --template react-ts
   ```

2. **依赖安装失败**
   ```bash
   # 清理并重新安装
   make clean
   make install
   ```

3. **构建失败**
   ```bash
   # 检查 Go 模块
   cd backend && go mod tidy
   
   # 检查前端依赖
   cd frontend && bun install
   ```

4. **Wire 生成失败**
   ```bash
   # 确保已安装 wire
   go install github.com/google/wire/cmd/wire@latest
   ```

---

## 📚 相关文档

- [构建说明](./build.md) - 详细的构建和部署说明
- [服务配置](./service.md) - 服务部署和配置说明
- [项目README](../README.md) - 项目概述和快速开始

---

*最后更新：2025-08-13*
