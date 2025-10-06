package main

import (
	"fmt"
	"math/rand"
)

//func main() {
//	randNum := rand.Intn(1)
//
//	//arr := [][]int{{1, 2}, {1, 2}, {2, 3}}
//	//res := unique2D(arr)
//	//fmt.Println(res)
//	//a := nil
//	//i := map[int]*string{}
//	//i[1] = nil
//	//fmt.Println(i[1])
//	//fmt.Println(a)
//	//ch := make(chan int)
//	//select {
//	//case i := <-ch:
//	//	fmt.Println(i)
//	//default:
//	//	fmt.Println("default")
//	//}
//	//h := func() func() int {
//	//	i := func() int {
//	//		fmt.Println("h")
//	//		return 0
//	//	}
//	//	i()
//	//	return i
//	//}
//	//res := h()
//	//fmt.Println(res)
//}

func main() {
	res := 0
	count := 0
	for res < 3 {
		if res == 0 {
			count++
			res++
			fmt.Println("当前等级：0 → 消耗1元宝，结果：升级成功！当前等级：1")
			continue
		}
		randNum := rand.Float64()
		if res == 1 {
			if float64(randNum) < 1.0/3.0 {
				count++
				res--
				fmt.Println("当前等级：1 → 消耗1元宝，结果：升级失败！当前等级：0")
			} else if float64(randNum) < 2.0/3.0 {
				count++
				fmt.Println("当前等级：1 → 消耗1元宝，结果：停留原级当前等级：1")
			} else {
				count++
				res++
				fmt.Println("当前等级：1 → 消耗1元宝，结果：升级成功！当前等级：2")
			}
			continue
		}
		if res == 2 {
			if float64(randNum) < 2.0/3.0 {
				count++
				fmt.Println("当前等级：2 → 消耗1元宝，结果：停留原级当前等级：2")
			} else {
				count++
				res++
				fmt.Println("当前等级：2 → 消耗1元宝，结果：升级成功！当前等级：3")
			}
		}
	}
	fmt.Printf("英雄升到3级，总消耗元宝：%d", count)
}

func unique2D(arr [][]int) [][]int {
	seen := make(map[string]bool)
	result := [][]int{}
	for _, sub := range arr {
		// 将子切片转换为字符串表示
		str := fmt.Sprintf("%v", sub)
		if !seen[str] {
			seen[str] = true
			result = append(result, sub)
		}
	}
	return result
}
