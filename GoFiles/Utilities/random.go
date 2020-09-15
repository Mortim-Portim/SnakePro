package Utilities

import (
	"math/rand"
	"math"
)
func GetRandomFloat(l,u float64) float64 {
	return rand.Float64()*(u-l)+l
}
func GetRandomInt(l,u float64) int {
	return int(math.Round(GetRandomFloat(l,u)))
}
func GetRandomPosInScreen(W,H int) (int, int) {
	return GetRandomPos(0,0,W,H)
}
func GetRandomPosOutScreen(W,H int) (int, int) {
	X, Y := GetRandomPos(0,0,W,H)
	if GetRandomFloat(0,1) < 0.5 {
		if GetRandomFloat(0,1) < 0.5 {
			Y = 1
		}else {
			Y = H-1
		}
	}else {
		if GetRandomFloat(0,1) < 0.5 {
			X = 1
		}else {
			X = W-1
		}
	}
	return X,Y
}
func GetRandomPos(x,y,w,h int) (int, int) {
	return int(GetRandomFloat(float64(x),float64(x+w))), int(GetRandomFloat(float64(y),float64(y+h)))
}