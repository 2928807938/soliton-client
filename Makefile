.PHONY: build run test clean tidy docker-up docker-down

# 构建
build:
	go build -o bin/api ./cmd/api

# 运行
run:
	go run ./cmd/api/main.go

# 测试
test:
	go test -v ./...

# 测试覆盖率
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 清理
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# 同步依赖
tidy:
	cd bom && go mod tidy
	cd share && go mod tidy
	cd user/domain && go mod tidy
	cd user/infrastructure && go mod tidy
	cd user && go mod tidy
	cd api/user-api && go mod tidy
	cd api && go mod tidy
	cd cmd/api && go mod tidy
	go work sync

# 启动 Docker 服务
docker-up:
	docker-compose up -d

# 停止 Docker 服务
docker-down:
	docker-compose down

# 查看 Docker 日志
docker-logs:
	docker-compose logs -f

# 重新构建并启动
docker-rebuild:
	docker-compose up -d --build
