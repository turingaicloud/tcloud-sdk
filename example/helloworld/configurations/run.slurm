#!/bin/bash
#SBATCH --output=/mnt/sharefs/home/htianab/USERDIR/slurm_log/hello.out
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=1
export TACC_WORKDIR=/mnt/sharefs/home/htianab/WORKDIR/helloworld
export TACC_USERDIR=/mnt/sharefs/home/htianab/USERDIR
srun /mnt/sharefs/home/htianab/WORKDIR/helloworld/run.sh
