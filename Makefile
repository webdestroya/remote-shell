build:
	docker build -t cloud87/remote-shell:latest .

buildtest:
	docker build -t cloud87/shelltest:latest -f Dockerfile.usage .

runtest:
	docker run --rm -p 8722:8722 cloud87/shelltest:latest /cloud87/remote_shell_init -u webdestroya