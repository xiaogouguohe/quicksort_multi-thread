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
        	//一定要在这里计数器加2，而不是每个quickSort开始的时候计数器加1，因为如果这样写，可能主函数会在quickSort执行add之前到达wait，这样主函数就无法阻塞在wait处，就会提前退出
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
