package utils

import (
	"math/rand"
	"time"
)

// 辗转相除法求最大公约数
func LargestCommonDivisor(m, n int) int {
	defer PanicHandler()
	/*辗转相除法求最大公约数 */
	a := m
	b := n
	c := 0
	// /* 余数不为0，继续相除，直到余数为0 */
	for b != 0 {
		c = a % b
		a = b
		b = c
	}
	return a
}

// 产生一个不同的随机数
func Random() int64 {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
	// 返回参数的随机数
	return rand.Int63()
}
