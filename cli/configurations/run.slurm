#!/bin/bash
#SBATCH --node=2
#SBATCH --ntasks-per-node=24
#SBATCH --cpus-per-task=2
#SBATCH --use_gpu=1
module load cuda=9.0
export OMP_NUM_THREADS=1
