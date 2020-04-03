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
