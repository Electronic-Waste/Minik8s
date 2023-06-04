#include <iostream>

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <math.h>

#include "cuda_runtime.h"
#include "device_launch_parameters.h"
#include "cublas_v2.h"

#define M 32 // 
#define N 32 // 

#define BLOCK_SIZE 32  // 每个Block的线程数

void initial_array(float *array, int size)
{
    for(int i=0; i<size; i++)
    {
        array[i] = (float)(i);
    }
}

void print_array(float *array, int rows, int cols)
{
    for(int i=0; i<rows; i++)
    {
        for(int j=0; j<cols; j++)
        {
            std::cout << array[i*cols+j] << " ";
        }
        std::cout << std::endl;
    }
    std::cout << std::endl;
}

__global__ void matrix_add(float a[M][N],float b[M][N],float c[M][N])
{
    int i = threadIdx.x + blockIdx.x * blockDim.x; 

    int j = threadIdx.y + blockIdx.y * blockDim.y; 

    if (i < M && j < N) 

    { 
        c[i][j] = a[i][j] + b[i][j]; 
    } 
}

void wrapper_add() 
{
    cudaError_t cudaStat;
    // 申请内存
    int Axy = M * N;
    int Bxy = M * N;
    int Cxy = M * N;
    float *h_A, *h_B, *h_C;
    h_A = (float*)malloc(Axy * sizeof(float));
    h_B = (float*)malloc(Bxy * sizeof(float));
    h_C = (float*)malloc(Cxy * sizeof(float));

    // 初始化数组
    initial_array(h_A, Axy);
    initial_array(h_B, Bxy);

    // 申请显存
    float (*d_A)[N];
    float (*d_B)[N];
    float (*d_C)[N];
    cudaStat = cudaMalloc((void**)&d_A, Axy * sizeof(float));
    if (cudaStat != cudaSuccess) {
        printf ("device memory allocation failed\n");
        return;
    }
    cudaStat = cudaMalloc((void**)&d_B, Bxy * sizeof(float));
    if (cudaStat != cudaSuccess) {
        printf ("device memory allocation failed\n");
        return;
    }
    cudaStat = cudaMalloc((void**)&d_C, Cxy * sizeof(float));
    if (cudaStat != cudaSuccess) {
        printf ("device memory allocation failed\n");
        return;
    }
    cudaMemcpy(d_A, h_A, Axy * sizeof(float), cudaMemcpyHostToDevice);
    cudaMemcpy(d_B, h_B, Bxy * sizeof(float), cudaMemcpyHostToDevice);

    dim3 DimGrid(1, 1); 

    dim3 DimBlock(32, 32);

    matrix_add <<<DimGrid, DimBlock>>>(d_A, d_B, d_C);
    cudaMemcpy(h_C,d_C,sizeof(float)*M*N,cudaMemcpyDeviceToHost);
    std::cout << "finish matrix adding" << std::endl;
    print_array(h_C,M,N); 
}

int main() {
    wrapper_add();
    return 0;
}