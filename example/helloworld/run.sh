#!/bin/bash
source /mnt/sharefs/home/user6/.Miniconda3/etc/profile.d/conda.sh
conda activate hello-be682e6a2d1fbd7357a70dc84651e9fd

python ${TACC_WORKDIR}/main.py \
