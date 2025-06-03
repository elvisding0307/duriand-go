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