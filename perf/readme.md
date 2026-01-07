## 性能分析

### 准备工作：引入pprof

引入pprof包：

```
_ "net/http/pprof"
```

启动std http server。

### 定位CPU消耗严重的 goroutine 或 runtime调用

程序运行十几秒后，然后执行（指定地址端口）：

```
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

等待几十秒，进入CLI，输入top：

```
(pprof) top
Showing nodes accounting for 29.33s, 99.76% of 29.40s total
Dropped 15 nodes (cum <= 0.15s)
      flat  flat%   sum%        cum   cum%
    29.33s 99.76% 99.76%     29.36s 99.86%  go-practice/perf.Test1.func1
```

这里输出的条目可能是函数路径，也可能是某个runtime调用。

> 注意：top指令是打印【函数“自身”占用的 CPU 时间】，不含子函数调用。而`top -cum`则是统计的是包含了子函数调用的累计耗时。

### 进一步定位 CPU 消耗最多的代码行

在pprof CLI：

```
(pprof) list func1
Total: 29.40s
ROUTINE ======================== go-practice/perf.Test1.func1 in C:\Users\yangr\Desktop\git\go-practice\perf\pprof_test.go
    29.33s     29.36s (flat, cum) 99.86% of Total
         .          .     18:   go func() {
         .          .     19:           var i int
    27.30s     27.33s     20:           for {
         .          .     21:                   i++
         .          .     22:                   i--
```

### 定位阻塞耗时最多的代码

代码中添加：`runtime.SetBlockProfileRate(1)`

执行：go tool pprof http://localhost:6060/debug/pprof/block?seconds=30

再执行top指令。

```
(pprof) top
Showing nodes accounting for 27.93us, 100% of 27.93us total
      flat  flat%   sum%        cum   cum%
   27.93us   100%   100%    27.93us   100%  sync.(*Cond).Wait
         0     0%   100%    27.93us   100%  net/http.(*Server).Serve.gowrap3
         0     0%   100%    27.93us   100%  net/http.(*conn).serve
         0     0%   100%    27.93us   100%  net/http.(*connReader).abortPendingRead
         0     0%   100%    27.93us   100%  net/http.(*response).finishRequest
```

block 采集的是【goroutine 在“同步原语”上阻塞的时间】，即：

- channel send / recv
- mutex / RWMutex
- Cond.Wait
- select

> - 也可以直接访问 http://localhost:6060/debug/pprof/block ，通过肉眼查找。
> - 注意：sleep不在block分析的范围内，只能通过trace分析。

### 定位内存消耗最多的函数

todo

### 定位开辟 goroutine 最多的函数

todo
