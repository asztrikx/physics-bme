package main

import "math/big"

func bigfloatCreate(f float64) *big.Float {
	return big.NewFloat(f).SetPrec(1024)
}
