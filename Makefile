.PHONY: mock
# Makefile中的命令必须使用Tab缩进，而不是空格！！！
mock:
	@mockgen -source="webook/internal/repository/articles/article.go" -package="artrepomocks" -destination="webook/internal/repository/articles/mocks/article.mock.go"
	@mockgen -source="webook/internal/service/article.go" -package="artsvcmocks" -destination="webook/internal/service/mocks/article.mock.go"
	@go mod tidy

#.PHONY: grpc
#grpc:
#	@buf generate webook/api/proto