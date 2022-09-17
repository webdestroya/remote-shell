# Remote SSH Docker Shell

Allows you to connect to a container running your application. Once you exit the SSH session, the server exits and the container dies.

This is similar to Heroku's one-off `heroku run bash` command. It is not meant to connect to long running containers for your app, but rather an ephemeral instance of your app's image. This program is very similar to AWS ECS Exec, but without the idle timeout or need to install any extra programs on the client. It's pure SSH.

## Authentication
This will pull your public keys from GitHub and use those for SSH authentication.

## Requirements
* Operating Systems: Linux (both glibc/musl supported)
* Architectures: AMD64, ARM64

## Configuration
| Flag  | Description | 
| ------------- | ------------- |
| `user`  | The GitHub user to pull keys for<br>**Required unless keys provided via env var** | 
| `port`  | The remote port for the SSH server<br>Default: `8722` |
| `shell`  | The shell command to execute.<br>Default: `/bin/bash` or `/bin/sh` |
| `idletime`  | If the connection is idle for more than X seconds, terminate the connection. Setting to `0` disables.<br>Default: `0` (disabled) |
| `maxtime`  | Maximum duration of a session.<br>Default: `12h` |
| `grace`  | How long to wait for a connection before we just terminate.<br>Default: `30m` |
| `insecure` | If you do not have CA Certificates installed, you can bypass SSL verification.<br>Not Recommended in production<br>Default: `false`

> Note: Any of the arguments can be provided using environment variables by prefixing the flag with `C87RS_` (i.e. `C87RS_PORT`)

You can also provide a single SSH key via the environment variable: `C87_RSHELL_AUTHORIZED_KEY`. The value is the same format used in a normal authorized key file. (`ssh-rsa XXXXX`)

## Docker Image


#### Tagging Scheme
* `v#` - updated to the latest version of this major release
* `v#.#` - updated to the latest version of this minor release
* `v#.#.#` - will not change, is locked to this specific tag.
* `latest` - updated to the latest release

Images are multi-architecture and are available for the following platforms:
* `linux/amd64`
* `linux/arm64`


## Usage

In your `Dockerfile`:

```dockerfile
FROM anything
# ... 

# This can go anywhere in your Dockerfile
COPY --from=ghrc.io/webdestroya/remote-shell:latest /cloud87 /cloud87
```

Then you can launch the container:

```
$ docker run -p 8722:8722 myapp:latest /cloud87/bin/remote-shell -user your-github-username
```

And connect to it on your client:

```sh
ssh -p 8722 \
  -o StrictHostKeyChecking=no \
  -o "UserKnownHostsFile=/dev/null" \
  IP_OF_THE_CONTAINER
```