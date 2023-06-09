# go-template
Golang project template - jgreat style.

### Config

Configuration options are read from environment variables.
| Variable | Default | Type | Description |
| --- | --- | --- | --- |
| `LISTEN_HOST` | `127.0.0.1` | string | Set the prometheus metrics listen host interface address |
| `LISTEN_PORT` | `9090` | int | Set the prometheus metrics listen port |
| `LOG_LEVEL` | `info` | string | Log level |


## Development

See [Setup](#go_setup) for requirements and creating a go development environment.

### Doing development

When you start work, source the go project `source_me.sh` file to setup your environment.

```
cd ~/gopath/%REPO%
source ./source_me.sh
```

Use `docker-compose` to build/run the app with sample DB and follow the logs:

```
cd ~/gopath/%REPO%/src/github.com/mobilecoinofficial/%REPO%
docker-compose up --build
```

The app available on localhost.

* App endpoint: http://127.0.0.1:9090/metrics

When you make changes to the code the app container should automatically restart and rebuild the app binary.

### Debugging

When launched from `docker-compose` the app is run with headless `dlv` listening on `:2345`.

There is a pre-configured `.vscode/launch.json` profile ready to attach to `dlv` to remote debug.

In `vscode` select the debug option and run `Attach Remote`, set your break points and have fun.

---

## Go Setup

### Prerequisites

This environment is pre-configured for running/compiling in docker with remote debugging and automatic rebuilds on code changes.

- `vscode` - https://code.visualstudio.com/download
- `go` 1.20 - https://golang.org/dl/
- `docker` - https://docs.docker.com/get-docker/
- `docker-compose` - https://docs.docker.com/compose/install/

### Install go

Download latest go 1.20 for your system: https://golang.org/dl/

Extract tar to `~/bin` (this will over wite the contents of the current `go`)

```
cd ~/bin
tar xvzf ~/Downloads/go1.20.3.linux-amd64.tar.gz
```

Move `go` to a versioned directory

```
mv go go-1.20.3
```

### Set up development environment

These instructions will help you create an isolated project path in you home directory.  More details on gopath setup can be found here: https://golang.org/doc/gopath_code.html

```
mkdir -p ~/gopath/%REPO%
cd ~/gopath/%REPO%
```

Add this script to your project base and point `GOROOT` at the version of go:

`source_me.sh`

```
export GOROOT="${HOME}/bin/go-1.20.3"
export GOPATH="$(pwd)"
export PATH="${PATH}:${GOROOT}/bin:${GOPATH}/bin"
```

Create src path and clone the repo:

```
# change into the base directory and source the script
cd ~/gopath/%REPO%
. source_me.sh

# make the src directory with git path
mkdir -p src/github.com/mobilecoinofficial

# cd and clone the repo
cd src/github.com/mobilecoinofficial
git clone git@github.com:git@github.com:mobilecoinofficial/%REPO%.git
```
