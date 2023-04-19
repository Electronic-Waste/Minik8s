# Minik8s 开题报告

### 人员组成

|  姓名  |     学号     |                         邮箱                          |
| :----: | :----------: | :---------------------------------------------------: |
| 魏靖霖 | 520021911429 |    [weijinglin@sjtu.edu.cn](mailto:weijinglin@sjtu.edu.cn)    |
| 王劭 | 520021911427 | [shaowang@sjtu.edu.cn](mailto:shaowang@sjtu.edu.cn) |
| 宋峻青 | 520021910991 |  [songjunqing@sjtu.edu.cn](mailto:songjunqing@sjtu.edu.cn)  |


### 自选功能选择

我们选择**MicroService**进行实现。

### 任务的时间安排

- 第一次迭代 2022/04/05 - 2022/04/30
  - 熟悉k8s的基本功能，设计方式
  - 实现Pod抽象，Service抽象，Pod ReplicaSet抽象和动态伸缩功能
  - 使用CI/CD进行部署
  - 准备中期答辩

- 第二次迭代 2022/05/01 - 2022/05/15

  - 完成DNS功能，容错功能
  - 使用交我算平台，支持GPU应用功能
  - 在多机上实现容器编排的功能

- 第三次迭代 2022/05/16 - 2022/05/31

  - 完成MicroService功能
  - 查漏补缺，准备最终答辩
  
### 人员分工

- 魏靖霖：实现Pod抽象，Pod ReplicaSet抽象（或者Deployment），容错功能，部署CI/CD测试，⽀持⾼级流量控制功能
- 王劭：实现CNI功能，动态伸缩功能，容错功能，GPU应用功能，对Pod流量进⾏劫持
- 宋峻青：实现Service抽象，DNS与转发功能，Serverless Workflow抽象，多机minik8s，⽀持⾃动化服务发现

### gitee仓库地址

https://gitee.com/jinglinwei/minik8s.git