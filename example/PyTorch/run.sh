#!/bin/bash
source /mnt/sharefs/home/htianab/.Miniconda3/etc/profile.d/conda.sh
conda activate torch-env

. ${TACC_WORKDIR}/get_nic_name.sh && python ${TACC_WORKDIR}/mnist.py --epoch=100 \
