# write-file
## 编译启动
```shell
go build -o writer
./writer -v 3
```
## 调整日志级别
```shell
# 动态调整日志级别，减少日志输出
kill -SIGUSR1 $PID
# 动态调整日志级别，打印更详细信息
kill -SIGUSR2 $PID
```