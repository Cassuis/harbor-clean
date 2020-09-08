# clean-harbor

半智能清理harbor中的镜像，主要用于CI中产生的镜像的自动清理。

通过对镜像的**tag前缀**的分类排序过滤来清理，如该镜像有dev-\*，test-\*等前缀的tag，则按照构建时间分别保留dev和test keepNum个镜像，其余的全部清理。



### Usage

```shell
Usage of clean-harbor:
  -h    help message
  -keepNum int
        每个repo保留的tag个数 (default 5)
  -password string
        harbor账号
  -projectName string
        projectName
  -url string
        harbor地址
  -user string
        harbor账号
```

### build

```shell
go build .
```

### crontab

```shell
for example
#> crontab -l
0 2 * * * /root/clean-harbor -url http://10.0.0.1 -user clean -password password -projectName cloud -keepNum 5 >> /var/log/cleanHarbor.log

# 需要在harbor中开启垃圾清理
```




