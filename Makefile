##
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/:.*##/:/' | sed 's/^##//g'

.PHONY: build

build: ## Builds binaries for all platforms
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/aws-ec2-assign-elastic-ip-linux-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/aws-ec2-assign-elastic-ip-darwin-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o build/aws-ec2-assign-elastic-ip-linux-arm
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/aws-ec2-assign-elastic-ip-windows-amd64

release:
	aws s3 sync build/ s3://aws-ec2-assign-elastic-ip/`git rev-parse --short HEAD`/
