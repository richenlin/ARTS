## 数组

数组是具有相同 **唯一类型** 的一组已编号且长度固定的数据项序列（这是一种同构的数据结构）

* 数组类型可以是任意的原始类型例如整型、字符串或者自定义类型
* 数组长度必须是一个常量表达式，并且必须是一个非负整数
* 数组长度也是数组类型的一部分，所以 [5] int 和 [10] int 是属于不同类型的
* 数组元素可以通过 索引（位置）来读取（或者修改），索引从 0 开始，第一个元素索引为 0
* 元素的数目，也称为长度或者数组大小必须是固定的并且在声明该数组时就给出（编译时需要知道数组长度以便分配内存）；数组长度最大为 2Gb
* 数组本身的类型是值类型（跟其他语言不同）

如果我们想让数组元素类型为任意类型的话可以使用空接口作为类型

```
var arr1 [5]int
// var arr1 = new([5]int)

for i := 0; i < len(arr1); i++ {
    arr1[i] = i * 2
}

for i := 0; i < len(arr1); i++ {
    fmt.Printf("Array at index %d is %d\n", i, arr1[i])
}

for i,_:= range arr1 {
    //...
}

arr2 := *arr1   //值拷贝
arr2[2] = 100   //不会影响到arr1

arr3 := &arr1   //指针
arr3[2] = 100   //arr1受影响
```

### 数组常量

如果数组值已经提前知道了，那么可以通过 数组常量 的方法来初始化数组，而不用依次使用 []= 方法（所有的组成元素都有相同的常量语法）。

```
//在这里数组长度可以写成 ... 或者直接忽略
var arrLazy = [...]int{5, 6, 7, 8, 22}

var arrAge = [5]int{18, 20, 15, 22, 16}

var arrKeyValue = [5]string{3: "Chris", 4: "Ron"}
```

### 多维数组

数组通常是一维的，但是可以用来组装成多维数组，例如：[3][5]int，[2][2][2]float64。

内部数组总是长度相同的。Go 语言的多维数组是矩形式的（唯一的例外是切片的数组）

```
const (
    WIDTH  = 1920
    HEIGHT = 1080
)

type pixel int
var screen [WIDTH][HEIGHT]pixel

for y := 0; y < HEIGHT; y++ {
    for x := 0; x < WIDTH; x++ {
        screen[x][y] = 0
    }
}

```

### 将数组传递给函数

把一个大数组传递给函数会消耗很多内存。有两种方法可以避免这种现象：

* 传递数组的指针
* 使用数组的切片

```
package main
import "fmt"

func main() {
    array := [3]float64{7.0, 8.5, 9.1}
    x := Sum(&array) // Note the explicit address-of operator
    // to pass a pointer to the array
    fmt.Printf("The sum of the array is: %f", x)
}

func Sum(a *[3]float64) (sum float64) {
    for _, v := range a { // derefencing *a to get back to the array is not necessary!
        sum += v
    }
    return
}
```

## 切片

切片（slice）是对数组一个连续片段的引用（该数组我们称之为相关数组，通常是匿名的）

* 切片是一个引用类型。
* 切片是可索引的，并且可以由 len() 函数获取长度
* 给定项的切片索引可能比相关数组的相同元素的索引小。和数组不同的是，切片的长度可以在运行时修改，最小为 0 最大为相关数组的长度：切片是一个 长度可变的数组
* 切片提供了计算容量的函数 cap() 可以测量切片最长可以达到多少。切片s: 0 <= len(s) <= cap(s)
* 多个切片如果表示同一个数组的片段，它们可以共享数据
* 因为切片是引用，所以它们不需要使用额外的内存并且比使用数组更有效率。Go中切片更常用

声明切片:
```
var identifier []type //不需要说明长度
```
切片的初始化格式:
```
var slice1 []type = arr1[start:end]
```

对于每一个切片（包括 string），以下状态总是成立的：
```
s == s[:i] + s[i:] // i是一个整数且: 0 <= i <= len(s)
len(s) <= cap(s)

slice1[0]  => arr1[start]

slice1 []type = arr1[:] // slice1 就等于完整的 arr1 数组，arr1[0:len(arr1)] 的一种缩写

arr1[2:] //和 arr1[2:len(arr1)] 相同，都包含了数组从第三个到最后的所有元素

arr1[:3] //和 arr1[0:3] 相同，包含了从第一个到第三个元素（不包括第三个）

```

### 创建切片

当相关数组还没有定义时，我们可以使用 make()或 new() 函数来创建一个切片 同时创建好相关数组。

```
slice1 := make([]type, len)
slice1 := new([]type)[start:end]

//下面两种方式等同
make([]int, 50, 100)
new([100]int)[0:50]

```

new () 和 make () 的区别:

* new (T) 为每个新的类型 T 分配一片内存，初始化为 0 并且返回类型为 * T 的内存地址：这种方法 返回一个指向类型为 T，值为 0 的地址的指针，它适用于值类型如数组和结构体；它相当于 &T{}。
* make(T) 返回一个类型为 T 的初始值，它只适用于 3 种内建的引用类型：切片、map 和 channe

## map集合

Map 是一种无序的键值对的集合。Map 最重要的一点是通过 key 来快速检索数据，key 类似于索引，指向数据的值。
Map 是一种集合，所以我们可以像迭代数组和切片那样迭代它。不过，Map 是无序的，我们无法决定它的返回顺序，这是因为 Map 是使用 hash 表来实现的。

* map是引用类型
* 在声明的时候不需要知道 map 的长度，map 是可以动态增长的
* 未初始化的 map 的值是 nil
* key 可以是 string、int、float。数组、切片和结构体不能作为 key (译者注：含有数组切片的结构体不能作为 key，只包含内建类型的 struct 是可以作为 key 的）但是指针和接口类型可以。
* 如果要用结构体作为 key 可以提供 Key() 和 Hash() 方法，这样可以通过结构体的域计算出唯一的数字或者字符串的 key
* value 可以是任意类型，通过使用空接口类型
* map 传递给函数的代价很小：在 32 位机器上占 4 个字节，64 位机器上占 8 个字节，无论实际上存储了多少数据。通过 key 在 map 中寻找值是很快的，但是比数组和切片要慢100倍，所以在性能敏感场景尽量使用切片
* map可以使用函数作为值，可以实现分支结构
* 常用的 len(map1) 方法可以获得 map 中的 pair 数目，这个数目是可以伸缩的，因为 map-pairs 在运行时可以动态添加和删除
  
```
var mapLit map[string]int
//var mapCreated map[string]float32
var mapAssigned map[string]int

mapLit = map[string]int{"one": 1, "two": 2}
mapCreated := make(map[string]float32)
mapAssigned = mapLit

mapCreated["key1"] = 4.5
mapCreated["key2"] = 3.14159
mapAssigned["two"] = 3

fmt.Printf("Map literal at \"one\" is: %d\n", mapLit["one"])
fmt.Printf("Map created at \"key2\" is: %f\n", mapCreated["key2"])
fmt.Printf("Map assigned at \"two\" is: %d\n", mapAssigned["two"])
fmt.Printf("Map literal at \"ten\" is: %d\n", mapLit["ten"])
```


## iota枚举

go里面有一个关键字iota， 用来声明enum的时候采用，默认开始值是0，const中每增减一行加1
```
package main
import(
    "fmt"
)

const(
    x = iota //x == 0
    y = iota //y == 1
    z = iota //z == 2
    w //声明省略了赋值，默认和之前一个值的字面相同。这里隐式的说w = iota，因此 w == 3
)

const v = iota //每遇到一个const关键字，iota就会重置，此时v == 0

const (
    h, i, j = iota, iota, iota // h == 0, i == 0, j == 0 iota在同一行内值相同
)

const (
	a       = iota //a=0
	b       = "B"
	c       = iota             //c=2
	d, e, f = iota, iota, iota //d=3,e=3,f=3
	g       = iota             //g = 4
)

//除非被显式设置为其它值或iota，每个const分组的第一个常量被默认设置为它的0值，第二及后续的常量被默认设置为它前面那个常量的值

//如果前面那个常量的值是iota，则它也被设置为iota

func main() {
	fmt.Println(a, b, c, d, e, f, g, h, i, j, x, y, z, w, v)
}
```