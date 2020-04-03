# G语言实现多线程快排

## 1 单线程快排的实现

想要实现多线程的快速排序，首先要会写单线程的快排。快速排序的算法参照《算法导论》的快速排序这一章。

~~~go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	len := 50
	maxNum := 100
	var nums = make([]int, len)    
    	//大小为len的切片（注意不是数组，切片nums存放的是数组的起始地址，相当于一级指针）
	rand.Seed(time.Now().UnixNano())    //设置时间种子
	for i := 0; i < len; i++ {
        nums[i] = rand.Intn(maxNum)    //生成[0, maxNum)的随机数
	}
	fmt.Println(nums)
	quickSort(nums,0, len - 1)
	fmt.Println(nums)
}

func quickSort(nums[] int, p int, q int) {    
    //因为切片是一级指针，所以直接传参相当于拷贝指针，这个指针指向的还是数组的起始地址
	if p < q {
		r := partition(nums, p, q)
		quickSort(nums, p, r - 1)
		quickSort(nums, r + 1, q)

	}
}

func partition(nums[] int, p int, q int)int {
	i := p - 1
	x := nums[q]
	for j := p; j < q; j++ {
		if nums[j] < x {
			i++
			nums[i], nums[j] = nums[j], nums[i]
		}
	}
	i++
	nums[i], nums[q] = nums[q], nums[i]
	return i
}
~~~

这就是go语言的快速排序的实现，切片和数组的区别要搞清楚。

## 2 多线程快排是否可行？

在试图实现多线程快排是否可行之前，首先要验证一下多线程快排是否可行？快排的核心思想是分治算法，在单线程的快排中，我们是先排序比哨兵小的一半，再排序比哨兵大的另一边，而这两个过程是互不干扰的，因此这两者并行实现是完全可行的。

## 3 等待组实现多线程快排

~~~go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	len := 50
	maxNum := 100
	var nums = make([]int, len)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		nums[i] = rand.Intn(maxNum)
	}

	var wg sync.WaitGroup    //在这里的作用相当于信号量
	wg.Add(1)    //quickSort执行完main函数才能return，因此必须等待组加1，等quickSort执行完再减1

	fmt.Println(nums)

	quickSort(nums, 0, len - 1, &wg)
	wg.Wait()    //等到等待组为0才往下执行，否则阻塞
	fmt.Println(nums)
}

func quickSort(nums[] int, p int, q int, wg *sync.WaitGroup) {    //wg必须是引用类型
	defer wg.Done()    
    	//回调函数，quickSort函数return后再执行它，好处在于就算有很多个return，也只需要一条defer语句
    	//Done使得wg计数器减1

	if p < q {
		r := partition(nums, p, q)

		wg.Add(2)    
        	//要等下面两个quickSort执行完才return，因此等待组加2
        	//一定要在这里计数器加2，而不是每个quickSort开始的时候计数器加1，因为如果这样写，可能主函数会在quickSort
          //执行add之前到达wait，这样主函数就无法阻塞在wait处，就会提前退出
		go quickSort(nums, p, r - 1, wg)
		go quickSort(nums, r + 1, q, wg)
	}
}

func partition(nums[] int, p int, q int)int {
	i := p - 1
	x := nums[q]
	for j := p; j < q; j++ {
		if nums[j] < x {
			i++
			nums[i], nums[j] = nums[j], nums[i]
		}
	}
	i++
	nums[i], nums[q] = nums[q], nums[i]
	return i
}
~~~

搞清楚了等待组其实可以当作计数器之后，其实就不难实现了。之前我们说单线程快排是先执行比哨兵小的一边，再执行比哨兵大的一边，现在就两边一起执行。而**等待组的作用就是通过计数器记录排序子任务的个数，在每个排序子任务发生时计数器加1，每个排序子任务完成时计数器减1，到计数器为0的时候就说明所有任务都执行完了，主函数可以退出了。要注意的是计数器加减的时机，如注释所提到的那样。**值得注意的是，**计数器只保证阻塞在主函数的wait那里，在递归树中，可能出现上游函数已经返回而下游函数还在继续执行的情况。**

## 4 通道实现多线程快排

~~~go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	len := 50
	maxNum := 100
	var nums = make([]int, len)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		nums[i] = rand.Intn(maxNum)
	}

	ch := make(chan bool)

	fmt.Println(nums)

	go quickSort(nums,0, len - 1, &ch)    //一定要go，否则
	_ = <-ch

	fmt.Println(nums)
}

func quickSort(nums[] int, p int, q int, ch* chan bool) {
	/*defer func() {
		ch <- true
	}*/
	if p < q {
		r := partition(nums, p, q)
		lCh := make(chan bool)
		rCh := make(chan bool)
		go quickSort(nums, p, r - 1, &lCh)
		go quickSort(nums, r + 1, q, &rCh)
		_ = <-lCh    //等待第一个quickSort通过lCh传来数据，在此之前阻塞
		_ = <-rCh    //等待第二个quickSort通过rCh传来数据，在此之前阻塞
	}
	*ch <- true    
    	//通过通道ch传数据给上游，在上游接收数据之前阻塞，相当于一个信号，告诉上游，自己执行完毕
}

func partition(nums[] int, p int, q int)int {
	i := p - 1
	x := nums[q]
	for j := p; j < q; j++ {
		if nums[j] < x {
			i++
			nums[i], nums[j] = nums[j], nums[i]
		}
	}
	i++
	nums[i], nums[q] = nums[q], nums[i]
	return i
}
~~~

**通道在这里的作用是实现阻塞，防止递归树中的上游函数在下游函数执行完毕之前就返回。这一点和之前的等待组实现是不一样的**，当然在这个问题里，这两者都是可行的，但是别的问题就不一定了。为什么通道能做到这一点？因为通道有这样的特性，**一端发送数据到通道，并且发送端阻塞直到接收端接收数据；一端接受通道的数据，并且接收端阻塞直到发送端发送数据。**go语言有一个很经典的说法，**”不要通过共享内存实现通信， 而要通过通信实现共享内存“，所以可以把通道理解成一块大小为1的共享内存。**

更多精彩内容，敬请期待。

2020.4.3
