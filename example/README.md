# TCloud Examples

## Hello world

## TensorFlow

+ Dataset: mnist

+ Task: image classification

+ Code: [mnist.py](https://github.com/xcwanAndy/tcloud-sdk/blob/master/example/TensorFlow/mnist.py)

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

## PyTorch
+ Dataset: mnist

+ Task: image classification

+ Code: [mnist.py](https://github.com/xcwanAndy/tcloud-sdk/blob/master/example/Pytorch/mnist.py)

+ Configuration

  + TACC ENV

    ~~~shell
    TACC_WORKERDIR #repo directory
    ~~~

  + TuXiv configuration

    ~~~yaml
    # tuxiv.conf
    
    entrypoint:
        - python ${TACC_WORKDIR}/mnist.py --epoch=3
    environment:
        name: torch-env
        dependencies:
            - pytorch=1.6.0
            - torchvision=0.7.0
        channel: pytorch
    job:
        name: test
        general:
            - nodes=2
    ~~~

+ Training process:

  + Enter the `TACC_WORKDIR` directory and follow the steps.
  + Build environment: `tcloud build tuxiv.conf`
  + Submit job: `tcloud submit`
  + Monitor job: `tcloud show [job id]`
  + Cancel job: `tcloud cancel [job id]`

## MXNet
+ Dataset: mnist

+ Task: image classification

+ Code: [mnist.py](https://github.com/xcwanAndy/tcloud-sdk/blob/master/example/MXNET/mnist.py)

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
    environment:
        name: mxnet-env 
        dependencies:
            - mxnet=1.5.0
    job:
        name: test
        general:
            - nodes=2
    ~~~

+ Training process:

  + Enter the `TACC_WORKDIR` directory and follow the steps.
  + Build environment: `tcloud build tuxiv.conf`
  + Submit job: `tcloud submit`
  + Monitor job: `tcloud show [job id]`
  + Cancel job: `tcloud cancel [job id]`
