# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 提供在操作此代码库时的指导信息。

## 项目概览
**Whisky** 是一个专门针对 OpenWrt 和嵌入式系统的全栈监控系统。采用 Go 后端和 React 前端技术栈，提供网络监控、系统状态追踪和管理功能。

## 快速启动命令
```bash
# 安装所有依赖
make install

# 开发模式 - 后端 (localhost:8080)
make server

# 开发模式 - 前端 (localhost:5173)  
make web

# 完整开发环境，需要在不同终端分别运行
make dev  # 显示启动指引
```

## 构建和测试命令
```bash
# 依赖注入代码生成
make wire

# 构建前端项目
make build-front

# 本地构建
make build

# 跨平台构建
make build-linux      # Linux AMD64
make build-windows    # Windows x64  
make build-arm64      # Linux ARM64
make build-darwin     # macOS AMD64
make release          # 使用 GoReleaser 构建所有平台

# 清理构建产物
make clean

# 测试和代码质量
cd backend && go test ./...          # 运行 Go 测试
cd backend && go fmt ./...           # Go 代码格式化
cd backend && go vet ./...           # Go 静态检查
cd frontend && bun run lint          # 前端 ESLint 检查
cd frontend && bun run build         # 前端构建测试
```

## 架构设计

### 后端分层架构 (基于 Wire 依赖注入)
- **API 层** (`internal/api/`): HTTP 路由处理器，包含 Home、Router、Network、YAcd、SCrash、CronTask 等 API
- **Business 层** (`internal/biz/`): 业务逻辑处理
- **DAO 层** (`internal/dao/`): 数据访问对象  
- **Service 层** (`internal/service/`): 内部服务，包含网络监控、定时任务、路由管理等
- **Model 层** (`internal/vmodel/`): 视图模型定义

### 技术栈
- **后端**: Go 1.23+ + Fiber web框架 + Wire 依赖注入 + Koanf 配置管理 + LevelDB 数据库
- **前端**: React 19 + TypeScript + Vite 构建工具 + Bun 包管理
- **工具包**: rose 日志框架、gocron 定时任务、gjson JSON 解析

### 依赖注入配置
- Wire 配置位于 `backend/internal/wire.go` 和 `backend/cmd/wire.go`
- 修改依赖结构后需运行 `make wire` 重新生成注入代码
- 主入口通过 `cmd.BuildInjector()` 初始化所有依赖

### 关键目录结构
```
backend/
├── cmd/           # 应用启动和依赖注入
├── internal/      # 内部业务代码 (API/Biz/DAO/Service 四层架构)
├── pkg/           # 可复用工具包
├── config/        # 配置管理
├── web/          # 前端静态文件嵌入
└── data/         # 运行时数据目录
```

## 开发工作流
1. 后端运行在 8080 端口 (可通过 `-p` 参数修改)
2. 前端开发服务器在 5173 端口，支持热模块更新(HMR)
3. 配置文件位于 `backend/config/config.toml.default` (运行时复制到 `backend/data/config.toml`)
4. 生产环境中静态文件通过 `//go:embed` 嵌入到 Go 二进制文件中
5. 开发环境使用 `-tags dev` 从文件系统读取静态文件

## 配置和部署
- 默认配置支持网络监控、LevelDB 存储、定时任务等功能
- 支持命令行参数: `-p` (端口), `-y` (yacd端口), `-w` (webhook), `-d` (debug模式)
- 版本信息通过构建时的 ldflags 注入

## 已知问题
- 后端存在格式化问题需要修复 (`go vet` 报错)
- API 层部分日志格式化指令缺失

## 目标平台
- OpenWrt (主要目标)
- Linux AMD64/ARM64  
- Windows
- macOS