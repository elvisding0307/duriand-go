# README

## Docker

1.构建镜像：

```
docker build -t duriand .
```

2.运行容器：

```
docker run -p 7224:7224 --env-file .env duriand
```

3.自定义相关配置

```
docker run -d --name duriand --network mlan --ip 172.18.0.129 -p 7224:7224 -e MYSQL_HOST=172.18.0.2 -e DURIAND_SECRET_KEY=P@ssw0rd duriand:latest
```

