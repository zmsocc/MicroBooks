.PHONY: docker
docker:
	@del webook 2>NUL || exit 0
	@set GOOS=linux&& set GOARCH=arm&& set CGO_ENABLED=0&& go build -o webook .
	@docker rmi -f flycash/webook:v0.0.1
	@docker build -t flycash/webook:v0.0.1 .

.PHONY: grpc
# 生成protoc文件命令
grpc:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto
	@buf generate webook/api/proto