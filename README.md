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
#### TUXIV.CONF

You can use `tcloud init` to initialize the job configuration. `tcloud init` will create a template configuration file named `tuxiv.conf` . There are four parts in `tuxiv.conf`, config different parts of job submission. Noted that `tuxiv.conf` follows the yaml format.

+ Entrypoint

  In this section, you should insert you shell commands to run your code line-by-line. The tcloud CLI will run the job as your configurations.

  ~~~yaml
  entrypoint:
      - python ${TACC_WORKDIR}/mnist.py --epoch=3
  ~~~

+ Environment

  In this section, you can specify your conda configurations for virtual environment used in the cluster, including environment name, dependencies, source channels and so on.

  ~~~yaml
  environment:
      name: torch-env
      dependencies:
          - pytorch=1.6.0
          - torchvision=0.7.0
      channels: pytorch
  ~~~

+ Job

  In this section, you can specify your slurm configurations for slurm cluster resources, including number of nodes, CPUs, GPUs and so on. All the slurm cluster configuration should be set in the general part.

  ~~~yaml
  job:
      name: test
      general:
          - nodes=2
  ~~~

+ Datasets

  In this section, you can specify the data location or url for your job.

#### TACC VARIABLES

+ `TACC_WORKDIR`: TACC job workspace directory, each job has a different workspace directory.
+ `TACC_MASTER_IP`: Master node IP for initialize the communication group.
+ `TACC_WORLD_SIZE`: Size of communication group.
+ `TACC_RANK`: Rank of nodes for multi-node jobs.
+ `TACC_LOCAL_RANK`: Local rank of worker in same node.
+ `TACC_CLUSTER_LIST`: IP of each node in communication group.

## Example

Basic examples are provided under the [example](example) folder. These examples include: [TensorFlow](example/TensorFlow), [PyTorch](example/PyTorch) and [MXNet](example/MXNet).

