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

### TensorFlow

+ Dataset: mnist

+ Task: image classification

+ Code: [mnist.py](https://github.com/xcwanAndy/tcloud-sdk/blob/master/examples/TuXiv_example/mnist.py)

+ Configuration

  + TACC ENV

    ~~~shell
    TACC_WORKERDIR #repo directory
    ~~~

  + TuXiv configuration

    ~~~yaml
    # tuxiv.conf
    entrypoint:
        - python ${TACC_WORKDIR}/mnist.py 
        - --task_index=0
        - --data_dir=${TACC_WORKDIR}/datasets/mnist_data
        - --batch_size=1
    environment:
        name: tf 
        dependencies:
            - tensorflow=1.15
    job:
        name: mnist
        general:
            - nodes=2
    ~~~

+ Training process:

  + Enter the `TACC_WORKDIR` directory and follow the steps.
  + Build environment: `tcloud build tuxiv.conf`
  + Submit job: `tcloud submit`
  + Monitor job: `tcloud show [job id]`
  + Cancel job: `tcloud cancel [job id]`

