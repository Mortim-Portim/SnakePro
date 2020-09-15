package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"fmt"
	"math"
)

//Laser

var (
	Laser_Item_Image, Laser_Item_Rotten_Image, Laser_Item_Image_Corner, Laser_Obj_Image, Laser_Smoke_Obj_Image *Utilities.ImageObj
)

type LaserFactory struct {
	prob float64
}

func (f *LaserFactory) Init() {
	f.prob = 1.0
	path := f_Resources+res_dir+f_Images+f_Items
	Laser_Item_Image = Utilities.LoadImgObj(path+f_Laser_img, TileWidth, TileHeight, 0, 0, 0)
	Laser_Item_Image_Corner = Utilities.LoadImgObj(path+f_Laser_img_c, TileWidth, TileHeight, 0, 0, 0)
	Laser_Obj_Image = Utilities.LoadImgObj(path+f_Shoot_img, TileWidth, TileHeight, 0, 0, 0)
	Laser_Smoke_Obj_Image = Utilities.LoadImgObj(path+f_Shoot_S_img, TileWidth, TileHeight, 0, 0, 0)
	Laser_Item_Rotten_Image = Utilities.LoadImgObj(path+f_Laser_Rotten_img, TileWidth, TileHeight, 0, 0, 0)
}

func (f *LaserFactory) GetRandomItemInBounds(W,H int) Item {
	X,Y := Utilities.GetRandomPosInScreen(W,H)
	s := &Laser{&Utilities.Point{X,Y}, false}
	//Do something
	
	return s
}
func (f *LaserFactory) GetProbability() float64 {
	return f.prob
}

type LaserObj struct {
	Owner *Player
	Points *Utilities.SnakeTiles
	Timer, FullTime int
	Rotation *Direction
	idx int
}

func (o *LaserObj) Update(g *InGame) {
	o.Timer --
	for i,_ := range(g.Pls) {
		plPnts := g.Pls[i].GetAllPoints()
		collPnts := Utilities.CollidePnts(plPnts, o.Points.Ts[:o.idx+1])
		if len(collPnts) > 0 {
			if g.Pls[i].GetHead().IsContained(collPnts) && g.Pls[i].ID != o.Owner.ID {
				g.Pls[i].OnCollision(LASER, 0, 1)
				g.Pls[i].SetDead(g.frame)
				o.Owner.OnCollision(LASER, 1, 0)
			}else{
				for s,_ := range(plPnts) {
					if plPnts[s].Equals(collPnts[len(collPnts)-1]) {
						if g.Pls[i].ID == o.Owner.ID {
							val := len(o.Owner.GetTiles().Ts)-int(o.Owner.GetSpeed()+1)
							if s < val {
								g.Pls[i].GetTiles().RemFront(s)
							}
						}else{
							g.Pls[i].GetTiles().RemFront(s)
						}
						break
					}
				}
			}
		}
	}
}
func (o *LaserObj) Exists() bool {
	if o.Timer <= 0 {
		return false
	}
	return true
}
func (o *LaserObj) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	if o.Timer > 0 {
		o.idx = len(o.Points.Ts)-int(math.Round((float64(o.Timer)/float64(o.FullTime))*float64(len(o.Points.Ts)-1)))-1
		for i := 0; i < o.idx; i++ {
			pnt, ok := (*mapper)[*o.Points.Ts[i]]
			if ok {
				Laser_Smoke_Obj_Image.X = float64(pnt.X)
				Laser_Smoke_Obj_Image.Y = float64(pnt.Y)
				Laser_Smoke_Obj_Image.Angle = o.Rotation.GetRotation()
				Laser_Smoke_Obj_Image.DrawImageObj(screen)
			}
		}
		
		pnt, ok := (*mapper)[*o.Points.Ts[o.idx]]
		if ok {
			Laser_Obj_Image.X = float64(pnt.X)
			Laser_Obj_Image.Y = float64(pnt.Y)
			Laser_Obj_Image.Angle = o.Rotation.GetRotation()
			Laser_Obj_Image.DrawImageObj(screen)
		}
	}
}

type Laser struct {
	Position *Utilities.Point
	IsRotten bool
}

func (s *Laser) Use(p *Player, g *InGame) {
	if p.GetState() == ALIVE {
		SoundEffects["Laser"].PR()
		
		laser := &Utilities.SnakeTiles{make([]*Utilities.Point, 1)}
		laser.Ts[0] = p.GetHead().Copy()
		
		laserDir := p.GetCurrDir()
		if s.IsRotten {
			laserDir = GetRandomDir()
		}
		for laserDir.Equals(p.GetCurrDir().Inverse()) {
			laserDir = GetRandomDir()
		}
		
		xdif, ydif := laserDir.GetXYDif()
		laser.AddBack(int(xdif), int(ydif), LaserLength+1)
		
		laser.RemFront(2)
		time := int(LaserTime*FPS)
		g.World.SObjs = append(g.World.SObjs, &LaserObj{p, laser, time,time,laserDir,0})
		
		p.CurrentItem = nil
	}
}
func (s *Laser) Pos() (*Utilities.Point) {
	return s.Position
}
func (s *Laser) OnCollision(ownTile, collTile Ptype, collPcs interface{}) {
	GH, ok := collPcs.(string)
	if ok && GH == "GH" {
		s.IsRotten = true
	}
}
func (s *Laser) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	pnt, ok := (*mapper)[*s.Position]
	if ok {
		if s.IsRotten {
			Laser_Item_Rotten_Image.X = float64(pnt.X)
			Laser_Item_Rotten_Image.Y = float64(pnt.Y)
			Laser_Item_Rotten_Image.DrawImageObj(screen)
		}else{
			Laser_Item_Image.X = float64(pnt.X)
			Laser_Item_Image.Y = float64(pnt.Y)
			Laser_Item_Image.DrawImageObj(screen)
		}
	}
}
func (s *Laser) DrawCorner(screen *ebiten.Image, pnt *Utilities.Point) {
	Laser_Item_Image_Corner.X = float64(pnt.X)
	Laser_Item_Image_Corner.Y = float64(pnt.Y)
	Laser_Item_Image_Corner.DrawImageObj(screen)
}