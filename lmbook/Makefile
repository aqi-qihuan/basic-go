.PHONY: docker
docker:
	@rm lmbook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o lmbook .
	@docker rmi -f flycash/lmbook:v0.0.1
	@docker build -t flycash/lmbook:v0.0.1 .
