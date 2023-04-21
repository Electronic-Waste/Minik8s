# Source Reading

-   crictl的使用：
    -   `crictl pull image_name`：可以拉取一个镜像，但是仍然有具体的配置没有设置
-   kubeadm源码解读：
    -   `SubCmdRun()`方法简单的说就是要求一定后面要有子命令，否则打印用法并正常返回：如kubeadm config + subcmd
    -   `kubeadm config images pull --config xxx.yaml`指令学习：
        -   在最内层的`newCmdConfigImagesPull`函数内部解析了最后的pull子命令，通过`cmd.PersistentFlags()`方法返回flagset对象（表示这个flag可以被子命令共用），通过`AddImagesCommonConfigFlags`接口添加了很多flag，但是我们仅仅需要config flag，相当于通过yaml文件来初始化配置，通过`AddCRISocketFlag`添加对应的容器运行时的连接的url。
        -   通过一层层调用，在最后，我们需要调用最后的一个yaml解析器去解析yaml文件，在`staging/src/k8s.io/apimachinery/pkg/util/yaml`目录下，这里面会对于yaml文件的读取以及解析进行操作。
        -   yaml的文件解析是以`---`作为分隔符的。`YAMLReader`一次只会读出一个yaml document (以`---`分隔).