package Utilities

import (
	"fmt"
	"image/color"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/vector"
	"marvin/SnakeProGo/Utilities/Physics"
)
type Points []*Physics.Vector
func (ps *Points) Fill(screen *ebiten.Image, col color.Color) {
	var path vector.Path
	for i,p := range(*ps) {
		if i == 0 {
			path.MoveTo(float32(p.X), float32(p.Y))
		}else{
			path.LineTo(float32(p.X), float32(p.Y))
		}
	}
	op := &vector.FillOptions{
		Color: col,
	}
	path.Fill(screen, op)
}

func DrawBoxToImage(img *ebiten.Image, x, y, w, h float64, col color.Color) {
	box := make([]*Physics.Vector, 4)
	box[0] = &Physics.Vector{x,y,0}
	box[1] = &Physics.Vector{x+w,y,0}
	box[2] = &Physics.Vector{x+w,y+h,0}
	box[3] = &Physics.Vector{x,y+h,0}
	Box := Points(box)
	Box.Fill(img, col)
}
func DrawLine(img *ebiten.Image, pnt1, pnt2 *Point, w float64, col color.Color) {
	ps := make([]*Physics.Vector, 4)
	loc := pnt1.ToVector()
	Vec := pnt2.ToVector().Sub(loc)
	right := Vec.CrossProduct(&Physics.Vector{0,0,1}).Normalize().Mul(w)
	left := Vec.CrossProduct(&Physics.Vector{0,0,-1}).Normalize().Mul(w)
	ps[0] = loc.Add(right)
	ps[1] = loc.Add(right).Add(Vec)
	ps[2] = loc.Add(left).Add(Vec)
	ps[3] = loc.Add(left)
	newPs := Points(ps)
	
	newPs.Fill(img, col)
}

func MaxFloat(vals ...float64) float64 {
	maxV := 0.0; max := -1
	for i,v := range(vals) {
		if v > maxV || max < 0 {
			max = i
			maxV = v
		}
	}
	return maxV
}

func GetBarChart(Width, Heigth, X, Y float64, backCol color.Color, colors []color.RGBA, values ...float64) *ImageObj {
	graph := &ImageObj{X:X, Y:Y, W:Width, H:Heigth, Angle:0}
	graphBack, _ := ebiten.NewImage(int(Width), int(Heigth), ebiten.FilterDefault)
	graphBack.Fill(backCol)
	
	bars := len(values)
	bW := Width/(0.5+float64(bars)*1.5)
	bD := bW/4
	
	HeigthFac := (Heigth-2*bD)/MaxFloat(values...)
	
	DrawLine(graphBack, &Point{int(bD),int(bD)}, &Point{int(bD),int(Heigth-bD)}, Width/80, color.RGBA{0, 0, 0, 255})
	DrawLine(graphBack, &Point{int(bD),int(Heigth-bD)}, &Point{int(Width-bD),int(Heigth-bD)}, Width/80, color.RGBA{0, 0, 0, 255})
	
	for i,v := range(values) {
		x := 0.75*bW+bW*1.5*float64(i)
		y := Heigth-bD-v*HeigthFac
		DrawBoxToImage(graphBack, x, y, bW, v*HeigthFac, colors[i])
		if y > Heigth/4*3 {
			y = Heigth-bD*1.6
		}
		vStr := fmt.Sprintf("%0.0f",v)
		PrintTextAt(graphBack, vStr, int(x+bW/2), int(y+bD*1.4), 255, 255, 255, GetSizeForText(vStr, int(bW*1.5),int(bD)))
	}
	
	graph.Img = graphBack
	return graph
}