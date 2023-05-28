#!/bin/bash

#SBATCH --job-name=gpujob
#SBATCH --partition=dgx2
#SBATCH -N 1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=6
#SBATCH --gres=gpu:1
#SBATCH --output=result.out
#SBATCH --error=result.out
module load gcc/8.3.0 cuda/10.2
nvcc test.cu -o test -lcublas
./test