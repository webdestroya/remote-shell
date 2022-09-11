build:
	docker build --progress plain -t cloud87/remote-shell:latest .

buildtest:
	docker build -t cloud87/shelltest:latest -f Dockerfile.usage .

runtest:
	# docker run --rm -p 8722:8722 --privileged -it cloud87/shelltest:latest /bin/bash
	docker run --rm -p 8722:8722 --privileged cloud87/shelltest:latest /cloud87/bin/remote-shell -user webdestroya

runversion:
	docker run --rm cloud87/shelltest:latest /cloud87/bin/remote-shell -version

compile:
	go build -a -o remote_shell

compile-version:
	go build -ldflags="-X 'main.buildVersion=v1'" -a -o remote-shell

gomod:
	go mod tidy
