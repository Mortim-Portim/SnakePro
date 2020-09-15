package Game

import (
	"math"
	"marvin/SnakeProGo/Utilities"
)

type Direction struct { //0=oben, 1=rechts, 2=unten, 3=links
	Type int
}
func GetRandomDir() *Direction {
	rand := Utilities.GetRandomFloat(0,1)
	if rand < 0.25 {
		return &Direction{0}
	}else if rand < 0.5 {
		return &Direction{1}
	}else if rand < 0.75 {
		return &Direction{2}
	}else if rand < 1.0 {
		return &Direction{3}
	}
	return nil
}
func (d *Direction) Copy() *Direction {
	dir := Direction{d.Type}
	return &dir
}
func (d *Direction) Equals(d2 *Direction) bool {
	return d.Type == d2.Type
}
func (d *Direction) Inverse() *Direction {
	dir := Direction{d.Type}
	dir.Type += 2
	if dir.Type > 3 {
		dir.Type -= 4
	}
	return &dir
}
func (d *Direction) Left() *Direction {
	dir := Direction{d.Type}
	dir.Type -= 1
	if dir.Type < 0 {
		dir.Type = 3
	}
	return &dir
}
func (d *Direction) Right() *Direction {
	dir := Direction{d.Type}
	dir.Type += 1
	if dir.Type > 3 {
		dir.Type = 0
	}
	return &dir
}
func GetDirFromXYDif(xdif, ydif float64) *Direction {
	dir := Direction{0}
	if math.Abs(xdif) > math.Abs(ydif) {
		if xdif > 0 {
			dir.Type = 1
		}else{
			dir.Type = 3
		}
	}else{
		if ydif > 0 {
			dir.Type = 2
		}else{
			dir.Type = 0
		}
	}
	return &dir
}
func (d *Direction) GetXYDif() (float64, float64) {
	switch d.Type {
	case 1:
		return 1, 0
	case 2:
		return 0, 1
	case 3:
		return -1, 0
	}
	return 0, -1
}
func DetDirFrom2Ts(p1, p2 *Utilities.Point) *Direction {
	return GetDirFromXYDif(float64(p2.X-p1.X), float64(p2.Y-p1.Y))
}
func (d *Direction) IsLeft(d2 *Direction) bool {
	newDir := d.Left()
	if newDir.Type == d2.Type {
		return true
	}
	return false
}
func (d *Direction) IsRight(d2 *Direction) bool {
	newDir := d.Right()
	if newDir.Type == d2.Type {
		return true
	}
	return false
}
func (d *Direction) IsInvers(d2 *Direction) bool {
	newDir := d.Inverse()
	if newDir.Type == d2.Type {
		return true
	}
	return false
}
func (d *Direction) GetRotation() float64 {
	switch d.Type {
	case 1:
		return 90
	case 2:
		return 180
	case 3:
		return 270
	}
	return 0
}