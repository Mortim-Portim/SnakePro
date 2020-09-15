package Utilities

import (
	
)

func IsInsideBox(bx, by, bw, bh, ox, oy, ow, oh float64) bool {
	if ox-ow/2.0 > bx && oy-oh/2.0 > by && ox+ow/2.0 < bx + bw && oy+oh/2.0 < by + bh {
		return true
	}
	return false
}

func GetStartPos4P(ID int, ScreenWidth, ScreenHeight, WallDis float64) (float64, float64, int) {
	W := ScreenWidth;H := ScreenHeight
	XPos := 0.0;YPos := 0.0
	XSteps := W/15.0*WallDis
	YSteps := H/15.0*WallDis
	dir := 0
	switch ID {
	case 0:
		XPos += XSteps
		YPos += YSteps
		dir = 1
	case 1:
		XPos += W-XSteps
		YPos += YSteps
		dir = 3
	case 2:
		XPos += W-XSteps
		YPos += H-YSteps
		dir = 3
	case 3:
		XPos += XSteps
		YPos += H-YSteps
		dir = 1
	default:
		XPos += W/2
		YPos += H/2
	}
	return float64(XPos), float64(YPos), dir
}