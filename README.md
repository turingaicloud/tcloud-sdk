# TCLOUD-SDK
## Command-line Interface used for TACC development.
```
TCLOUD Command-line Interface v0.0.2

Usage:
tcloud [command] [flags] [args]

Available Commands:
    tcloud init
    tcloud config [-u/-f] [args]
    tcloud download [<url>]
    tcloud add
    tcloud submit
    tcloud ps [-j] [<JOB_ID>]
    tcloud install

Use "tcloud [command] --help" for more information about a command.
```

## Installation
You can try out the latest features by directly install from master branch:

```
git clone https://github.com/xcwanAndy/tcloud-sdk
cd tcloud-sdk
echo 'export GOPATH=$PWD' >> ~/.bash_profile
make install
```

## Configuration
### CLI Configuration
Before using the tcloud CLI and submit ML jobs to TACC, you need to configure your TACC credentials. You can do this by running the `tcloud config` command:
```
$ tcloud config [-u/--username] MYUSERNAME
$ tcloud config [-f/--file] MYFILEPATH
```

### Job Configuration
TODO(SECTION in TUXIV.CONF)

TODO(TACC VARIABLES)

## Example
Basic examples are provided under the [example](example) folder. These examples include: [TensorFlow](example/TensorFlow), [PyTorch](example/PyTorch) and [MXNet](example/MXNet).