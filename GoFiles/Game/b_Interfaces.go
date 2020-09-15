package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
)

type Ptype int

const (
	P_SnkH = Ptype(0)
	P_SnkB = Ptype(1)
	P_SnkT = Ptype(2)
	P_Apl = Ptype(3)
	P_Itm = Ptype(4)
	P_Wall = Ptype(5)
)

type Piece interface {
	Ptype() (Ptype)
	GetImg() *Utilities.ImageObj 
}
func DrawPieceToScreen(ImageObj *Utilities.ImageObj, screen *ebiten.Image, pnt *Utilities.Point, mapper *map[Utilities.Point]Utilities.Point) {
	newPnt, ok := (*mapper)[Utilities.Point{pnt.X, pnt.Y}]
	if ok {
		ImageObj.X = float64(newPnt.X);ImageObj.Y = float64(newPnt.Y)
		ImageObj.DrawImageObj(screen)
	}
}

type ItemFactory interface {
	GetRandomItemInBounds(W,H int) Item
	GetProbability() float64
}

type Item interface {
	Use(*Player, *InGame)
	Pos() (*Utilities.Point)
	OnCollision(ownTile, collTile Ptype, collPcs interface{})
	Draw(*ebiten.Image, *map[Utilities.Point]Utilities.Point)
	DrawCorner(*ebiten.Image, *Utilities.Point)
}
type StaticObj interface {
	Update(*InGame)
	Draw(*ebiten.Image, *map[Utilities.Point]Utilities.Point)
	Exists() bool
}