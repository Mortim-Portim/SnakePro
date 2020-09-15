package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

//Revive

var (
	Revive_Item_Image, Revive_Item_Image_Corner *Utilities.ImageObj
)

type ReviveFactory struct {
	prob float64
}

func (f *ReviveFactory) Init() {
	f.prob = 0.6
	path := f_Resources+res_dir+f_Images+f_Items
	Revive_Item_Image = Utilities.LoadImgObj(path+f_Revive_img, TileWidth, TileHeight, 0, 0, 0)
	Revive_Item_Image_Corner = Utilities.LoadImgObj(path+f_Revive_img_c, TileWidth, TileHeight, 0, 0, 0)
}

func (f *ReviveFactory) GetRandomItemInBounds(W,H int) Item {
	X,Y := Utilities.GetRandomPosInScreen(W,H)
	s := &Revive{&Utilities.Point{X,Y}}
	//Do something
	
	return s
}
func (f *ReviveFactory) GetProbability() float64 {
	return f.prob
}

type ReviveObj struct {
	Position *Utilities.Point
	RemainingTime, FullTime int
	parent *Player
}
func (o *ReviveObj) Update(g *InGame) {
	o.RemainingTime --
	if o.RemainingTime == 0 {
		SoundEffects["Revive"].PR()
		o.parent.Reset(o.parent.ID, o.parent.ConID, o.parent.GetSpeed(), ALIVE)
		o.parent.SetUndestroyable(g.frame+int(FPS*ReviveImmortalTime))
	}
}
func (o *ReviveObj) Exists() bool {
	if o.RemainingTime <= 0 {
		return false
	}
	return true
}
func (o *ReviveObj) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	newPnt := (*mapper)[*o.Position]
	Revive_Item_Image.X = float64(newPnt.X)
	Revive_Item_Image.Y = float64(newPnt.Y)
	Revive_Item_Image.DrawImageObjAlpha(screen, float64(o.RemainingTime)/float64(o.FullTime))
}

type Revive struct {
	Position *Utilities.Point
}

func (s *Revive) Use(p *Player, g *InGame) {
	if p.GetState() == GHOST {
		p.SetState(WAITING)
	}
	
	g.World.SObjs = append(g.World.SObjs, &ReviveObj{p.ItemPos, FPS*1.5, FPS*1.5, p})
	p.CurrentItem = nil
}
func (s *Revive) Pos() (*Utilities.Point) {
	return s.Position
}
func (s *Revive) OnCollision(ownTile, collTile Ptype, collPcs interface{}) {
	
}
func (s *Revive) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	pnt, ok := (*mapper)[*s.Position]
	if ok {
		Revive_Item_Image.X = float64(pnt.X)
		Revive_Item_Image.Y = float64(pnt.Y)
		Revive_Item_Image.DrawImageObj(screen)
	}
}
func (s *Revive) DrawCorner(screen *ebiten.Image, pnt *Utilities.Point) {
	Revive_Item_Image_Corner.X = float64(pnt.X)
	Revive_Item_Image_Corner.Y = float64(pnt.Y)
	Revive_Item_Image_Corner.DrawImageObj(screen)
}