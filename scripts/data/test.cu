#include <iostream>

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <math.h>

#include "cuda_runtime.h"
#include "device_launch_parameters.h"
#include "cublas_v2.h"

#define M 8 // 矩阵行
#define K 8 // 矩阵列、矩阵行
#define N 8 // 矩阵列

#define BLOCK_SIZE 32  // 每个Block的线程数

// 初始化数组
void initial_array(float *array, int size)
{
    for(int i=0; i<size; i++)
    {
        array[i] = (float)(rand()%10+1);
    }
}

// 打印数组
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

void matrix_multiplication_cublas(int dimx_t, int dimy_t)
{
    cudaError_t cudaStat;
    // 申请内存
    int Axy = M * K;
    int Bxy = K * N;
    int Cxy = M * N;
    float *h_A, *h_B, *h_C;
    h_A = (float*)malloc(Axy * sizeof(float));
    h_B = (float*)malloc(Bxy * sizeof(float));
    h_C = (float*)malloc(Cxy * sizeof(float));

    // 初始化数组
    initial_array(h_A, Axy);
    initial_array(h_B, Bxy);

    // 申请显存
    float *d_A, *d_B, *d_C;
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

    // 设置参数
    int dimx = dimx_t;
    int dimy = dimy_t;
    dim3 block(dimx, dimy);
    dim3 grid((M+block.x-1)/block.x, (N+block.y-1)/block.y);

    // 设置参数
    cudaEvent_t gpustart, gpustop;
    float elapsedTime = 0.0;

    // 创建句柄
    cublasHandle_t handle;
    cublasCreate(&handle);
    elapsedTime = 0.0;
    cudaEventCreate(&gpustart);
    cudaEventCreate(&gpustop);
    cudaEventRecord(gpustart, 0);

    // 二维矩阵乘法-CUBLAS计算
    float a = 1, b = 0;
    cublasSgemm(
            handle,
            CUBLAS_OP_T,   // 矩阵A的属性参数，转置，按行优先
            CUBLAS_OP_T,   // 矩阵B的属性参数，转置，按行优先
            M,             // 矩阵A行数、矩阵C行数
            N,             // 矩阵B列数、矩阵C列数
            K,             // 矩阵A列数、矩阵B行数
            &a,            // alpha的值
            d_A,           // 左矩阵，为A
            K,             // A的leading dimension，此时选择转置，按行优先，则leading dimension为A的列数
            d_B,           // 右矩阵，为B
            N,             // B的leading dimension，此时选择转置，按行优先，则leading dimension为B的列数
            &b,            // beta的值
            d_C,           // 结果矩阵C
            M              // C的leading dimension，C矩阵一定按列优先，则leading dimension为C的行数
    );
    cudaMemcpy(h_C, d_C, Cxy * sizeof(float), cudaMemcpyDeviceToHost); // 显存拷贝到内存
    cudaDeviceSynchronize();
    cudaEventRecord(gpustop, 0); // 记录结束时间
    cudaEventSynchronize(gpustop);
    cudaEventElapsedTime(&elapsedTime, gpustart, gpustop); // 计算耗时
    cudaEventDestroy(gpustart);
    cudaEventDestroy(gpustop);

    // 打印计算结果
    std::cout << "Matrix_A: " << M << "x" << K << std::endl;
    print_array(h_A, M, K);
    std::cout << "Matrix_B: " << K << "x" << N << std::endl;
    print_array(h_B, K, N);
    std::cout << "Matrix_C: " << M << "x" << N << std::endl;
    print_array(h_C, M, N);

    // 打印耗时
    printf("matrix_multiplication_cublas: ");
    printf("gridx: %4d, gridy: %4d, blockx: %4d, blocky: %4d", grid.x, grid.y, block.x, block.y);
    printf(", runtime: %8fs\n", elapsedTime/1000);

    // 释放显存
    cudaFree(d_A);
    cudaFree(d_B);
    cudaFree(d_C);
    // 释放内存
    free(h_A);
    free(h_B);
    free(h_C);
    // 释放设备
    cudaDeviceReset();
}

int main()
{
    matrix_multiplication_cublas(2, 2);
    return 0;
}