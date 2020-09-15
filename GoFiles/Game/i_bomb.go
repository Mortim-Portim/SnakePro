package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

//Bomb

var (
	Bomb_Item_Image,Bomb_Item_Rotten_Image, Bomb_Item_Image_Corner, Bomb_Obj_Image_Radius, Bomb_Obj_Image *Utilities.ImageObj
	
	//Needed:
	//f_Bomb_img, f_Bomb_img_c string
	//SoundEffects["Bomb"]
)

type BombFactory struct {
	prob float64
}

func (f *BombFactory) Init() {
	f.prob = 1.0
	path := f_Resources+res_dir+f_Images+f_Items
	Bomb_Item_Image = Utilities.LoadImgObj(path+f_Bomb_img, TileWidth, TileHeight, 0, 0, 0)
	Bomb_Item_Image_Corner = Utilities.LoadImgObj(path+f_Bomb_img_c, TileWidth, TileHeight, 0, 0, 0)
	Bomb_Obj_Image_Radius = Utilities.LoadImgObj(path+f_Bomb_img_r, (TileWidth/SubImageScale)*float64(BombRadius*2+1), (TileHeight/SubImageScale)*float64(BombRadius*2+1), 0, 0, 0)
	Bomb_Obj_Image = Utilities.LoadImgObj(path+f_Bomb_Obj_img, TileWidth, TileHeight, 0, 0, 0)
	Bomb_Item_Rotten_Image = Utilities.LoadImgObj(path+f_Bomb_Rotten_img, TileWidth, TileHeight, 0, 0, 0)
}

func (f *BombFactory) GetRandomItemInBounds(W,H int) Item {
	X,Y := Utilities.GetRandomPosInScreen(W,H)
	s := &Bomb{&Utilities.Point{X,Y}, false}
	//Do something
	
	return s
}
func (f *BombFactory) GetProbability() float64 {
	return f.prob
}

type BombObj struct {
	Owner *Player
	Timer, fullTime int
	Position *Utilities.Point
}
func (o *BombObj) Update(g *InGame) {
	o.Timer --
	if o.Timer == 0 {
		o.Explode(g)
		return
	}
}
func (o *BombObj) Exists() bool {
	if o.Timer <= 0 {
		return false
	}
	return true
}
func (o *BombObj) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	if o.Timer > 0 {
		alpha := float64(o.fullTime-o.Timer)/float64(o.fullTime)
		Bomb_Obj_Image_Radius.X = (float64(o.Position.X-BombRadius))*(TileWidth/SubImageScale)
		Bomb_Obj_Image_Radius.Y = (float64(o.Position.Y-BombRadius))*(TileHeight/SubImageScale)
		Bomb_Obj_Image_Radius.DrawImageObjAlpha(screen, alpha)
		
		pnt, ok := (*mapper)[*o.Position]
		if ok {
			Bomb_Obj_Image.X = float64(pnt.X)
			Bomb_Obj_Image.Y = float64(pnt.Y)
			Bomb_Obj_Image.DrawImageObj(screen)
		}
	}
}
func (o *BombObj) Explode(g *InGame) {
	SoundEffects["Explosion"].sounds[1].Play()
	bombPnts := make([]*Utilities.Point, 0)
	for x := -BombRadius; x <= BombRadius; x++ {
		for y := -BombRadius; y <= BombRadius; y++ {
			bombPnts = append(bombPnts, &Utilities.Point{x+o.Position.X,y+o.Position.Y})
		}
	}
	for i,_ := range(g.Pls) {
		plPnts := g.Pls[i].GetAllPoints()
		collPnts := Utilities.CollidePnts(plPnts, bombPnts)
		if len(collPnts) > 0 {
			
			if g.Pls[i].GetHead().IsContained(collPnts) {
				g.Pls[i].OnCollision(BOMB, 0, 1)
				g.Pls[i].SetDead(g.World.frame)
				if g.Pls[i].ID != o.Owner.ID {
					o.Owner.OnCollision(BOMB, 1, 0)
				}
			}else{
				for s,_ := range(plPnts) {
					if plPnts[s].Equals(collPnts[len(collPnts)-1]) {
						g.Pls[i].GetTiles().RemFront(s)
						break
					}
				}
			}
		}
	}
}

type Bomb struct {
	Position *Utilities.Point
	IsRotten bool
}

func (s *Bomb) Use(p *Player, g *InGame) {
	if p.GetState() == ALIVE {
		SoundEffects["Explosion"].sounds[0].Play()
		
		head := p.GetHead().Copy()
		
		xdif, ydif := p.GetCurrDir().GetXYDif()
		
		if s.IsRotten {
			xdif, ydif = GetRandomDir().GetXYDif()
		}
		
		bombPnt := &Utilities.Point{head.X+int(xdif)*BombLength, head.Y+int(ydif)*BombLength}
		
		time := int(BombTime*FPS)
		bomb := &BombObj{p, time, time, bombPnt}
		
		g.World.SObjs = append(g.World.SObjs, bomb)
		
		p.CurrentItem = nil
	}
}
func (s *Bomb) Pos() (*Utilities.Point) {
	return s.Position
}
func (s *Bomb) OnCollision(ownTile, collTile Ptype, collPcs interface{}) {
	GH, ok := collPcs.(string)
	if ok && GH == "GH" {
		s.IsRotten = true
	}
}
func (s *Bomb) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	pnt, ok := (*mapper)[*s.Position]
	if ok {
		if s.IsRotten {
			Bomb_Item_Rotten_Image.X = float64(pnt.X)
			Bomb_Item_Rotten_Image.Y = float64(pnt.Y)
			Bomb_Item_Rotten_Image.DrawImageObj(screen)
		}else{
			Bomb_Item_Image.X = float64(pnt.X)
			Bomb_Item_Image.Y = float64(pnt.Y)
			Bomb_Item_Image.DrawImageObj(screen)
		}
	}
}
func (s *Bomb) DrawCorner(screen *ebiten.Image, pnt *Utilities.Point) {
	Bomb_Item_Image_Corner.X = float64(pnt.X)
	Bomb_Item_Image_Corner.Y = float64(pnt.Y)
	Bomb_Item_Image_Corner.DrawImageObj(screen)
}