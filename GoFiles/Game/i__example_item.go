package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

//ITEMNAME

var (
	ITEMNAME_Item_Image, ITEMNAME_Item_Image_Corner *Utilities.ImageObj
	
	//Needed:
	f_ITEMNAME_img, f_ITEMNAME_img_c string
	//SoundEffects["ITEMNAME"]
	//Statistic.AddHead("ITEMNAME", 1)
)

type ITEMNAMEFactory struct {
	prob float64
}

func (f *ITEMNAMEFactory) Init() {
	f.prob = 0.8
	path := f_Resources+res_dir+f_Images+f_Items
	ITEMNAME_Item_Image = Utilities.LoadImgObj(path+f_ITEMNAME_img, TileWidth, TileHeight, 0, 0, 0)
	ITEMNAME_Item_Image_Corner = Utilities.LoadImgObj(path+f_ITEMNAME_img_c, TileWidth, TileHeight, 0, 0, 0)
}

func (f *ITEMNAMEFactory) GetRandomItemInBounds(W,H int) Item {
	X,Y := Utilities.GetRandomPosInScreen(W,H)
	s := &ITEMNAME{&Utilities.Point{X,Y}}
	//Do something
	
	return s
}
func (f *ITEMNAMEFactory) GetProbability() float64 {
	return f.prob
}


type ITEMNAME struct {
	Position *Utilities.Point
}

func (s *ITEMNAME) Use(p *Player, g *InGame) {
	SoundEffects["ITEMNAME"].PR()
	
	//Do something
	
	p.CurrentItem = nil
}
func (s *ITEMNAME) Pos() (*Utilities.Point) {
	return s.Position
}
func (s *ITEMNAME) OnCollision(ownTile, collTile Ptype, collPcs interface{}) {
	
}
func (s *ITEMNAME) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	pnt, ok := (*mapper)[*s.Position]
	if ok {
		ITEMNAME_Item_Image.X = float64(pnt.X)
		ITEMNAME_Item_Image.Y = float64(pnt.Y)
		ITEMNAME_Item_Image.DrawImageObj(screen)
	}
}
func (s *ITEMNAME) DrawCorner(screen *ebiten.Image, pnt *Utilities.Point) {
	ITEMNAME_Item_Image_Corner.X = float64(pnt.X)
	ITEMNAME_Item_Image_Corner.Y = float64(pnt.Y)
	ITEMNAME_Item_Image_Corner.DrawImageObj(screen)
}