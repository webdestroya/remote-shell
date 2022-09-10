build:
	docker build -t cloud87/remote-shell:latest .

buildtest:
	docker build -t cloud87/shelltest:latest -f Dockerfile.usage .

runtest:
	# docker run --rm -p 8722:8722 --privileged -it cloud87/shelltest:latest /bin/bash
	docker run --rm -p 8722:8722 --privileged cloud87/shelltest:latest /cloud87/bin/remote_shell -user webdestroya

compile:
	go build -a -o remote_shell

compile-version:
	go build -ldflags="-X 'main.Version=v1'" -a -o remote_shell

gomod:
	go mod tidy
