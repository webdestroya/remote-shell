build:
	docker build --progress plain -t webdestroya/remote-shell:latest -f docker/Dockerfile .

buildtest:
	docker build -t shelltest:latest -f docker/Dockerfile.usage .

runtest:
	docker run --rm -p 8722:8722 --privileged cloud87/shelltest:latest /cloud87/bin/remote-shell -user webdestroya

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