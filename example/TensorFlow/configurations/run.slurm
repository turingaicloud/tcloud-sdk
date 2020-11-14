#!/bin/bash
#SBATCH --nodes=2
export TACC_WORKDIR=/mnt/sharefs/home/testuser/WORKDIR/TensorFlow
export TACC_MODELDIR=/mnt/sharefs/home/testuser/models/WORKDIR/TensorFlow
export TACC_LOGDIR=/mnt/sharefs/home/testuser/logs/WORKDIR/TensorFlow
srun /mnt/sharefs/home/testuser/WORKDIR/TensorFlow/run.sh
