# TCLOUD-SDK
## Command-line Interface used for TACC development.
```
TCLOUD Command-line Interface v0.0.1

Usage:
tcloud [command]

Available Commands:
    tcloud init
    tcloud config
    tcloud download
    tcloud build
    tcloud add
    tcloud submit
    tcloud ps
    tcloud attach
    TODO[add more]

Use "tcloud [command] --help" for more information about a command.
```

## Build from source code
You can try out the latest features by directly install from master branch:

```
git clone https://github.com/xcwanAndy/tcloud-sdk
cd tcloud-sdk
echo 'export GOPATH=$PWD' >> ~/.bash_profile
make install
```

## XCompile     [TODO]
Support Linux / macOS / Windows.

## Example  [TODO]
Basic examples are provided under the [example](example) folder. These examples include: [helloworld](example/helloworld), [TensorFlow](example/TensorFlow), [PyTorch](example/PyTorch) and [MXNet](example/MXNet).

