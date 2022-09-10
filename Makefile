build:
	docker build -t cloud87/remote-shell:latest .

buildtest:
	docker build -t cloud87/shelltest:latest -f Dockerfile.usage .

runtest:
	# docker run --rm -p 8722:8722 --privileged -it cloud87/shelltest:latest /bin/bash
	docker run --rm -p 8722:8722 --privileged cloud87/shelltest:latest /cloud87/remote_shell_init -u webdestroya

go-docker-build:
	docker build -t c87rsgo:latest -f Dockerfile.golang .

gobuild:
	# go build -a remote_shell.go
	go build -o remote_shell

gotest: gobuild
	rm -rf tmp
	mkdir -p tmp
	./remote_shell -u webdestroya -h tmp

gomod:
	go mod tidy

buildssh:
	go build -a sshtest.go

# sshtest: buildssh
# 	./sshtest