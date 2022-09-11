AWS_ACCOUNTID=$$(aws sts get-caller-identity --query Account --output text)

build:
	docker build --progress plain -t ghcr.io/webdestroya/remote-shell:latest -f docker/Dockerfile .

buildtest:
	docker build -t shelltest:latest -f docker/Dockerfile.usage .

runtest:
	docker run --rm -p 8722:8722 --privileged shelltest:latest /cloud87/bin/remote-shell -user webdestroya

runversion:
	docker run --rm shelltest:latest /cloud87/bin/remote-shell -version

compile:
	rm -f remote-shell
	go build -a -o remote-shell

compile-version:
	go build -ldflags="-X 'main.buildVersion=v1'" -a -o remote-shell

gomod:
	go mod tidy

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.49.0 golangci-lint run -v

pushtest: build buildtest
	aws ecr get-login-password | docker login --username AWS --password-stdin $(AWS_ACCOUNTID).dkr.ecr.us-east-1.amazonaws.com
	docker tag shelltest:latest $(AWS_ACCOUNTID).dkr.ecr.us-east-1.amazonaws.com/cloud87/remote-shell-test:latest
	docker push $(AWS_ACCOUNTID).dkr.ecr.us-east-1.amazonaws.com/cloud87/remote-shell-test:latest
