package main

import "fmt"

func main() {
	//DeferClosure()
	//DeferClosureV1()
	//println(DeferReturn())
	//println(DeferReturnV1())
	//println(DeferReturnV2())
	//DeferClosureLoopV1()
	//DeferClosureLoopV2()
	//DeferClosureLoopV3()
	ShareSlice()
}

func DeferClosure() {
	i := 0
	defer func() {
		println(i)
	}()
	i = 1
}

func DeferClosureV1() {
	i := 0
	defer func(i int) {
		println(i)
	}(i)
	i = 1
}

// DeferReturn 不带名字的返回值，修改不了
func DeferReturn() int {
	i := 0
	defer func() {
		i = 1
	}()
	return i
}

// DeferReturnV1 带名字的返回值可以修改
func DeferReturnV1() (i int) {
	i = 0
	defer func() {
		i = 1
	}()
	return i
}

func DeferReturnV2() *MyStruct {
	a := &MyStruct{
		name: "Jerry",
	}
	defer func() {
		a.name = "Tom"
	}()
	return a
}

type MyStruct struct {
	name string
}

func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}

func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}

func DeferClosureLoopV3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}

func ShareSlice() {
	s1 := []int{1, 2, 3, 4}
	s2 := s1[2:]
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s2[0] = 99
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s2 = append(s2, 199)
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s2[1] = 1999
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))
}
