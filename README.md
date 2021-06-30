# TCLOUD-SDK
## Command-line Interface used for TACC job submission.
```
TCLOUD Command-line Interface v0.2.6

Usage:
tcloud [command] [flags] [args]

Available Commands:
    tcloud init
    tcloud config [-u/-f] [args]
    tcloud upload [-c] <local_dirpath> [<remote_dirpath>]
    tcloud download [<filepath>]
    tcloud add [<dependency_name>]
    tcloud submit [<path_to_repo>]
    tcloud ps [-j] [<JOB_ID>]
    tcloud install [<path_to_repo>]
    tcloud log [-j] [<JOB_ID>]
    tcloud cancel [-j] [<JOB_ID>]
    tcloud ls [<dirpath>]

Use "tcloud [command] --help" for more information about commands.
```

## Installation
You can try out the latest features by directly install from master branch:

```
git clone https://github.com/Turing-AI-Cloud/tcloud-sdk.git
cd tcloud-sdk
echo 'export GOPATH=$PWD' >> ~/.bash_profile
make install
```

## Configuration
### CLI Configuration
1. Before using the tcloud CLI and submit ML jobs to TACC, you need to configure your TACC credentials. You can do this by running the `tcloud config` command:
```
$ tcloud config [-u/--username] MYUSERNAME
$ tcloud config [-f/--file] MYPRIVATEFILEPATH
```
2. You need to run `tcloud init` command to obtain the latest cluster hardware information from TACC cluster.

### Job Configuration
#### TUXIV.CONF

You can use `tcloud init` to pull the latest cluster configurations from remote. There are four parts in `tuxiv.conf`, config different parts of job submission. Noted that `tuxiv.conf` follows the yaml format.

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

  In this section, you can specify your slurm configurations for slurm cluster resources, including number of nodes, CPUs, GPUs, output file and so on. All the slurm cluster configuration should be set in the general part.

  ~~~yaml
  job:
      name: test
      general:
          - nodes=2
          - output=${TACC_SLURM_USERLOG}/output.log
  ~~~

+ Datasets

  In this section, you can specify the needed CityNet dataset name, and tcloud will help place the dataset access in `TACC_USERDIR`. You can view the table of CityNet datasets at [CityNet Dataset Info](https://docs.google.com/spreadsheets/d/18qi2YpYvuXkWns7KY9pHYQclhS1Yyt5ysqgZ4plYcTg/edit#gid=0).

  ~~~yaml
  datasets:
    - OpenRoadMap
  ~~~

#### TACC VARIABLES

+ `TACC_WORKDIR`: TACC job workspace directory, each job has a different workspace directory.
+ `TACC_USERDIR`: TACC User directory.
+ `TACC_SLURM_USERLOG`: Slurm log directory, the default value is `${TACC_USERDIR}/slurm_log`.

## Example

Basic examples are provided under the [example](example) folder. These examples include: [HelloWorld](example/helloworld), [TensorFlow](example/TensorFlow), [PyTorch](example/PyTorch) and [MXNet](example/MXNet).

