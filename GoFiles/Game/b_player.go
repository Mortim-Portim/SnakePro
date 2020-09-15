package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"math"
	"fmt"
)
var (
	PlayerSnakeBodies [][][]*Utilities.ImageObj
)
func GetSnakeTilesFromPlayerSnakeBodies(ID int) ([]*Utilities.ImageObj,[]*Utilities.ImageObj,[]*Utilities.ImageObj,[]*Utilities.ImageObj,[]*Utilities.ImageObj) {
	return PlayerSnakeBodies[ID][0],PlayerSnakeBodies[ID][1],PlayerSnakeBodies[ID][2],PlayerSnakeBodies[ID][3],PlayerSnakeBodies[ID][4]
}
type PlayerDirectionEventListener struct {
	x1,x2,y1,y2 int
	p *Player
}
func (l *PlayerDirectionEventListener) GetAxis() (int,int,int,int) {
	return l.x1,l.x2,l.y1,l.y2
}
func (l *PlayerDirectionEventListener) OnDirectionEvent(xdif, ydif float64) {
	conDir := GetDirFromXYDif(xdif, ydif)
	if !conDir.Equals(l.p.snake.currentDir) && !conDir.Inverse().Equals(l.p.snake.currentDir) {
		l.p.SetNextDir(conDir)
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
	WAITING = PlayerState(2)
)

type Player struct {
	snake *Snake
	ID, ConID int
	CurrentItem Item
	ItemPos *Utilities.Point
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
func GetPlayerTiles(ID int) [][]*Utilities.ImageObj {
	SnkH,SnkB,SnkT,SnkR,SnkL := make([]*Utilities.ImageObj, 3),make([]*Utilities.ImageObj, 3),make([]*Utilities.ImageObj, 3),make([]*Utilities.ImageObj, 3),make([]*Utilities.ImageObj, 3)
	SnkH_path := f_Resources+res_dir+f_Images+f_Snks+f_Player+fmt.Sprintf("%v",ID)
	SnkH[0] = Utilities.LoadImgObj(SnkH_path+f_SnkH+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	SnkB[0] = Utilities.LoadImgObj(SnkH_path+f_SnkB+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	SnkT[0] = Utilities.LoadImgObj(SnkH_path+f_SnkT+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	SnkL[0] = Utilities.LoadImgObj(SnkH_path+f_SnkL+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	SnkR[0] = Utilities.LoadImgObj(SnkH_path+f_SnkR+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	SnkH[1] = Utilities.LoadImgObj(SnkH_path+f_SnkH+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	SnkB[1] = Utilities.LoadImgObj(SnkH_path+f_SnkB+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	SnkT[1] = Utilities.LoadImgObj(SnkH_path+f_SnkT+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	SnkL[1] = Utilities.LoadImgObj(SnkH_path+f_SnkL+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	SnkR[1] = Utilities.LoadImgObj(SnkH_path+f_SnkR+f_GHOST+f_snkFormat, TileWidth, TileHeight, 0,0, 0)
	return [][]*Utilities.ImageObj{SnkH,SnkB,SnkR,SnkL,SnkT}
}
func (p *Player) Init(g *InGame, ID int, speed float64) {
	p.Statistic = GetNewPlayerStatistic(lastPlayerNames.Get(ID), []string{"Wins","Kills","Deaths",
			"Apples","Items",
			"Stupid_D","Skill_D","Laser_D","Bomb_D","Bot_D","Skill_K","Laser_K","Bomb_K","Bot_K","Fart_D","Fart_K"})
	p.ItemPos = GetItemPos(ID)
	SnkH,SnkB,SnkR,SnkL,SnkT := GetSnakeTilesFromPlayerSnakeBodies(ID)
	p.snake = CreateSnake(g, SnkH,SnkB,SnkT,SnkR,SnkL)
	p.Reset(ID,ID,speed,ALIVE)
}
func (p *Player) Reset(ID, conID int, speed float64, state PlayerState) {
	X,Y, startDir := Utilities.GetStartPos4P(ID, XTiles-1, YTiles-1, 0)
	p.CurrentItem = nil
	p.ID = ID; p.ConID = conID
	p.snake.Reset(speed, int(state), int(X), int(Y), PlayerStartLenght, &Direction{startDir})
}
func (p *Player) Update(frame int) {
	if p.GetState() == ALIVE {
		p.SetScore(frame)
	}
	if p.GetState() != WAITING {
		p.snake.Update(frame)
	}
}
func (p *Player) SetUndestroyable(t int) {
	p.snake.undestroyable = t
}
func (p *Player) GetTiles() *Utilities.SnakeTiles {
	return p.snake.tiles
}
func (p *Player) GetCurrDir() *Direction {
	return p.snake.currentDir
}
func (p *Player) SetCurrDir(dir *Direction) {
	p.snake.currentDir = dir
}
func (p *Player) GetNextDir() *Direction {
	return p.snake.nextDir
}
func (p *Player) SetNextDir(dir *Direction) {
	p.snake.nextDir = dir
}
func (p *Player) GetState() PlayerState {
	return PlayerState(p.snake.state)
}
func (p *Player) SetState(state PlayerState) {
	p.snake.state = int(state)
}
func (p *Player) GetScore() int {
	return p.snake.Score
}
func (p *Player) SetScore(score int) {
	p.snake.Score = score
}
func (p *Player) GetSpeed() float64 {
	return p.snake.speed
}
func (p *Player) SetSpeed(speed float64) {
	p.snake.speed = speed
}
func (p *Player) GetAllPieces() ([]Piece, []*Utilities.Point) {
	return p.snake.currentPieces, p.snake.tiles.Ts
}
func (p *Player) GetAllPoints() ([]*Utilities.Point) {
	return p.snake.tiles.Ts
}
func (p *Player) GetHead() (*Utilities.Point) {
	return p.snake.tiles.GetLast()
}
func (p *Player) GetTail() (*Utilities.Point) {
	return p.snake.tiles.GetFirst()
}
func (p *Player) EatFood(NutritionValue int) {
	p.snake.currentFett += NutritionValue
}
func (p *Player) Draw(screen *ebiten.Image, PntToScrn *map[Utilities.Point]Utilities.Point) {
	if p.GetState() != WAITING {
		pcs, pnts := p.GetAllPieces()
		for i,pnt := range(pnts) {
			DrawPieceToScreen(pcs[i].GetImg(), screen, pnt, PntToScrn)
		}
		if p.CurrentItem != nil {
			pnt := (*PntToScrn)[*p.ItemPos]
			p.CurrentItem.DrawCorner(screen, &pnt)
		}
	}
}

func (p *Player) SetDead(frame int) {
	if p.snake.undestroyable <= frame {
		if p.GetState() == ALIVE {
			SoundEffects["Dead"].PR()
		}
		oldItm := p.CurrentItem
		p.Reset(p.ID, p.ConID, p.GetSpeed(), GHOST)
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
		case FART:
			p.Statistic.AddHead("Fart_D", dead)
			p.Statistic.AddHead("Fart_K", kill)
		case GH_PL__N_APL:
		
		case GH_PL__ANY_ITM:
		
		case GH_PL__WALL:
			
	}
	p.Statistic.Set()
}