package Old
/**
package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"math"
	"fmt"
)
var (
	FullScalePlayerHeads []*Utilities.ImageObj
)
type PlayerDirectionEventListener struct {
	x1,x2,y1,y2 int
	p *Player
}
func (l *PlayerDirectionEventListener) GetAxis() (int,int,int,int) {
	return l.x1,l.x2,l.y1,l.y2
}
func (l *PlayerDirectionEventListener) OnDirectionEvent(xdif, ydif float64) {
	conDir := GetDirFromXYDif(xdif, ydif)
	if !conDir.Equals(l.p.currentDir) && !conDir.Inverse().Equals(l.p.currentDir) {
		l.p.nextDir = conDir
	}
}
type PlayerButtonEventListener struct {
	btn int
	p *Player
	state bool
}
func (l *PlayerButtonEventListener) GetButton() (int) {
	return l.btn
}
func (l *PlayerButtonEventListener) OnButtonDown() {
	if !l.state {
		if l.p.CurrentItem != nil {
			l.p.CurrentItem.Use(l.p, SnakeProInGame)
		}
	}
	l.state = true
}
func (l *PlayerButtonEventListener) OnButtonUp() {
	l.state = false
}

type PlayerState int
const (
	ALIVE = PlayerState(0)
	GHOST = PlayerState(1)
	DEAD = PlayerState(2)
)

type PlayerPiece struct {
	Type Ptype
	dir *Direction
	parent *Player
	img *Utilities.ImageObj
}
func (p *PlayerPiece) Ptype() (Ptype) {
	return p.Type
}
func (p *PlayerPiece) GetImg() *Utilities.ImageObj {
	p.img.Angle = p.dir.GetRotation()
	return p.img
}


type Player struct {
	ID, ConID int
	tiles Utilities.SnakeTiles //Back == Kopf
	currentPieces []Piece
	CurrentItem Item
	ItemPos *Utilities.Point
	currentDir, nextDir *Direction
	undestroyable int
	
	currentFett, Score int
	speed float64
	SnkH,SnkB,SnkT,SnkR,SnkL []*Utilities.ImageObj
	
	state PlayerState
	parentGame *InGame
	
	Statistic *PlayerStats
}
func GetItemPos(ID int) (pnt *Utilities.Point) {
	pnt = &Utilities.Point{}
	switch ID {
		case 0:
			pnt.X = 1; pnt.Y = 1
		case 1:
			pnt.X = int(XTiles)-2; pnt.Y = 1
		case 2:
			pnt.X = int(XTiles)-2; pnt.Y = int(YTiles)-2
		case 3:
			pnt.X = 1; pnt.Y = int(YTiles)-2
	}
	return
}

func (p *Player) Init(g *InGame, ID int, speed float64) {
	p.SnkH,p.SnkB,p.SnkT,p.SnkR,p.SnkL = make([]*Utilities.ImageObj, 3),make([]*Utilities.ImageObj, 3),make([]*Utilities.ImageObj, 3),make([]*Utilities.ImageObj, 3),make([]*Utilities.ImageObj, 3)
	p.parentGame = g
	p.Reset(ID, ID, speed, ALIVE)
	p.InitPlayerTileImgs()
	p.Statistic = GetNewPlayerStatistic(lastPlayerNames.Get(ID), []string{"Wins","Kills","Deaths",
			"Apples","Items",
			"Stupid_D","Skill_D","Laser_D","Bomb_D","Bot_D","Skill_K","Laser_K","Bomb_K","Bot_K"})
	p.ItemPos = GetItemPos(ID)
}
func (p *Player) InitPlayerTileImgs() {
	SnkH_path := f_Resources+res_dir+f_Images+f_Snks+f_Player+fmt.Sprintf("%v",p.ID)
	p.SnkH[0] = Utilities.LoadImgObj(SnkH_path+f_SnkH+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	p.SnkB[0] = Utilities.LoadImgObj(SnkH_path+f_SnkB+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	p.SnkT[0] = Utilities.LoadImgObj(SnkH_path+f_SnkT+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	p.SnkL[0] = Utilities.LoadImgObj(SnkH_path+f_SnkL+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	p.SnkR[0] = Utilities.LoadImgObj(SnkH_path+f_SnkR+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	p.SnkH[1] = Utilities.LoadImgObj(SnkH_path+f_SnkH+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	p.SnkB[1] = Utilities.LoadImgObj(SnkH_path+f_SnkB+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	p.SnkT[1] = Utilities.LoadImgObj(SnkH_path+f_SnkT+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	p.SnkL[1] = Utilities.LoadImgObj(SnkH_path+f_SnkL+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	p.SnkR[1] = Utilities.LoadImgObj(SnkH_path+f_SnkR+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
}
func (p *Player) Reset(ID, conID int, speed float64, state PlayerState) {
	X,Y, startDir := Utilities.GetStartPos4P(ID, XTiles-1, YTiles-1, 0)
	p.currentDir = &Direction{startDir}; p.nextDir = nil;
	p.CurrentItem = nil;p.currentFett = PlayerStartLenght
	p.tiles.Ts = make([]*Utilities.Point, 1)
	p.tiles.Ts[0] = &Utilities.Point{int(X),int(Y)}
	p.speed = speed; p.state = state; p.ID = ID; p.ConID = conID
	p.undestroyable = 0
	
	p.compAllPieces()
}
func (p *Player) Update(frame int) {
	if p.state == ALIVE {
		p.Score = frame
	}
	if frame%int(FPS/p.speed) == 0{
		p.UpdateCurrDir()
		p.UpdateSnakeHead()
		p.UpdateSnakeTail()
		p.compAllPieces()
	}
}
func (p *Player) UpdateCurrDir() {
	if p.nextDir != nil {
		p.currentDir = p.nextDir
		p.nextDir = nil
	}
}
func (p *Player) UpdateSnakeHead() {
	xdif, ydif := p.currentDir.GetXYDif()
	p.tiles.AddBack(int(xdif), int(ydif), 1)
}
func (p *Player) UpdateSnakeTail() {
	if p.currentFett >= 1 {
		p.currentFett --
	}else if p.currentFett == 0 {
		p.tiles.RemFront(1)
	}else if p.currentFett < 0 {
		for ;p.currentFett != 0; p.currentFett ++ {
			if len(p.tiles.Ts) > 1 {
				p.tiles.RemFront(1)
			}
		}
	}
}
func (p *Player) compAllPieces() {
	Pcs := make([]Piece, len(p.tiles.Ts))
	Pcs[len(Pcs)-1] = &PlayerPiece{P_SnkH, p.currentDir.Copy(), p, p.SnkH[int(p.state)]}
	for i := len(p.tiles.Ts)-2; i >= 0; i-- {
		var piece *PlayerPiece
		if i > 0 {
			tL := p.tiles.Ts[i+1]; tN := p.tiles.Ts[i-1];t := p.tiles.Ts[i]
			TtoTN := DetDirFrom2Ts(t, tN)
			TtoTL := DetDirFrom2Ts(t, tL)
			var t_img *Utilities.ImageObj
			var t_dir *Direction
			if TtoTN.IsLeft(TtoTL) {
				t_img = p.SnkR[int(p.state)]
				t_dir = DetDirFrom2Ts(t, tL)
			}else if TtoTN.IsRight(TtoTL) {
				t_img = p.SnkL[int(p.state)]
				t_dir = DetDirFrom2Ts(tL, t)
			}else if TtoTN.IsInvers(TtoTL) || TtoTN.Equals(TtoTL) {
				t_img = p.SnkB[int(p.state)]
				t_dir = DetDirFrom2Ts(t, tL).Left()
			}
			piece = &PlayerPiece{P_SnkB, t_dir, p, t_img}
		}else {
			t_dir := DetDirFrom2Ts(p.tiles.Ts[0], p.tiles.Ts[1]).Left()
			piece = &PlayerPiece{P_SnkT, t_dir, p, p.SnkT[int(p.state)]}
		}
		Pcs[i] = piece
	}
	p.currentPieces = Pcs
}
func (p *Player) GetAllPieces() ([]Piece, []*Utilities.Point) {
	return p.currentPieces, p.tiles.Ts
}
func (p *Player) GetAllPoints() ([]*Utilities.Point) {
	return p.tiles.Ts
}
func (p *Player) GetHead() (*Utilities.Point) {
	return p.tiles.GetLast()
}
func (p *Player) EatFood(NutritionValue int) {
	p.currentFett += NutritionValue
}
func (p *Player) Draw(screen *ebiten.Image, PntToScrn *map[Utilities.Point]Utilities.Point) {
	pcs, pnts := p.GetAllPieces()
	for i,pnt := range(pnts) {
		DrawPieceToScreen(pcs[i].GetImg(), screen, pnt, PntToScrn)
	}
	if p.CurrentItem != nil {
		pnt := (*PntToScrn)[*p.ItemPos]
		p.CurrentItem.DrawCorner(screen, &pnt)
	}
}

func (p *Player) SetDead(frame int) {
	if p.undestroyable <= frame {
		if p.state == ALIVE {
			SoundEffects["Dead"].PR()
		}
		oldItm := p.CurrentItem
		p.Reset(p.ID, p.ConID, p.speed, GHOST)
		if oldItm != nil {
			oldItm.Use(p, SnakeProInGame)
		}
	}
}
func (p *Player) OnCollision(coll COLLISION, kill, dead int) {
	p.Statistic.AddHead("Kills", kill)
	p.Statistic.AddHead("Deaths",dead)
	switch coll {
		case AL_PL__N_APL:
			p.Statistic.AddHead("Apples", 1)
		case AL_PL__SELF:
			p.Statistic.AddHead("Stupid_D", 1)
		case AL_PL__WALL:
			p.Statistic.AddHead("Stupid_D", 1)
		case AL_PL__AL_PL:
			p.Statistic.AddHead("Skill_D", dead)
			p.Statistic.AddHead("Skill_K", kill)
		case AL_PL__ANY_ITM:
			p.Statistic.AddHead("Items", 1)
		case VICTORY:
			p.Statistic.AddHead("Wins", 1)
		case SKILL:
			p.Statistic.AddHead("Skill_D", dead)
			p.Statistic.AddHead("Skill_K", kill)
		case LASER:
			p.Statistic.AddHead("Laser_D", dead)
			p.Statistic.AddHead("Laser_K", kill)
		case BOMB:
			p.Statistic.AddHead("Bomb_D", dead)
			p.Statistic.AddHead("Bomb_K", kill)
		case BOT:
			p.Statistic.AddHead("Bot_D", dead)
			p.Statistic.AddHead("Bot_K", kill)
		case GH_PL__N_APL:
		
		case GH_PL__ANY_ITM:
		
		case GH_PL__WALL:
			
	}
	p.Statistic.Set()
}
**/