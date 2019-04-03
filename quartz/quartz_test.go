package quartz

import (
	"fmt"
	"sort"
	"testing"
)

// 规则

// 6个字段 用空格间隔
// (second) (minute) (hour) (day of month, optional) (month) (day of week, optional)
// ?只能用于dom和dow，表示不关心，必须有且只有一个?

// 每个字段可以由多个表达式组成，表达式间用逗号间隔

// 支持的表达式
// "*"
// "?"
// "5"
// "30/6"
// "2-10" 闭区间
func TestQuartz(t *testing.T) {
	// a := time.Now()
	// fmt.Println(a)
	// b := time.Now().Local()
	// fmt.Println(b)
	a := []int{2, 3, 1, 0, 3, -1, 5, 7}
	sort.Sort(intPool(a))
	fmt.Println(a)
}

type intPool []int

func (ip intPool) Len() int {
	return len(ip)
}

func (ip intPool) Less(i, j int) bool {
	return ip[i] < ip[j]
}

func (ip intPool) Swap(i, j int) {
	temp := ip[i]
	ip[i] = ip[j]
	ip[j] = temp
}
