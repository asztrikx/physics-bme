package main

import (
	"math"
	"math/big"
)

func degreeToRadian(degree *big.Float) *big.Float {
	// radian/(2π) == degree/360

	radian := bigfloatCreate(0).Quo(degree, bigfloatCreate(360))

	//Multiply separately to keep pi's high precision
	radian.Mul(radian, bigfloatCreate(math.Pi))
	radian.Mul(radian, bigfloatCreate(2))

	return radian
}

func radianToDegree(radian *big.Float) *big.Float {
	// radian/(2π) == degree/360

	//Quo separately to keep pi's high precision
	degree := bigfloatCreate(0).Quo(radian, bigfloatCreate(math.Pi))
	degree.Quo(degree, bigfloatCreate(2.0))

	degree.Mul(degree, bigfloatCreate(360))

	return degree
}
