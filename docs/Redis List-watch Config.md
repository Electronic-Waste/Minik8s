# Redis List-watch Config

## 1. 首先安装Redis

> 每台机子上都要安装

```shell
$ sudo apt-get install redis-server
```



## 2. 配置Redis

> 默认情况下Redis无密码，为了方便起见我们也不设密码

- 找到Redis的配置文件，地址为`/etc/redis/redis.conf`

- 将`bind 127.0.0.1`这一行注释掉

  > 这一行会让redis-server仅能通过`127.0.0.1`进行连接

- 将`protected-mode`由`yes`改为`no`

  > 若为`yes`则会禁止其它hosts上的`redis-cli`连接

- 找到redis-server进程，并且kill掉，redis-server会自动重启，这时就已经配置好了