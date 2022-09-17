AWS_ACCOUNTID=$$(aws sts get-caller-identity --query Account --output text)

.PHONY: build
build:
	docker build --progress plain -t ghcr.io/webdestroya/remote-shell:latest -f docker/Dockerfile .

.PHONY: buildtest
buildtest:
	docker build -t shelltest:latest -f docker/Dockerfile.usage .

.PHONY: clean
clean:
	rm -f remote-shell

.PHONY: compile
compile: clean
	go build -a -o remote-shell -ldflags="-s -w"
	stat remote-shell

.PHONY: tidy
tidy:
	go mod verify
	go mod tidy
	@if ! git diff --quiet go.mod go.sum; then \
		echo "please run go mod tidy and check in changes, you might have to use the same version of Go as the CI"; \
		exit 1; \
	fi

# lint:
# 	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.49.0 golangci-lint run -v

.PHONY: lint-install
lint-install:
	@echo "Installing golangci-lint"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2

.PHONY: lint
lint:
	@which golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found, please run: make lint-install"; \
		exit 1; \
	}
	golangci-lint run

# This is for easy testing
.PHONY: pushtest
pushtest: build buildtest
	aws ecr get-login-password | docker login --username AWS --password-stdin $(AWS_ACCOUNTID).dkr.ecr.us-east-1.amazonaws.com
	docker tag shelltest:latest $(AWS_ACCOUNTID).dkr.ecr.us-east-1.amazonaws.com/cloud87/remote-shell-test:latest
	docker push $(AWS_ACCOUNTID).dkr.ecr.us-east-1.amazonaws.com/cloud87/remote-shell-test:latest


.PHONY: test-release
test-release:
	goreleaser release --skip-publish --rm-dist --snapshot --debug