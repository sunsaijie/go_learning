# goroutine

简单来说就是go语言层面上支持并发，比如Python中的线程是在调用系统内核级的线程，而go语言在语言层面上实现了线程的调用，这样可以使用为更小的线程的开销，或者理解为python中的协程，但是支持并发。

### 并发模型

go语言推荐一种叫CSP并发模型的并发形式, 意思就是通信来共享内存。

> DO NOT COMMUNICATE BY SHARING MEMORY; INSTEAD, SHARE MEMORY BY COMMUNICATING.
> “不要以共享内存的方式来通信，相反，要通过通信来共享内存。

go语言中通过通道的方式进行每个goroutine中间的通信。

### 调用方式

+ 语法: `go 函数调用`  , go关键字后面跟上调用函数。
+ 使用示例

```go
import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup // 这里申明了一个goroutine的计数器

func runner(name int) {
	defer wg.Done()  // 这个goroutine执行完了，计数器自动减一
	fmt.Printf("这是第%d子goroutine的运行\n", name)
}

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)   //新建一个goroutine增加一个计数器值
		go runner(i)
	}
	wg.Wait()  // 阻塞，等待计数器归零
	fmt.Println("这是主goroutine的运行")
}

```



```
打印结果:
这是第9子goroutine的运行
这是第5子goroutine的运行
这是第2子goroutine的运行
这是第0子goroutine的运行
这是第7子goroutine的运行
这是第8子goroutine的运行
这是第6子goroutine的运行
这是第4子goroutine的运行
这是第3子goroutine的运行
这是第1子goroutine的运行
这是主goroutine的运行
```

> 这里因为需要等待子goroutine全部运行完，所以加了一个全局的sync.WaitGroup计数器，在主goroutine等待全部执行完。

# chan(通道)

### 通道简介

- 简介
个人理解跟Python的管道有点像，用于进程间通信，go语言的通道用于goruntine之间的数据通信。通道具有类型，只能放入符合通道类型的数据。
- 声明
```
  var 变量名 chan 通道类型
```
> 变量声明，如果没有初始化，没有开辟内存空间，从通道取值会阻塞。
- 初始化
使用`make`函数进行内存空间的申请，
```
  ch1 := make(chan int, 0)    //申请了int类型的通道，容量为0
  ch2 := make(chan int, 10) //申请了容量为10的通道
```
> 如果容量申请为0，没有相应得goroutine取值的话，程序会阻塞。

```
func main(){
  ch1 := make(chan int, 0)
  ch1 <- 10 //这里往里面存值，但是容量空间为0,并且没有其他的goroutine取值，就会产生死锁。 
}
```

### 通道的使用方法

通道只有一个符号两个动作，加上关闭方法，自己总结的。
- 一个符号`<-`，
- 两个动作就是存值和取值，符号在通道变量右边就是存值，在左边就是取值，应该也比较形象。
```
ch1:= make(chan int, 10)
ch1 <- 10 //存值
v := <- ch1 //取值
```
- 关闭 `close(ch1)`

### 通道的取值情况

| 通道 | nil, 只初始化，没有申请内存 | 非空                                                     | 空(申请了内存空间，但是没有值) | 满(超过最大容量)                                         | 没满                                                     |
| ---- | --------------------------- | -------------------------------------------------------- | ------------------------------ | -------------------------------------------------------- | -------------------------------------------------------- |
| 接收 | 阻塞                        | 正常取值                                                 | 阻塞                           | 正常取值                                                 | 正常取值                                                 |
| 发送 | 阻塞                        | 正常存值                                                 | 正常存值                       | 阻塞                                                     | 正常存值                                                 |
| 关闭 | panic(报错)                 | 关闭后，不能写入，但是可以读，读取完所有数据后，返回零值 | 关闭成功，再取数据返回零值     | 关闭后，不能写入，但是可以读，读取完所有数据后，返回零值 | 关闭后，不能写入，但是可以读，读取完所有数据后，返回零值 |

### 单向通道

主要用于通道参数的类型申明，限定函数内部对与通道的操作是取值还是存值

+ 只能取值: `<-chan 类型`
+ 只能存值:  `chan<- 类型`

```go
func onlySave(ch1 chan<- int){
  // a := <-ch1 // 这个编译会报错
  ch1 <- 10 // 对的
}

func onlyGet(ch1 <-chan int){
  a := <ch1 //对的
  // ch1 <- 10 这个编译会报错
}
```



### 代码示例

```go
import (
	"fmt"
	"math/rand"
)

// 随机生成数字，放入task通道,
// 新建一个20个goroutine的池，计算该数字 * 2，将结果存入result通道
// 主goroutine从result通道中取值后，打印结果

var task chan int64
var result chan int64

func calc(worker int, task <-chan int64, result chan<- int64) {
	for {
		v, ok := <-task
		if !ok {
			break
		}
		fmt.Printf("第%d个worker开始工作\n", worker)
		ret := v * 2
		result <- ret
		fmt.Printf("第%d个worker结束工作\n", worker)
	}
}

func main() {
	task = make(chan int64, 100)
	result = make(chan int64, 100)
	for t := 0; t < 5; t++ {
		go calc(t, task, result)
	}
	for i := 0; i < 100; i++ {
		a := rand.Int63()
		task <- int64(a)
	}
	close(task)

	for {
		a := <-result
		fmt.Printf("结果为%d\n", a)
	}

}


```



# select

go语言中的内置关键字，类似于switch关键字，有多个case和默认线路，每个case对应了一条的通道接收和发送过程，随机在非阻塞的通道中选择一个接受或发送的过程。



```go
import "fmt"

func main() {
	ch1 := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case ch1 <- i:        // 第二次存值的时候，因为容量为1，里面已经有了值，所有不会走这条
			fmt.Println("存入", i)
		case a := <-ch1:      // 第一次取值的时候，因为通道里面是空，所有不会走这条
			fmt.Printf("取值%d\n", a)
		}
	}
}
```



```
最终结果:
存入 0
取值0
存入 2
取值2
存入 4
取值4
存入 6
取值6
存入 8
取值8
```



# 锁

go语言中包括互斥锁，读写互斥锁

+ 互斥锁

  + 申明: ` var lock sync.Mutex`

  + 使用: `Lock()`和`Unlock`方法

  + 示例：

    ```go
    import (
    	"fmt"
    	"sync"
    )
    
    var wg sync.WaitGroup
    var lock sync.Mutex
    
    var x int = 0
    
    func add() {
    	defer wg.Done()
    	for i := 0; i < 50000; i++ {
    		lock.Lock()
    		x = x + 1
    		lock.Unlock()
    	}
    }
    
    func main() {
    	wg.Add(2)
    	go add()
    	go add()
    	wg.Wait()
    	fmt.Println(x)
    }
    ```

    

+ 读写互斥锁

  - 读操作不是独占，写操作是独占，比如读写数据库，如果是读锁，则针对所有的读操作，是没必要都加锁的，但是如果要写数据库，则必须加锁。使用条件是读操作高于写操作几个数量级以上的场景。

  - 变量申明: `var rwlock sync.RWMutex`

  - 使用方法:

    - 加读锁: `rwlock.RLock`和`rwlock.RUnlock`
    - 加写锁: `rwlock.Lock`和 `rwlock.Unlock`

  - 使用示例:

    ```go
    import (
    	"fmt"
    	"sync"
    	"time"
    )
    
    var x int = 0
    var rwlock sync.RWMutex
    var wg sync.WaitGroup
    
    func write() {
    	defer wg.Done()
    	rwlock.Lock()
    	x = x + 1
    	time.Sleep(time.Microsecond * 10)
    	rwlock.Unlock()
    }
    
    func read() {
    	defer wg.Done()
    	rwlock.RLock() //加读锁
    	a := x
    	fmt.Println(a)
    	time.Sleep(time.Microsecond * 5)
    	rwlock.RUnlock()
    
    }
    
    func main() {
    	for i := 0; i < 10; i++ {
    		wg.Add(1)
    		go write()
    	}
    
    	for i := 0; i < 100; i++ {
    		wg.Add(1)
    		go read()
    	}
    	wg.Wait()
    
    }
    ```

    