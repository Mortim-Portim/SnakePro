package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

//Bot

var (
	Bot_Item_Image, Bot_Item_Rotten_Image, Bot_Item_Image_Corner, Bot_Obj_Image_H, Bot_Obj_Image_B *Utilities.ImageObj
)

type BotFactory struct {
	prob float64
}

func (f *BotFactory) Init() {
	f.prob = 1.0
	path := f_Resources+res_dir+f_Images+f_Items
	Bot_Item_Image = Utilities.LoadImgObj(path+f_Bot_img, TileWidth, TileHeight, 0, 0, 0)
	Bot_Item_Image_Corner = Utilities.LoadImgObj(path+f_Bot_img_c, TileWidth, TileHeight, 0, 0, 0)
	Bot_Obj_Image_H = Utilities.LoadImgObj(f_Resources+res_dir+f_Images+f_Snks+f_Bot+f_SnkH+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	Bot_Obj_Image_B = Utilities.LoadImgObj(f_Resources+res_dir+f_Images+f_Snks+f_Bot+f_SnkB+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	Bot_Item_Rotten_Image = Utilities.LoadImgObj(path+f_Bot_Rotten_img, TileWidth, TileHeight, 0, 0, 0)
}

func (f *BotFactory) GetRandomItemInBounds(W,H int) Item {
	X,Y := Utilities.GetRandomPosInScreen(W,H)
	s := &Bot{&Utilities.Point{X,Y}, false}
	return s
}
func (f *BotFactory) GetProbability() float64 {
	return f.prob
}

type BotSnake struct {
	Snake
	Owner *Player
	IsRotten bool
	Timer int
}
func (o *BotSnake) Init(p *Player, g *InGame) {
	o.Owner = p
	b_Hd := []*Utilities.ImageObj{Bot_Obj_Image_H}
	b_Bd := []*Utilities.ImageObj{Bot_Obj_Image_B}
	o.Snake = *CreateSnake(g, b_Hd,b_Bd,b_Bd,b_Bd,b_Bd)
	X, Y := p.GetHead().X, p.GetHead().Y
	o.Snake.Reset(BotSpeed, 0, X, Y, BotLength, p.GetCurrDir().Copy())
}
func (o *BotSnake) Update(g *InGame) {
	o.Timer --
	o.Snake.nextDir = o.GetNextDir(g)
	
	if o.IsRotten {
		if Utilities.GetRandomFloat(0,1) < 0.5 {
			o.Snake.nextDir = DetDirFrom2Ts(o.Snake.GetHead(), o.Owner.GetHead())
		}
	}
	o.Snake.Update(g.frame)
	
	for i,_ := range(g.Pls) {
		if g.Pls[i].ID != o.Owner.ID && g.Pls[i].GetState() == ALIVE {
			plPnts := g.Pls[i].GetAllPoints()
			collPnts := Utilities.CollidePnts(plPnts, o.Snake.GetAllPoints())
			if len(collPnts) > 0 {
				if g.Pls[i].GetHead().IsContained(collPnts) {
					g.Pls[i].OnCollision(BOT, 0, 1)
					g.Pls[i].SetDead(g.frame)
					o.Owner.OnCollision(BOT, 1, 0)
				}
				if o.Snake.GetHead().IsContained(collPnts) {
					o.Snake.state = 1
				}
			}
		}
	}
}
func (o *BotSnake) Exists() bool {
	if o.Snake.state != 0 {
		return false
	}
	if o.Timer <= 0 {
		return false
	}
	return true
}
func (o *BotSnake) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	o.Snake.Draw(screen, mapper)
}

func (o *BotSnake) GetNextDir(g *InGame) *Direction {
	return DetDirFrom2Ts(o.Snake.GetHead(), o.GetClosestPlayer(g).GetHead())
}
func (o *BotSnake) GetClosestPlayer(g *InGame) *Player {
	pnt := o.Snake.GetHead()
	closestPlayer := g.Pls[0]
	closestDis := -1.0
	for i,_ := range(g.Pls) {
		if g.Pls[i] != o.Owner && g.Pls[i].GetState() == ALIVE {
			hd := g.Pls[i].GetHead()
			dis := hd.GetDisTo(pnt)
			if dis < closestDis || closestDis < 0 {
				closestDis = dis
				closestPlayer = g.Pls[i]
			}
		}
	}
	return closestPlayer
}

type Bot struct {
	Position *Utilities.Point
	IsRotten bool
}

func (s *Bot) Use(p *Player, g *InGame) {
	if p.GetState() == ALIVE {
		Bot := &BotSnake{IsRotten:s.IsRotten, Timer:FPS*10}
		Bot.Init(p,g)
		g.World.SObjs = append(g.World.SObjs, Bot)
		p.CurrentItem = nil
	}
}
func (s *Bot) Pos() (*Utilities.Point) {
	return s.Position
}
func (s *Bot) OnCollision(ownTile, collTile Ptype, collPcs interface{}) {
	GH, ok := collPcs.(string)
	if ok && GH == "GH" {
		s.IsRotten = true
	}
}
func (s *Bot) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point) {
	pnt, ok := (*mapper)[*s.Position]
	if ok {
		if s.IsRotten {
			Bot_Item_Rotten_Image.X = float64(pnt.X)
			Bot_Item_Rotten_Image.Y = float64(pnt.Y)
			Bot_Item_Rotten_Image.DrawImageObj(screen)
		}else{
			Bot_Item_Image.X = float64(pnt.X)
			Bot_Item_Image.Y = float64(pnt.Y)
			Bot_Item_Image.DrawImageObj(screen)
		}
	}
}
func (s *Bot) DrawCorner(screen *ebiten.Image, pnt *Utilities.Point) {
	Bot_Item_Image_Corner.X = float64(pnt.X)
	Bot_Item_Image_Corner.Y = float64(pnt.Y)
	Bot_Item_Image_Corner.DrawImageObj(screen)
}