package Utilities

import (
	"fmt"
	"math"
	"marvin/SnakeProGo/Utilities/Physics"
)

type Point struct {
	X, Y int
}
func (p *Point) ToVector() *Physics.Vector {
	return &Physics.Vector{float64(p.X),float64(p.Y),0}
}

func (p *Point) Print() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}
func (p *Point) Equals(p2 *Point) bool {
	if p.X == p2.X && p.Y == p2.Y {
		return true
	}
	return false
}
func (p *Point) Copy() *Point {
	return &Point{p.X, p.Y}
}
func (p *Point) CollidesWithWall(XTiles,YTiles float64) bool {
	if p.X >= 0 && p.Y >= 0 && p.X < int(XTiles) && p.Y < int(YTiles) {
		return false
	}
	return true
}
func (p *Point) CollidesWithPnts(pnts []*Point) *Point {
	for _,pnt := range(pnts) {
		if p.Equals(pnt) {
			return pnt
		}
	}
	return nil
}
func (p *Point) CollidesWithPnt(pnt *Point) bool {
	if p.Equals(pnt) {
		return true
	}
	return false
}
func (p *Point) IsContained(pnts []*Point) bool {
	for _,pnt := range(pnts) {
		if p.Equals(pnt) {
			return true
		}
	}
	return false
}
func (p *Point) GetDisTo(p2 *Point) float64 {
	return math.Sqrt(float64((p.X-p2.X)*(p.X-p2.X)+(p.Y-p2.Y)*(p.Y-p2.Y)))
}
func CollidePnts(pnts1,pnts2 []*Point) (collPnts []*Point) {
	for _,pnt1 := range(pnts1) {
		for _,pnt2 := range(pnts2) {
			if pnt1.Equals(pnt2) {
				collPnts = append(collPnts, pnt1)
			}
		}
	}
	return
}


type SnakeTiles struct {
	Ts []*Point
}

func (t *SnakeTiles) Print() string {
	out := "["
	for i,_ := range(t.Ts) {
		out += t.Ts[i].Print()
		if i < len(t.Ts)-1 {
			out += ","
		}
	}
	out += "]"
	return out
}
func (t *SnakeTiles) GetLast() *Point {
	if len(t.Ts) > 0 {
		return t.Ts[len(t.Ts)-1]
	}
	return &Point{0,0}
}
func (t *SnakeTiles) GetFirst() *Point {
	if len(t.Ts) > 0 {
		return t.Ts[0]
	}
	return &Point{0,0}
}

func (t *SnakeTiles) AddBack(xdir, ydir, tiles int) {
	var LastP *Point
	for i := 0; i < tiles; i++ {
		LastP = t.GetLast()
		t.Ts = append(t.Ts, &Point{LastP.X+xdir, LastP.Y+ydir})
	}
}
func (t *SnakeTiles) RemBack(tiles int) {
	t.Ts = t.Ts[:len(t.Ts)-tiles]
}

func (t *SnakeTiles) AddFront(xdir, ydir, tiles int) {
	FirstP := t.GetFirst()
	appList := make([]*Point, tiles)
	for i := tiles-1; i >= 0; i-- {
		appList[i] = &Point{FirstP.X+xdir, FirstP.Y+ydir}
		FirstP = appList[i]
	}
	newList := make([]*Point, 0)
	for i,_ := range(appList) {
		newList = append(newList, appList[i])
	}
	for i,_ := range(t.Ts) {
		newList = append(newList, t.Ts[i])
	}
	t.Ts = newList
}
func (t *SnakeTiles) RemFront(tiles int) {
	t.Ts = t.Ts[tiles:]
}
