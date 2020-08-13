# defer函数

## 不执行的情况

defer 声明的代码在goroutine 没有执行：

```go
package main

import (
	"fmt"
	"time"
)
var ch chan  int
func main() {

	ch =make(chan  int)
	for i:=0;i<10 ;i++  {
		go test(i)
	}

	for i:=0;i<10 ;i++  {
		<-ch
	}
	
}
func test(index int ){
	time.Sleep(time.Second*1)
	defer fmt.Println("退出test ",index)
	ch<-index
}
```
运行结果：

```
结果一：退出test  5

结果二：什么都没有
```

不是说好的defer函数是总会执行的吗？结果并没有执行

分析：
> 当goroutine的主线程结束，goroutine 的子线程也会结束，剩下的无论什么的代码都不会执行。因为子线程已经结束。 

## 执行的顺序

```go
package main

import (
	"fmt"
)

func main() {
	defer_call()
}

func defer_call() {
	defer func() { fmt.Println("1") }()
	defer func() { fmt.Println("2") }()
	defer func() { fmt.Println("3") }()
}
```
执行结果：

```
3
2
1
```

分析：

> defer函数的执行顺序是后进先出（感觉是一个栈）。写代码的时候要注意defer函数的执行顺序