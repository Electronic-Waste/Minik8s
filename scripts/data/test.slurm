#!/bin/bash

#SBATCH --job-name=testoo        # 作业名
#SBATCH --partition=dgx2       # dgx2 队列
#SBATCH -N 1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=6     # 1:6 的 GPU:CPU 配比
#SBATCH --gres=gpu:1           # 1 块 GPU
#SBATCH --output=%j.out
#SBATCH --error=%j.err

module load gcc/8.3.0 cuda/10.2
nvcc test.cu -o test -lcublas
./test