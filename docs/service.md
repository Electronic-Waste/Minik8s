# Service

## 1. Service目前支持的功能

- 通过`./bin/kubectl apply <service-yaml-flie>`来应用service规则
- 通过`./bin/kubectl delete service <service-name>`来删除service规则



## 2. Service尚未支持的功能

- 跨多机访问podIP
- 通过`./bin/kubectl get service <service-name>`来获取service相关的信息



### 3. yaml文件示例

- 一个service.yaml文件示例，通过指定selector的app来匹配Pod

  > 比如selector.app是`test-pod`，那么它就会去匹配metadata.labels.app是`test-pod`的Pod

```yaml
kind: Service
metadata:
  name: test-service
spec:
  ports:
  - name: go1-port
    protocol: TCP
    port: 8080
    targetPort: 8080
  - name: go2-port
    protocol: TCP
    port: 80
    targetPort: 80
  selector:
    app: test-pod
```

- 下面给出的的Pod yaml文件是上面的Service所匹配的

```yaml
kind: Pod
metadata:
  name: test
  labels:
    app: test-pod
...
```



## 4. 运行结果示例

![image-20230524191508529](./img/service-result.png)