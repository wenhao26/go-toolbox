# Video Parser 构建脚本
# 支持 Windows、Linux、macOS 多平台交叉编译

BINARY_NAME = videoparser
BUILD_DIR = build
CMD_DIR = ./cmd

# 版本信息
VERSION ?= 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go 编译参数
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# 默认目标
.PHONY: all
all: clean build

# 构建所有平台
.PHONY: build
build: build-windows build-linux build-darwin
	@echo "✓ All platforms built successfully"

# Windows (amd64)
.PHONY: build-windows
build-windows:
	@echo "Building Windows amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "✓ Windows build complete"

# Linux (amd64)
.PHONY: build-linux
build-linux:
	@echo "Building Linux amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	@echo "✓ Linux build complete"

# macOS (amd64)
.PHONY: build-darwin
build-darwin:
	@echo "Building macOS amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	@echo "✓ macOS build complete"

# 仅构建当前平台（开发调试用）
.PHONY: build-local
build-local:
	@echo "Building for current platform..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✓ Local build complete"

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "✓ Clean complete"

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "✓ Tests complete"

# 代码格式化
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Format complete"

# 帮助信息
.PHONY: help
help:
	@echo "Video Parser Build Commands:"
	@echo "  make all          - Clean and build all platforms"
	@echo "  make build        - Build all platforms"
	@echo "  make build-local  - Build for current platform only"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format code"
	@echo "  make help         - Show this help"