package usecase

import (
	"math"
)

type calcFunc func(float64, float64) float64

func calcLog(value, void float64) float64 {
	void = 0
	return math.Log10(value)
}

func calcLn(value, void float64) float64 {
	void = 0
	return math.Log(value)
}

func calcCos(value, void float64) float64 {
	void = 0
	return math.Cos(value)
}

func calcAcos(value, void float64) float64 {
	void = 0
	return math.Acos(value)
}

func calcSin(value, void float64) float64 {
	void = 0
	return math.Sin(value)
}

func calcAsin(value, void float64) float64 {
	void = 0
	return math.Asin(value)
}

func calcTan(value, void float64) float64 {
	void = 0
	return math.Tan(value)
}

func calcAtan(value, void float64) float64 {
	void = 0
	return math.Atan(value)
}

func calcSqrt(value, void float64) float64 {
	void = 0
	return math.Sqrt(value)
}

func calcPlus(x, y float64) float64 {
	return x + y
}

func calcMinus(x, y float64) float64 {
	return x - y
}

func calcMult(x, y float64) float64 {
	return x * y
}

func calcDiv(x, y float64) float64 {
	return x / y
}

func calcMod(x, y float64) float64 {
	return float64(int(x) % int(y))
}

func calcPow(x, y float64) float64 {
	return math.Pow(x, y)
}
