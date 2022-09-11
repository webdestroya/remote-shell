# Remote SSH Docker Shell

Allows you to connect to a docker instance of your app.

## Authentication
This will pull your public keys from GitHub and use those for SSH authentication.

## Requirements
* Operating Systems: Debian, Alpine
* Architectures: AMD64, ARM64

## Configuration
| Flag  | Description | 
| ------------- | ------------- |
| `-user`  | The GitHub user to pull keys for<br>**Required** | 
| `-port`  | The remote port for the SSH server<br>Default: `8722` |
| `-shell`  | The shell command to execute.<br>Default: `/bin/bash` or `/bin/sh` |
| `-idletime`  | If the connection is idle for more than X seconds, terminate the connection. Setting to `0` disables.<br>Default: `0` (disabled) |
| `-maxtime`  | Maximum duration of a session.<br>Default: `12h` |
| `-grace`  | How long to wait for a connection before we just terminate.<br>Default: `30m` |
| `-insecure` | If you do not have CA Certificates installed, you can bypass SSL verification.<br>Not Recommended in production<br>Default: `false`

Note: Any of the arguments can be provided using environment variables by prefixing the flag with `C87RS_` (i.e. `C87RS_PORT`)

## Docker Image Variants

Each of the variants listed below will have the following scheme for `<version>`:
* `v#` - updated to the latest version of this major release
* `v#.#` - updated to the latest version of this minor release
* `v#.#.#` - will not change, is locked to this specific tag.
* `latest` - updated to the latest release

Additionally, images are multi-architecture and are available for the following platforms:
* `linux/amd64`
* `linux/arm64`


### `remote-shell:<version>`
This is the primary image, build on Ubuntu and will be compatible with nearly all Linux flavors.

### `remote-shell:<version>-alpine`
This image contains a binary that was built on Alpine Linux. If you are using an Alpine based image, use this variant.


## Usage

In your `Dockerfile`:

```docker
FROM anything
# ... 

# This can go anywhere in your Dockerfile
COPY --from=ghrc.io/webdestroya/remote-shell:latest /cloud87 /cloud87
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