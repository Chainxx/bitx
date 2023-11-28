package lotto

import (
	"github.com/chainxx/bitx/pkg/math"
)

// 博彩

/*
双色球的投注方式是6+1，也就是从33个红球中选取6个红球号码，再从16个蓝球中选择1个蓝球号码。这样的一组号码构成一个投注单元，即为1柱。

当你说的6+2，我猜测你可能指的是在常规的6+1基础上再多选择一个蓝色号码，那么这会产生2柱。因为你有两种不同的蓝色号码可以和6个红色号码配对。

对于7+1、7+2、8+2等更多数目的选项来说，总的计算方式是“组合数”，或者说是从n个不同元素中取出m个元素（不考虑顺序）有多少种方法。
公式为C(n, m) = n! / [m!(n-m)!]。

所以：

- 7+1：从7个红色号码中选取任意6个搭配1个蓝色号码，有C(7, 6) = 7种可能；
- 7+2：从7个红色号码中选取任意6个搭配任意1个蓝色号码，然后乘上蓝色号码可以有两种选择，即C(7, 6)*2 = 14种可能；
- 8+2：首先从8个红球中选择出任意6颗（C(8,6)=28种），然后从2个蓝球中选出1颗（C(2,1)=2种），所以总共有28*2=56种可能。
*/

func DoubleColorBallCount(redBallCount int, blueBallCount int) int64 {
	if redBallCount <= 5 || blueBallCount <= 0 {
		return 0
	}
	
	nf := math.Factorial(int64(redBallCount))
	mf := math.Factorial(int64(6))
	nmf := math.Factorial(int64(redBallCount - 6))
	return nf / (mf * nmf) * int64(blueBallCount)
}
