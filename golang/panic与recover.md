# panic与recover



> panic是用来停止当前函数的执行。当函数F调用panic时， F的正常执行立即停止，但是如果F定义了defer函数，那么defer会执行。

> recover内置函数允许程序管理一个正在paincing goroutine的行为

个人理解： panic函数就是抛出错误，如果没有被拦截，程序运行终止。为了防止程序终止运行，就得使用recover来接收panic错误或者说拦截panicing，并且recover函数可以将错误转化为error类型。因为panic错误不会影响defer函数运行，也就是说defer声明的函数即使遇到错误也会执行。

可以利用这个特性，在函数中用defer关键词声明一个函数，该函数内使用recover拦截panic错误，并返回给调用者error，从而使goroutine不挂掉。是不是有点像 try-catch了。

上码：

```go
package main


import (
    "fmt"
    "errors"
)

func main() {
    testError()
    afterErrorfunc()
}

func testError() {
    //defer catch()
    panic(" \"panic 错误\"")
    fmt.Println("抛出一个错误后继续执行代码")
}
func  catch()  {
    if r := recover(); r != nil {
        fmt.Println("testError() 遇到错误:", r)
        var err error
        switch x := r.(type) {
        case string:
            err = errors.New(x)
        case error:
            err = x
        default:
            err = errors.New("")
        }
        if err != nil {
            fmt.Println("recover后的错误:",err)
        }
    }
}

func afterErrorfunc(){
    fmt.Println("遇到错误之后 func ")
}
```
当panic 函数执行的时候导致后面函数 afterErrorfunc() 不能执行，main函数也抛出一个错误，整个程序异常退出。 

通过defer 关键字调用 catch函数，改造下：

```
func testError() {
    defer catch()
    panic(" \"panic 错误\"")
    fmt.Println("抛出一个错误后继续执行代码")
}
```

程序正常结束，没有因为panic(错误)而到导致程序终止挂掉。错误被recover 函数接收，转化为error类型的错误,最后输出“ recover后的错误:  “panic 错误” ” 而且后面 afterErrorfunc(）执行。
