package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

//Fart

var (
	Fart_Item_Image, Fart_Item_Rotten_Image, Fart_Obj_Image, Fart_Obj_Rotten_Image, Fart_Item_Image_Corner *Utilities.ImageObj
)

type FartFactory struct {
	prob float64
}
func (f *FartFactory) Init() {
	f.prob = 0.8
	path := f_Resources+res_dir+f_Images+f_Items
	Fart_Item_Image = Utilities.LoadImgObj(path+f_Fart_img, TileWidth, TileHeight, 0, 0, 0)
	Fart_Item_Image_Corner = Utilities.LoadImgObj(path+f_Fart_img_c, TileWidth, TileHeight, 0, 0, 0)
	Fart_Item_Rotten_Image = Utilities.LoadImgObj(path+f_Fart_Rotten_img, TileWidth, TileHeight, 0, 0, 0)
	Fart_Obj_Image = Utilities.LoadImgObj(path+f_Fart_Obj_img, (TileWidth/SubImageScale)*float64(FartRadius*2+1), (TileHeight/SubImageScale)*float64(FartRadius*2+1), 0, 0, 0)
	Fart_Obj_Rotten_Image = Utilities.LoadImgObj(path+f_Fart_Obj_Rotten_img, (TileWidth/SubImageScale)*float64(FartRadius*2+1), (TileHeight/SubImageScale)*float64(FartRadius*2+1), 0, 0, 0)
}
func (f *FartFactory) GetRandomItemInBounds(W,H int) Item {
	X,Y := Utilities.GetRandomPosInScreen(W,H)
	s := &Fart{&Utilities.Point{X,Y}, false}
	//Do something
	
	return s
}
func (f *FartFactory) GetProbability() float64 {
	return f.prob
}


type FartObj struct {
	Owner *Player
	Timer, fullTime int
	Position *Utilities.Point
	fartPnts []*Utilities.Point
	drawX, drawY float64
	IsRotten bool
}
func (o *FartObj) Update(g *InGame) {
	o.Timer --
	if g.frame%int(FPS/FartStrength) == 0{
		for i,_ := range(g.Pls) {
			if g.Pls[i] != o.Owner {
				plPnts := g.Pls[i].GetAllPoints()
				collPnts := Utilities.CollidePnts(plPnts, o.fartPnts)
				if len(collPnts) > 0 {
					if g.Pls[i].GetHead().IsContained(collPnts) {
						g.Pls[i].snake.tiles.RemFront(1)
						if len(g.Pls[i].snake.tiles.Ts) < 1 {
							g.Pls[i].OnCollision(FART, 0, 1)
							g.Pls[i].SetDead(g.World.frame)
							if o.Owner != nil {
								o.Owner.OnCollision(FART, 1, 0)
							}
						}
						g.Pls[i].snake.compAllPieces()
					}
				}
			}
		}
	}
}
func (o *FartObj) Exists() bool {
	if o.Timer <= 0 {
		return false
	}
	return true
}
func MinuxXSqurdPOne(x float64) float64 {
	xp := 2*x-1
	return -(xp*xp)+1
}
func (o *FartObj) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	if o.Timer > 0 {
		alpha := MinuxXSqurdPOne(float64(o.fullTime-o.Timer)/float64(o.fullTime))
		
		if !o.IsRotten {
			Fart_Obj_Image.X = o.drawX
			Fart_Obj_Image.Y = o.drawY
			Fart_Obj_Image.DrawImageObjAlpha(screen, alpha)
		}else{
			Fart_Obj_Rotten_Image.X = o.drawX
			Fart_Obj_Rotten_Image.Y = o.drawY
			Fart_Obj_Rotten_Image.DrawImageObjAlpha(screen, alpha)
		}
	}
}



type Fart struct {
	Position *Utilities.Point
	IsRotten bool
}

func (s *Fart) Use(p *Player, g *InGame) {
	if p.GetState() == ALIVE {
		SoundEffects["Farting"].PR()
		
		tl := p.GetTail()
		pl := p
		if s.IsRotten {
			xdif, ydif := p.GetCurrDir().GetXYDif()
			hdPnt := p.GetHead()
			tl = &Utilities.Point{hdPnt.X+int(xdif*float64(FartRadius)*1.2), hdPnt.Y+int(ydif*float64(FartRadius)*1.2)}
			pl = nil
		}
		
		fartPnts := make([]*Utilities.Point, 0)
		for x := -FartRadius; x <= FartRadius; x++ {
			for y := -FartRadius; y <= FartRadius; y++ {
				fartPnts = append(fartPnts, &Utilities.Point{x+tl.X,y+tl.Y})
			}
		}
		Xp := (float64(tl.X-FartRadius))*(TileWidth/SubImageScale)
		Yp := (float64(tl.Y-FartRadius))*(TileHeight/SubImageScale)
		fart := &FartObj{pl, int(FartTime*FPS), int(FartTime*FPS), tl, fartPnts, Xp, Yp, s.IsRotten}
		
		g.World.SObjs = append(g.World.SObjs, fart)
		
		p.CurrentItem = nil
	}
}
func (s *Fart) Pos() (*Utilities.Point) {
	return s.Position
}
func (s *Fart) OnCollision(ownTile, collTile Ptype, collPcs interface{}) {
	GH, ok := collPcs.(string)
	if ok && GH == "GH" {
		s.IsRotten = true
	}
}
func (s *Fart) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	pnt, ok := (*mapper)[*s.Position]
	if ok {
		if s.IsRotten {
			Fart_Item_Rotten_Image.X = float64(pnt.X)
			Fart_Item_Rotten_Image.Y = float64(pnt.Y)
			Fart_Item_Rotten_Image.DrawImageObj(screen)
		}else{
			Fart_Item_Image.X = float64(pnt.X)
			Fart_Item_Image.Y = float64(pnt.Y)
			Fart_Item_Image.DrawImageObj(screen)
		}
	}
}
func (s *Fart) DrawCorner(screen *ebiten.Image, pnt *Utilities.Point) {
	Fart_Item_Image_Corner.X = float64(pnt.X)
	Fart_Item_Image_Corner.Y = float64(pnt.Y)
	Fart_Item_Image_Corner.DrawImageObj(screen)
}