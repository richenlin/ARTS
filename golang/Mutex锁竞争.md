# Mutex锁竞争

 使用协程获取golang 对象锁，想当然的认为是按先按时间先后顺序获取的，写代码验证一下：

```go
package main

import (
	"sync"
	"fmt"
	"time"
)
//开启10个线程  同时去竞争一个互斥锁 谁有能力谁上

var mutex *sync.Mutex;
var  ch chan int


func main() {
	mutex=new(sync.Mutex)
	fmt.Println("start")
	ch=make(chan int )
	for i:=0; i<10; i++ {
		go TestMutex(ch,i)
	}
	for i:=0; i<10; i++ {
	<-ch
	}
}

func TestMutex(ch chan int,index int)  {
	fmt.Println("to enter mutex","index=",index)
	mutex.Lock();
	defer mutex.Unlock()
	defer fmt.Println("unLock","index=",index)
	fmt.Println("in mutex","index=",index)
	time.Sleep(2*time.Second)
	ch<-1
```

运行结果如下：

```shell
to enter mutex index= 3
to enter mutex index= 9
in mutex index= 3
to enter mutex index= 5
to enter mutex index= 6
to enter mutex index= 7
to enter mutex index= 8
to enter mutex index= 1
to enter mutex index= 0
to enter mutex index= 4
to enter mutex index= 2
unLock index= 3
in mutex index= 9
unLock index= 9
in mutex index= 5
unLock index= 5
in mutex index= 6
unLock index= 6
in mutex index= 7
unLock index= 7
in mutex index= 8
unLock index= 8
in mutex index= 1
unLock index= 1
in mutex index= 0
unLock index= 0
in mutex index= 4
in mutex index= 2
unLock index= 4
unLock index= 2
```

通过运行结果发现所有协程都开始执行，但是进入锁的协程编号变成了3，不是按12345678的顺序获取互斥锁

总结:

> ```javascript
> 协程获取锁先后顺序不是按时间来获取，而是竞争关系谁有能力谁上
> ```



## 读写锁（RWmutex）

### 模型一： 多个协程一起读

```go
var RWmutex *sync.RWMutex;
var  RWch chan int

func main() {
   RWmutex=new(sync.RWMutex)
   fmt.Println("start")
   RWch=make(chan int )
   for i:=0; i<10; i++ {
      go TestRWMutex(i)
   }
   for i:=0; i<10; i++ {
      <-RWch
   }
}

func TestRWMutex(index int)  {
   fmt.Println("进入读写锁，准备读","index=",index)
   RWmutex.RLock();
   //读取数据
   fmt.Println("读..","index=",index)
   time.Sleep(4*time.Second)
   RWmutex.RUnlock()
   fmt.Println("离开读写锁..","index=",index)
   RWch<-1
}
```

运行结果

```
start
进入读写锁，准备读 index= 0
进入读写锁，准备读 index= 2
读.. index= 2
进入读写锁，准备读 index= 4
读.. index= 4
进入读写锁，准备读 index= 9
读.. index= 9
进入读写锁，准备读 index= 3
读.. index= 3
进入读写锁，准备读 index= 8
读.. index= 8
读.. index= 0
进入读写锁，准备读 index= 1
读.. index= 1
进入读写锁，准备读 index= 7
读.. index= 7
进入读写锁，准备读 index= 6
读.. index= 6
进入读写锁，准备读 index= 5
读.. index= 5
离开读写锁.. index= 4
离开读写锁.. index= 3
离开读写锁.. index= 9
离开读写锁.. index= 2
离开读写锁.. index= 7
离开读写锁.. index= 1
离开读写锁.. index= 5
离开读写锁.. index= 0
离开读写锁.. index= 6
离开读写锁.. index= 8
```

分析：

> 读锁并不互斥，可以连续加锁



### 模型二：写保护

```go
package main

import (
	"fmt"
	"time"
	"sync"
)
var RWWmutex *sync.RWMutex;
var  RWWch chan int

func main() {
	RWWmutex=new(sync.RWMutex)
	fmt.Println("start")
	RWWch=make(chan int )
	for i:=0; i<10; i++ {
		go TestRWWMutex(RWWch,i)
	}
	for i:=0; i<10; i++ {
		<-RWWch
	}
}

func TestRWWMutex(ch chan int,index int) {
	fmt.Println("进入写锁","index=",index)
	RWWmutex.Lock();

	//写东西..
	fmt.Println("写东西.....","index=",index)
	time.Sleep(1*time.Second)
	RWWmutex.Unlock()
	fmt.Println("离开写锁","index=",index)
	RWWch<-1
}
```

运行结果:

```
start
进入写锁 index= 0
进入写锁 index= 1
进入写锁 index= 6
进入写锁 index= 7
进入写锁 index= 2
进入写锁 index= 3
进入写锁 index= 9
写东西..... index= 0
进入写锁 index= 4
进入写锁 index= 5
进入写锁 index= 8
离开写锁 index= 0
写东西..... index= 1
写东西..... index= 6
离开写锁 index= 1
写东西..... index= 7
离开写锁 index= 6
离开写锁 index= 7
写东西..... index= 2
离开写锁 index= 2
写东西..... index= 3
离开写锁 index= 3
写东西..... index= 9
离开写锁 index= 9
写东西..... index= 4
离开写锁 index= 4
写东西..... index= 5
离开写锁 index= 5
写东西..... index= 8
离开写锁 index= 8
```

分析：

> 写锁是互斥的，写的时候其他协程等待

### 模型三：混合读写

```go
package main

import (
	"sync"
	"fmt"
	"time"
)

//读写同时发生的模型
//读的时候的写 NO
//写的时候读 NO


var WR_Wmutex *sync.RWMutex;
//var WR_Rmutex *sync.RWMutex;
var  WRch chan int

func main() {
	WR_Wmutex=new(sync.RWMutex)
	fmt.Println("start")
	WRch=make(chan  int )

	for i:=0;i<5;i++  {
		go TestRW_WMutex(i)
		go TestRW_RMutex(i)
	}
	for i:=0;i<10 ; i++ {
		<-WRch
	}
}

func TestRW_WMutex(index int) {
	WR_Wmutex.Lock();
	//写东西..
	fmt.Println("WWW 正在写.....","index=",index)
	time.Sleep(1*time.Second)
	fmt.Println("WWW 离开写锁","index=",index)
	WR_Wmutex.Unlock()

	WRch<-1
}

func TestRW_RMutex(index int) {
	WR_Wmutex.RLock();
	//读东西..
	fmt.Println("RRR 正在读.....","index=",index)
	time.Sleep(1*time.Second)
	fmt.Println("RRR 离开读锁","index=",index)
	WR_Wmutex.RUnlock()

	WRch<-1
}
```

运行结果

```
start
WWW 正在写..... index= 0
WWW 离开写锁 index= 0
RRR 正在读..... index= 4
RRR 正在读..... index= 2
RRR 正在读..... index= 3
RRR 正在读..... index= 0
RRR 正在读..... index= 1
RRR 离开读锁 index= 4
RRR 离开读锁 index= 2
RRR 离开读锁 index= 1
RRR 离开读锁 index= 3
RRR 离开读锁 index= 0
WWW 正在写..... index= 1
WWW 离开写锁 index= 1
WWW 正在写..... index= 2
WWW 离开写锁 index= 2
WWW 正在写..... index= 3
WWW 离开写锁 index= 3
WWW 正在写..... index= 4
WWW 离开写锁 index= 4
```

分析:

> 协程可以并行读，但与写锁互斥，也就是写的时候不能读

### 模型四：先读后写

```go
package main

import (
	"sync"
	"fmt"
	"time"
)

//读写同时发生的模型
//读的时候的写 NO
//写的时候读 NO


var WR_Wmutex *sync.RWMutex;
//var WR_Rmutex *sync.RWMutex;
var  WRch chan int

func main() {
	WR_Wmutex=new(sync.RWMutex)
	fmt.Println("start")
	WRch=make(chan  int )

	for i:=0;i<5;i++  {
		go TestRW_RMutex(i)
		go TestRW_WMutex(i)
	}
	for i:=0;i<10 ; i++ {
		<-WRch
	}
}

func TestRW_WMutex(index int) {
	WR_Wmutex.Lock();
	//写东西..
	fmt.Println("WWW 正在写.....","index=",index)
	time.Sleep(1*time.Second)
	fmt.Println("WWW 离开写锁","index=",index)
	WR_Wmutex.Unlock()

	WRch<-1
}

func TestRW_RMutex(index int) {
	WR_Wmutex.RLock();
	//读东西..
	fmt.Println("RRR 正在读.....","index=",index)
	time.Sleep(1*time.Second)
	fmt.Println("RRR 离开读锁","index=",index)
	WR_Wmutex.RUnlock()

	WRch<-1
}
```

运行结果：

```
start
RRR 正在读..... index= 0
RRR 离开读锁 index= 0
WWW 正在写..... index= 4
WWW 离开写锁 index= 4
RRR 正在读..... index= 4
RRR 正在读..... index= 1
RRR 正在读..... index= 2
RRR 正在读..... index= 3
RRR 离开读锁 index= 3
RRR 离开读锁 index= 2
RRR 离开读锁 index= 4
RRR 离开读锁 index= 1
WWW 正在写..... index= 0
WWW 离开写锁 index= 0
WWW 正在写..... index= 1
WWW 离开写锁 index= 1
WWW 正在写..... index= 2
WWW 离开写锁 index= 2
WWW 正在写..... index= 3
WWW 离开写锁 index= 3
```

分析:

> 读的时候不能写，可以多次读；写的时候不能读，其他人也不能写

