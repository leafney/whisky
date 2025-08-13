# whisky

前后端一体化项目，后端使用 Go + Fiber，前端使用 React + TypeScript + Vite。

## 项目结构

```
whisky/
├── frontend/              # 前端项目 (React + TypeScript + Vite)
├── backend/               # 后端服务 (Go + Fiber)
│   └── web/              # 前端静态文件嵌入
│       ├── web.go        # 静态文件嵌入逻辑
│       └── dist/         # 前端构建产物目录
└── docs/                 # 项目文档
```

## 环境要求

- Go 1.22+
- Bun (推荐) 或 Node.js 18+
- curl

## 开发指南

### 安装依赖

```bash
# 安装后端依赖
cd backend && go mod tidy

# 安装前端依赖（当前端项目创建后）
cd frontend && bun install
```

### 本地开发

```bash
# 后端开发模式（开发时从文件系统读取静态文件）
cd backend && go run -tags dev main.go

# 前端开发模式（当前端项目创建后）
cd frontend && bun run dev
```

### 构建部署

```bash
# 构建前端（当前端项目创建后）
cd frontend && bun run build

# 复制前端构建产物到后端
mkdir -p backend/web/dist
rm -rf backend/web/dist/*
cp -R frontend/dist/* backend/web/dist/

# 构建后端（嵌入静态文件）
cd backend && go build -o ../bin/whisky main.go

# 使用 goreleaser 构建多平台版本
goreleaser build --snapshot --clean
```

### 发布流程

1. 提交代码到 Git
2. 创建版本标签：`git tag v1.0.0`
3. 推送标签：`git push origin v1.0.0`
4. GitHub Actions 会自动构建并发布

## 创建前端项目

```bash
# 在 frontend 目录中创建 Vite + React + TypeScript 项目
cd frontend
bun create vite . --template react-ts
bun install

# 或者使用 make 命令
make create-frontend
```

## 说明

- 生产环境下，前端静态文件会通过 `//go:embed` 嵌入到 Go 二进制文件中
- 开发环境下（使用 `-tags dev`），会从文件系统读取静态文件，便于前端热更新
- GitHub Actions 会自动处理前后端的构建和集成