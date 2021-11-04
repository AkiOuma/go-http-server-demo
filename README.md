# Http Server Demo

## 需求
基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

## 实现思路
1. 使用context.WithCancel创建一个context对象和对应的cancel方法
2. 在启动http服务的函数种传入一个context对象，在服务启动的函数(ServerRun)中创建一个go routine，在其中并使用select持续监听context的状态是否为Done。若context为Done时，则强行让http.Server对象执行Shutdown方法
3. 在启动监听Linux Signal的函数中传入一个context对象，并在监听函数(ReceiveSignal)中创建一个go routine，在其中并使用select持续监听context的状态是否为Done，与监听Linux Signal是否为Interrupt，若满足上面两个条件，则返回一个error
4. 在main函数中创建一个errgroup，使用errgroup.Go方法调用服务启动的函数(ServerRun)和监听函数(ReceiveSignal)。如果其中某个调用返回了错误，则直接执行cancel方法，将所有的go routine终止


## 项目结构
```
.
├── cmd
│   └── main.go 
├── go.mod
├── go.sum
├── README.md
├── server
│   ├── controller
│   │   ├── home.go   
│   │   └── stopper.go
│   └── server.go
└── signal-processor
    └── processor.go
```
* `server/controller/home.go`
  模拟一个主页的数据接口
* `server/controller/stopper.go`
  模拟一个可以终止http服务的接口
* `server/server.go`
  http服务启动
* `signal-processor/processor.go`
  Linux Signal监听函数

## 实现效果

### 任务启动
`go run ./cmd/ .`

### 获取数据
```bash
$ curl http://localhost:8080
{"message":"welcome"}
```
可以正常获取到home的json数据

### 停止http服务
```bash
$ curl http://localhost:8080/stop
```

任务进程打印以下内容：
```bash
http server    : 2021/11/04 23:43:28 HTTP server Shutdown by api
signal receiver: 2021/11/04 23:43:28 stop receiving signal: 
        context canceled
main           : 2021/11/04 23:43:28 Exit Reason: 
        http: Server closed
```
可以看到，http服务端在被终止后，Linux Signal接受服务的go routine也停止接受信号退出了

### 接受Interrupt信号
```bash
$ go run ./cmd/
^C
signal receiver: 2021/11/04 23:46:15 receive signal: 
        interrupt
http server    : 2021/11/04 23:46:15 HTTP server Shutdown: 
        context canceled
main           : 2021/11/04 23:46:15 Exit Reason: 
        interrupt
```

可以看到在接受了^C(Interrupt)信号后，Linux Signal接受服务停止后，http服务也终止了