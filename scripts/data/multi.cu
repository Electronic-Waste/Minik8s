#include <iostream>

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <math.h>

#include "cuda_runtime.h"
#include "device_launch_parameters.h"
#include "cublas_v2.h"

#define M 2 // 
#define N 1 // 

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

__global__ void matrix_product(float a[M][N],float b[N][M],float c[M][M])
{
    int i = threadIdx.x + blockIdx.x * blockDim.x; 

    int j = threadIdx.y + blockIdx.y * blockDim.y; 

    if (i < M && j < M) 

    { 
        float tmp = 0;
        for (int m = 0;m < N;m++) {
            tmp += a[i][m]*b[m][j];
        }
        c[i][j] = tmp; 
    } 
}

void wrapper_product() 
{
    cudaError_t cudaStat;
    // 申请内存
    int Axy = M * N;
    int Bxy = N * M;
    int Cxy = M * M;
    float *h_A, *h_B, *h_C;
    h_A = (float*)malloc(Axy * sizeof(float));
    h_B = (float*)malloc(Bxy * sizeof(float));
    h_C = (float*)malloc(Cxy * sizeof(float));

    // 初始化数组
    initial_array(h_A, Axy);
    initial_array(h_B, Bxy);

    // 申请显存
    float (*d_A)[N];
    float (*d_B)[M];
    float (*d_C)[M];
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

    dim3 DimBlock(M, M);

    matrix_product <<<DimGrid, DimBlock>>>(d_A, d_B, d_C);
    cudaMemcpy(h_C,d_C,sizeof(float)*M*M,cudaMemcpyDeviceToHost);
    std::cout << "finish matrix producting" << std::endl;
    print_array(h_C,M,M); 
}

int main() {
    wrapper_product();
    return 0;
}