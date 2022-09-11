# Remote SSH Docker Shell

Allows you to connect to a docker instance of your app.

## Authentication
This will pull your public keys from GitHub and use those for SSH authentication.

## Requirements
* Operating Systems: Debian, Alpine
* Architectures: AMD64, ARM64

## Usage

In your `Dockerfile`:

```docker
FROM ruby:3.1.2
# ...

# This can go anywhere in your image (to improve caching)
COPY --from=ghrc.io/webdestroya/remote-shell:v1 /cloud87 /cloud87

# For Alpine builds:
COPY --from=ghrc.io/webdestroya/remote-shell:v1 /cloud87-alpine /cloud87
```

Then you can launch the container:

```
$ docker run --rm -p 8722:8722 myapp:latest /cloud87/bin/remote-shell -user your-github-username
```

And connect to it on your client:

```
ssh -p 8722 \
  -o StrictHostKeyChecking=no \
  -o "UserKnownHostsFile=/dev/null" \
  IP_OF_THE_CONTAINER
```


## Configuration
| Option  | Arg | Env Var |   |
| ------------- | ------------- | ------------- | ------------- |
| GitHub User  | `-user`  | `C87RS_USER` | The user to pull keys for<br>**Required** | 
| Port  | `-port`  | `C87RS_PORT` | The remote port for the SSH server<br>Default: 8722 |
| Shell  | `-shell`  | `C87RS_SHELL` | The shell command to execute.<br>Default: `/bin/bash` or `/bin/sh` |
| Idle Timeout  | `-idletime`  | `C87RS_IDLETIME` | If the connection is idle for more than X seconds, terminate the connection.<br>`0` disables.<br>Default: disabled |
| Max Runtime  | `-maxtime`  | `C87RS_MAXTIME` | Maximum duration of a session.<br>Default: 12 hours |
| Grace Time  | `-grace`  | `C87RS_GRACE` | How long to wait for a connection before we just terminate.<br>Default: 10min |
| Insecure Mode | `-insecure` | _none_ | If you do not have CA Certificates installed, you can bypass SSL verification.<br>Not Recommended in production<br>Default: `false`

