package perf

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/trace"
	"testing"
	"time"
)

/*
Test

go tool pprof http://...pprof/profile?seconds=30 结果，即可以定位到具体函数
(pprof) top
flat  flat%   sum%        cum   cum%
29.37s 99.90% 99.90%     29.38s 99.93%  go-practice/perf.Test1.func1
*/
func TestHighCpuOccupied(t *testing.T) {
	go func() {
		var i int
		for {
			// i++ i--本身的执行耗时极低，不太可能被cpu采样
			i++
			i--
			// 加上sleep这一行，就不是高cpu消耗携程(cpu占用不会飙升)，仅仅是高调度型携程，sleep本质上是OS syscall，在tool pprof的top中会看到  runtime.stdcall 排名第一，而不是匿名函数本身
			// 详解：Sleep 之后，goroutine 不在 running 状态，而是进入睡眠（系统调用会让出cpu控制权），不会进 CPU 采样
			// - 不过sleep的耗时就会排到前面，以syscall的形式
			// - 切记，sleep是阻塞G的syscall，不会加剧cpu消耗，只会增加runtime.mcall调用次数，但这只是runtime 内部调度操作，不在OS层面。
			time.Sleep(time.Millisecond)
		}
	}()
	_ = http.ListenAndServe(":6060", nil)
}

/*
TestWithSleep：如何定位调用最频繁的方法（不一定是cpu密集、内存占用多），例如这里的sleep
- 说明：通过 http://.../debug/pprof 中列出的所有分析因素都无法观察到，只能通过trace

步骤：
1. 代码中添加trace逻辑，启动程序
2. 执行 go tool trace trace.out，会在浏览器打开trace页面，其中可以从goroutine、网络阻塞、同步阻塞、syscall、调度延迟方面去进行在线图形观察
3. 点击goroutines分析，可以看到一个根据总执行耗时排序的g列表，排除前面的trace调用代码，找到自定义的函数，点击进去进行分析
4. 在子页面中，可以单独分析此函数的各个方面，Breakdown小节中通过图形表格展示出了不同类型的执行耗时，可以直观的看到 block time(sleep) 是耗时最长的调用，长达十几秒。
*/
func TestWithSleep(t *testing.T) {
	// 打开 trace 文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 开始 trace
	if err = trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	go func() {
		var i int
		for {
			i++
			i--
			time.Sleep(time.Nanosecond)
		}
	}()
	_ = http.ListenAndServe(":6060", nil)
}
