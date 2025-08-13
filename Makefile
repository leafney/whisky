.PHONY: backend frontend build clean dev-backend dev-frontend build-frontend install help

# 默认目标
help:
	@echo "whisky 项目构建工具"
	@echo ""
	@echo "可用命令："
	@echo "  install        安装所有依赖"
	@echo "  dev-backend    启动后端开发服务器"
	@echo "  dev-frontend   启动前端开发服务器"
	@echo "  build-frontend 构建前端项目"
	@echo "  build-backend  构建后端项目（包含前端）"
	@echo "  build-release  使用goreleaser构建"
	@echo "  clean          清理构建产物"

# 安装依赖
install:
	@echo "安装后端依赖..."
	cd backend && go mod tidy
	@if [ -f "frontend/package.json" ]; then \
		echo "安装前端依赖..."; \
		cd frontend && bun install; \
	else \
		echo "前端项目尚未创建，跳过前端依赖安装"; \
	fi

# 后端开发模式（带dev标签，从文件系统读取静态文件）
dev-backend:
	@echo "启动后端开发服务器..."
	cd backend && go run -tags dev main.go

# 前端开发模式
dev-frontend:
	@if [ -f "frontend/package.json" ]; then \
		echo "启动前端开发服务器..."; \
		cd frontend && bun run dev; \
	else \
		echo "前端项目尚未创建，请先创建前端项目"; \
		echo "运行: cd frontend && bun create vite . --template react-ts"; \
	fi

# 构建前端
build-frontend:
	@if [ -f "frontend/package.json" ]; then \
		echo "构建前端项目..."; \
		cd frontend && bun run build && cd ..; \
		echo "复制前端构建产物到后端..."; \
		mkdir -p backend/web/dist; \
		rm -rf backend/web/dist/*; \
		if [ -d "frontend/dist" ] && [ "$$(ls -A frontend/dist)" ]; then \
			cp -R frontend/dist/* backend/web/dist/; \
			echo "前端构建产物复制完成"; \
		else \
			echo "前端构建产物为空"; \
		fi; \
		echo "前端构建完成！"; \
	else \
		echo "前端项目尚未创建，使用默认页面"; \
		mkdir -p backend/web/dist; \
	fi

# 构建后端（生产模式，嵌入静态文件）
build-backend: build-frontend
	@echo "构建后端项目..."
	cd backend && go build -o ../bin/whisky main.go
	@echo "构建完成！可执行文件位于: bin/whisky"

# 使用goreleaser构建
build-release: build-frontend
	@echo "使用goreleaser构建..."
	goreleaser build --snapshot --clean

# 清理构建产物
clean:
	@echo "清理构建产物..."
	rm -rf bin/ dist/ backend/web/dist/
	cd backend && go clean
	@if [ -f "frontend/package.json" ]; then \
		cd frontend && rm -rf dist/; \
	fi
	@echo "清理完成！"


