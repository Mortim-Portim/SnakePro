package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

var (
	Speed_Item_Image, Speed_Item_Image_Corner, Speed_Item_Rotten_Image *Utilities.ImageObj
)

type SpeedFactory struct {
	prob float64
}

func (f *SpeedFactory) Init() {
	f.prob = 1.0
	path := f_Resources+res_dir+f_Images+f_Items
	Speed_Item_Image = Utilities.LoadImgObj(path+f_Speed_img, TileWidth, TileHeight, 0, 0, 0)
	Speed_Item_Image_Corner = Utilities.LoadImgObj(path+f_Speed_img_c, TileWidth, TileHeight, 0, 0, 0)
	Speed_Item_Rotten_Image = Utilities.LoadImgObj(path+f_Speed_Rotten_img, TileWidth, TileHeight, 0, 0, 0)
}

func (f *SpeedFactory) GetRandomItemInBounds(W,H int) Item {
	X,Y := Utilities.GetRandomPosInScreen(W,H)
	s := &Speed{&Utilities.Point{X,Y}, false}
	return s
}
func (f *SpeedFactory) GetProbability() float64 {
	return f.prob
}
type SpeedObj struct {
	Owner *Player
	Timer, fullTime int
	IsRotten bool
	RottenDec int
}
func (o *SpeedObj) Update(g *InGame) {
	o.Timer --
	if o.Timer <= 0 {
		o.Owner.SetSpeed(PlayerSpeed)
	}else{
		o.Owner.SetSpeed(SpeedStrength)
	}
	if o.IsRotten && o.RottenDec > 0 {
		if Utilities.GetRandomFloat(0,1) < 2.0/float64(o.fullTime) {
			o.RottenDec --
			if Utilities.GetRandomFloat(0,1) < 0.5 {
				o.Owner.snake.nextDir = o.Owner.snake.currentDir.Left()
			}else{
				o.Owner.snake.nextDir = o.Owner.snake.currentDir.Right()
			}
		}
	}
	if g.frame%int(FPS/o.Owner.snake.speed) == 0 {
		o.Owner.snake.UpdateCurrDir()
		o.Owner.snake.UpdateSnakeHead()
		o.Owner.snake.UpdateSnakeTail()
		o.Owner.snake.compAllPieces()
	}
	g.World.CheckCollisionPlayer(o.Owner)
}
func (o *SpeedObj) Exists() bool {
	if o.Timer <= 0 {
		o.Owner.SetSpeed(PlayerSpeed)
		return false
	}
	return true
}
func (o *SpeedObj) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {}


type Speed struct {
	Position *Utilities.Point
	IsRotten bool
}

func (s *Speed) Use(p *Player, g *InGame) {
	if p.GetState() == ALIVE {
		SoundEffects["Speed"].PR()
		
		speed := &SpeedObj{p, int(FPS*SpeedLength), int(FPS*SpeedLength), s.IsRotten, 3}
		
		g.World.SObjs = append(g.World.SObjs, speed)
		
		p.CurrentItem = nil
	}
}
func (s *Speed) Pos() (*Utilities.Point) {
	return s.Position
}
func (s *Speed) OnCollision(ownTile, collTile Ptype, collPcs interface{}) {
	GH, ok := collPcs.(string)
	if ok && GH == "GH" {
		s.IsRotten = true
	}
}
func (s *Speed) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	pnt, ok := (*mapper)[*s.Position]
	if ok {
		if s.IsRotten {
			Speed_Item_Rotten_Image.X = float64(pnt.X)
			Speed_Item_Rotten_Image.Y = float64(pnt.Y)
			Speed_Item_Rotten_Image.DrawImageObj(screen)
		}else{
			Speed_Item_Image.X = float64(pnt.X)
			Speed_Item_Image.Y = float64(pnt.Y)
			Speed_Item_Image.DrawImageObj(screen)
		}
	}
}
func (s *Speed) DrawCorner(screen *ebiten.Image, pnt *Utilities.Point) {
	Speed_Item_Image_Corner.X = float64(pnt.X)
	Speed_Item_Image_Corner.Y = float64(pnt.Y)
	Speed_Item_Image_Corner.DrawImageObj(screen)
}




