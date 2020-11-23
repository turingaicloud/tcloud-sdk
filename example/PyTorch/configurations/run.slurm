#!/bin/bash
#SBATCH --nodes=2
#SBATCH --ntasks=2
#SBATCH --cpus-per-task=1
#SBATCH --gres=gpu:1
export TACC_WORKDIR=/mnt/sharefs/home/testuser4/WORKDIR/PyTorch
export TACC_USERDIR=/mnt/sharefs/home/testuser4/USERDIR
srun /mnt/sharefs/home/testuser4/WORKDIR/PyTorch/run.sh
