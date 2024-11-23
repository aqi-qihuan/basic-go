.PHONY: docker
docker:
	@rm lanmengbook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o lanmengbook .
	@docker rmi -f flycash/lanmengbook:v0.0.1
	@docker build -t flycash/lanmengbook:v0.0.1 .
