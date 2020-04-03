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
