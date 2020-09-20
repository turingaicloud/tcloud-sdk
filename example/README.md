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

## MXNet
