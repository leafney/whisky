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
```

## 构建命令
```bash
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
```

## 架构设计
- **后端**: Go 1.23+ 配合 Fiber web框架，Wire 依赖注入，Koanf 配置管理，LevelDB 数据库
- **前端**: React 19 + TypeScript 配合 Vite 构建工具，使用 Bun 包管理
- **关键目录**:
  - `backend/` - Go 后端代码 (internal/api 路由层, internal/biz 业务层, internal/dao 数据层)
  - `frontend/` - React 前端代码
  - `bin/` - 构建输出目录
  - `docs/` - 项目文档

## 开发工作流
1. 后端运行在 8080 端口
2. 前端开发服务器在 5173 端口，支持热模块更新(HMR)
3. 配置文件位于 `backend/config/config.toml.default`
4. 生产环境中静态文件嵌入到 Go 二进制文件中

## 测试与质量保障
- **后端**: `cd backend && go test ./...`
- **前端**: `cd frontend && npm run lint` + `npm run build`
- **代码质量**: Go fmt/vet 工具，ESLint 代码检查，TypeScript 严格类型模式

## 目标平台
- OpenWrt (主要目标)
- Linux AMD64/ARM64
- Windows
- macOS

## 关键说明
- 构建产物输出到 `bin/` 目录
- 所有依赖通过 Makefile 管理
- 生产环境构建将静态文件嵌入到二进制中
- 专为资源受限的嵌入式环境设计