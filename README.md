# Remote SSH Docker Shell

Allows you to connect to a docker instance of your app.

## Authentication
This will pull your public keys from GitHub and use those for SSH authentication.

## Requirements
This app assumes that the main docker image contains the following binaries:

* `bash`
* `cat`
* `chmod`
* `curl`
* `grep`
* `mkdir`
* `printenv`
* `sed`
* `sort`
* `timeout`
* `xargs`

If you do not have these, you can always write a custom init script to configure dropbear as needed.

## Usage

In your `Dockerfile`:

```docker
FROM ruby:3.1.2
# ...

# This can go anywhere in your image (to improve caching)
COPY --from=ghrc.io/webdestroya/docker-remote-shell:v1 /cloud87 /cloud87

```

Then you can launch the container:

```
$ docker run --rm -p 8722:8722 myapp:latest /cloud87/remote_shell -u your-github-username
```


## Configuration
| Option  | Arg | Env Var |   |
| ------------- | ------------- | ------------- | ------------- |
| GitHub User  | `-u`  | `C87_RSHELL_GH_USER` | The user to pull keys for<br>**Required** | 
| Port  | `-p`  | `C87_RSHELL_PORT` | The remote port for the SSH server<br>Default: 8722 |
| Idle Timeout  | `-i`  | `C87_RSHELL_IDLE_TIMEOUT` | If the connection is idle for more than X seconds, terminate the connection.<br>`0` disables.<br>Default: disabled |
| Keepalive  | `-k`  | `C87_RSHELL_KEEPALIVE` |keepalive interval.<br>`0` disables.<br>Default: 5 min |
| Max Runtime  | `-m`  | `C87_RSHELL_MAX_RUNTIME` | Maximum duration of a session.<br>Default: 12 hours |

