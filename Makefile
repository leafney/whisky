.PHONY: install server wire web build-front dev build build-windows build-linux build-arm64 build-darwin release clean help

# 定义可执行文件名称和版本信息
EXECUTABLE = whisky
BIN_DIR = bin
BACKEND_DIR = backend

# 获取版本信息
DEFAULT_VERSION = $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
VERSION ?= $(DEFAULT_VERSION)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME = $(shell date +"%Y-%m-%d %H:%M:%S")

# 构建标志
LDFLAGS = -s -w -X 'main.Version=$(VERSION)' -X 'main.GitBranch=$(GIT_BRANCH)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'

# 默认目标
help:
	@echo "whisky 项目构建工具"
	@echo ""
	@echo "开发命令："
	@echo "  install       安装所有依赖"
	@echo "  server        启动后端程序"
	@echo "  wire          运行wire生成依赖文件"
	@echo "  web           启动前端项目"
	@echo "  dev           同时启动前后端"
	@echo ""
	@echo "构建命令："
	@echo "  build-front   构建前端项目"
	@echo "  build         构建完整项目(本地平台)"
	@echo "  build-windows 构建Windows版本"
	@echo "  build-linux   构建Linux版本"
	@echo "  build-arm64   构建ARM64版本"
	@echo "  build-darwin  构建macOS版本"
	@echo "  release       使用goreleaser构建"
	@echo ""
	@echo "工具命令："
	@echo "  clean         清理构建产物"

# 安装依赖
install:
	@echo "安装后端依赖..."
	cd $(BACKEND_DIR) && go mod tidy
	@if [ -f "frontend/package.json" ]; then \
		echo "安装前端依赖..."; \
		cd frontend && bun install; \
	else \
		echo "前端项目尚未创建，跳过前端依赖安装"; \
	fi

# 启动后端程序
server:
	@echo "启动后端服务器..."
	cd $(BACKEND_DIR) && go run main.go

# 运行wire生成依赖文件
wire:
	@echo "运行Wire生成依赖注入文件..."
	cd $(BACKEND_DIR)/cmd && wire

# 启动前端项目
web:
	@if [ -f "frontend/package.json" ]; then \
		echo "启动前端开发服务器..."; \
		cd frontend && bun run dev; \
	else \
		echo "前端项目尚未创建，请先创建前端项目"; \
		echo "运行: cd frontend && bun create vite . --template react-ts"; \
	fi

# 构建前端项目
build-front:
	@if [ -f "frontend/package.json" ]; then \
		echo "构建前端项目..."; \
		cd frontend && bun run build && cd ..; \
		echo "复制前端构建产物到后端..."; \
		mkdir -p $(BACKEND_DIR)/web/dist; \
		rm -rf $(BACKEND_DIR)/web/dist/*; \
		if [ -d "frontend/dist" ] && [ "$$(ls -A frontend/dist)" ]; then \
			cp -R frontend/dist/* $(BACKEND_DIR)/web/dist/; \
			echo "前端构建产物复制完成"; \
		else \
			echo "前端构建产物为空"; \
		fi; \
		echo "前端构建完成！"; \
	else \
		echo "前端项目尚未创建，使用默认页面"; \
		mkdir -p $(BACKEND_DIR)/web/dist; \
	fi

# 同时启动前后端 (提示用户在不同终端运行)
dev:
	@echo "启动开发环境..."
	@echo "请在不同的终端窗口中运行以下命令："
	@echo "  终端1: make server"
	@echo "  终端2: make web"
	@echo ""
	@echo "或者使用后台方式启动："
	@echo "  make server &"
	@echo "  make web"

# 通用构建函数
define build-target
	@echo "构建 $(1) 平台..."
	mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && GOOS=$(2) GOARCH=$(3) CGO_ENABLED=0 go build -o ../$(BIN_DIR)/$(EXECUTABLE)_$(1) \
		-ldflags="$(LDFLAGS)" \
		./main.go
	@echo "构建完成: $(BIN_DIR)/$(EXECUTABLE)_$(1)"
endef

# 构建完整项目 (本地平台)
build: build-front
	@echo "构建后端项目..."
	mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(BIN_DIR)/$(EXECUTABLE) \
		-ldflags="$(LDFLAGS)" \
		./main.go
	@echo "构建完成！可执行文件位于: $(BIN_DIR)/$(EXECUTABLE)"

# 构建Windows版本
build-windows: build-front
	$(call build-target,windows_amd64.exe,windows,amd64)

# 构建Linux版本
build-linux: build-front
	$(call build-target,linux_amd64,linux,amd64)

# 构建ARM64版本
build-arm64: build-front
	$(call build-target,linux_arm64,linux,arm64)

# 构建macOS版本
build-darwin: build-front
	$(call build-target,darwin_amd64,darwin,amd64)

# 使用goreleaser构建
release: build-front
	@echo "使用goreleaser构建..."
	goreleaser build --snapshot --clean

# 清理构建产物
clean:
	@echo "清理构建产物..."
	rm -rf $(BIN_DIR)/ dist/ $(BACKEND_DIR)/web/dist/
	cd $(BACKEND_DIR) && go clean
	@if [ -f "frontend/package.json" ]; then \
		cd frontend && rm -rf dist/; \
	fi
	@echo "清理完成！"