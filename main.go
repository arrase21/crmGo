package main

import (
	"fmt"
	"slices"
)

func main() {
	list := []int{1, 2, 3, 4, 5}
	fmt.Println(maxnum(list))
	fmt.Println(min(list))
}
func maxnum(list []int) int {
	maxn := 0
	for _, i := range list {
		if i > maxn {
			maxn = i
		}
	}
	return maxn
}

func min(list []int) int {
	minNum := slices.Min(list)
	return minNum
}
