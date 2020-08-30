# TCLOUD-SDK
## Command-line Interface used for TACC development.
```
TCLOUD Command-line Interface v0.0.1

Usage:
tcloud [command]

Available Commands:
    tcloud upload
    tcloud download
    tcloud init
    tcloud clone
    TODO[add more]

Use "tcloud [command] --help" for more information about a command.
```

Build instructions (Linux)
* Project has to be cloned to somewhere in your `$GOPATH`
* To install your `$GOBIN` has to be set. 
```
$ cd $TCLOUD_DIR/cli
$ make build install
```

## XCompile     [TODO]
Support Linux / macOS / Windows.

## Example  [TODO]
Tcloud examples, which includes: helloworld, TensorFlow, PyTorch, MXNet.