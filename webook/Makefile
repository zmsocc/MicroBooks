.PHONY: docker
docker:
	@del webook 2>NUL || exit 0
	@set GOOS=linux&& set GOARCH=arm&& set CGO_ENABLED=0&& go build -o webook .
	@docker rmi -f flycash/webook:v0.0.1
	@docker build -t flycash/webook:v0.0.1 .