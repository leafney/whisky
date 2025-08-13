# whisky 项目文档

欢迎查阅 whisky 项目的技术文档。

## 📚 文档索引

### 开发相关

- **[Makefile 命令文档](./makefile.md)** - 详细的 Makefile 命令说明和使用指南
  - 开发命令：install, server, wire, web, dev
  - 构建命令：build-front, build, build-*平台*, release
  - 工具命令：clean

### 构建部署

- **[构建说明](./build.md)** - 项目构建和编译相关说明
  - Makefile 构建（推荐）
  - 手动编译方法
  - 文件上传和压缩
  - Shell 脚本执行

- **[服务配置](./service.md)** - OpenWrt 系统下的服务部署
  - 下载和安装
  - 配置文件设置
  - 自启动服务配置
  - 服务管理命令

### 功能模块

- **[网络监控](./network-monitor.md)** - 智能网络连通性监控系统
  - 基于定时任务的自动检测
  - 多重保护机制和智能重启
  - 配置说明和API接口
  - 故障排查和最佳实践

## 🚀 快速开始

### 开发环境

```bash
# 1. 安装依赖
make install

# 2. 启动开发服务（两个终端）
make server  # 终端1：后端服务
make web     # 终端2：前端服务
```

### 构建部署

```bash
# 构建本地版本
make build

# 构建 Linux 版本用于服务器部署
make build-linux

# 构建所有平台版本
make release
```

## 📋 项目结构

```
whisky/
├── frontend/              # 前端项目 (React + TypeScript + Vite)
├── backend/               # 后端服务 (Go + Fiber)
│   └── web/              # 前端静态文件嵌入
│       ├── web.go        # 静态文件嵌入逻辑
│       └── dist/         # 前端构建产物目录
├── bin/                  # 构建输出目录
├── docs/                 # 项目文档
├── Makefile              # 构建工具
└── README.md             # 项目说明
```

## 🔗 相关链接

- [项目主页](../README.md) - 项目概述和功能介绍
- [GitHub 仓库](https://github.com/leafney/whisky) - 源代码和问题反馈
- [发布页面](https://github.com/leafney/whisky/releases) - 下载最新版本

---

*文档维护：请在修改项目功能时同步更新相关文档*
