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
├── Makefile
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
```bash
$ make serve
go run ./cmd/ .
main           : 2021/11/05 21:38:48 starting http server at 8080
main           : 2021/11/05 21:38:48 starting http server at 8082
main           : 2021/11/05 21:38:48 starting http server at 8081
```

### 获取数据
```bash
$ curl http://localhost:8080
{"message":"welcome"}
$ curl http://localhost:8081
{"message":"welcome"}
$ curl http://localhost:8082
{"message":"welcome"}
```
可以正常获取到来自三个不同端口启动的http服务的home的json数据

### 停止http服务
```bash
$ curl http://localhost:8080/stop
```

任务进程打印以下内容：
```bash
http server    : 2021/11/05 21:41:26 HTTP server(8080) Shutdown by api
http server    : 2021/11/05 21:41:26 HTTP server(8080) Shutdown: 
        context canceled
signal receiver: 2021/11/05 21:41:26 stop receiving signal: 
        context canceled
http server    : 2021/11/05 21:41:26 HTTP server(8082) Shutdown: 
        context canceled
http server    : 2021/11/05 21:41:26 HTTP server(8081) Shutdown: 
        context canceled
main           : 2021/11/05 21:41:26 Exit Reason: 
        http: Server closed
```
可以看到，启动自8080端口的http服务端在被终止后，其余两个http服务以及Linux Signal接收服务的go routine也停止接受信号退出了。若终止请求来自8081或者8082端口的服务也会获得类似的结果

### 接受Interrupt信号
```bash
$ go run ./cmd/
^C
signal receiver: 2021/11/05 21:38:51 receive signal: 
        interrupt
http server    : 2021/11/05 21:38:51 HTTP server(8081) Shutdown: 
        context canceled
http server    : 2021/11/05 21:38:51 HTTP server(8080) Shutdown: 
        context canceled
http server    : 2021/11/05 21:38:51 HTTP server(8082) Shutdown: 
        context canceled
main           : 2021/11/05 21:38:51 Exit Reason: 
        interrupt
make: *** [serve] Error 1
```

可以看到在接受了^C(Interrupt)信号后，Linux Signal接受服务停止后，三个启动自不同端口的http服务也终止了