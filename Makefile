# 定义可执行文件的名称
EXECUTABLE = whisky

# 定义输出目录
DIST_DIR = ./dist

# 定义不同平台的可执行文件名称
WINDOWS = $(EXECUTABLE)_windows_amd64.exe
LINUX = $(EXECUTABLE)_linux_amd64
DARWIN = $(EXECUTABLE)_darwin_amd64
ARM64 = $(EXECUTABLE)_linux_arm64

# 获取当前版本为最近一次提交的短哈希值
DEFAULT_VERSION = $(shell git rev-parse --short HEAD)
VERSION ?= $(DEFAULT_VERSION)

GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_TIME = $(shell date +"%Y-%m-%d %H:%M:%S")

# 通用构建命令
define build
	@echo "Building $(TARGET)..."
	mkdir -p $(DIST_DIR)
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build -o $(DIST_DIR)/$(TARGET) \
		-ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.GitBranch=$(GIT_BRANCH)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'" \
		./main.go
endef

# 定义构建目标
windows: TARGET = $(WINDOWS)
windows:
	$(call build,windows,amd64)

linux: TARGET = $(LINUX)
linux:
	$(call build,linux,amd64)

darwin: TARGET = $(DARWIN)
darwin:
	$(call build,darwin,amd64)

arm64: TARGET = $(ARM64)
arm64:
	$(call build,linux,arm64)

# 定义清理目标
clean:
	rm -rf $(DIST_DIR)

wire:
	@echo "Running Wire..."
	cd $(shell pwd)/cmd && wire

run:
	@echo "Go Start..."
	go run main.go

# 定义默认目标
all: windows linux darwin arm64

# 添加提示信息
help:
	@echo "使用方法:"
	@echo "  make         - 构建所有平台的可执行文件"
	@echo "  make clean   - 清理生成的可执行文件"
	@echo "  make windows - 仅构建 Windows 平台的可执行文件"
	@echo "  make linux   - 仅构建 Linux 平台的可执行文件"
	@echo "  make arm64   - 仅构建 Linux arm64 平台的可执行文件"
	@echo "  make darwin  - 仅构建 macOS 平台的可执行文件"
	@echo "  make VERSION=<自定义版本> - 使用自定义版本构建"
	@echo "  make wire    - 运行 wire 生成依赖注入文件"
	@echo "  make run     - 运行 Go 程序"

# 设置默认目标为 help
.DEFAULT_GOAL := help

.PHONY: all clean windows linux darwin help arm64 wire run
