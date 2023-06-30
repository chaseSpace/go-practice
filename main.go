package main

func main() {
	c := make(chan int)
	close(c)
	println(IsChannelClosed(c))
}

func IsChannelClosed(ch chan int) bool {
	select {
	case _, ok := <-ch:
		return !ok // 返回通道已关闭的状态
	default:
		return false // 通道仍然可用
	}
}
