package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"math"
	//"fmt"
)

func GetPolynomOfPoints(pnts ...*Utilities.Point) []*Utilities.Point {
	allPnts := make([]*Utilities.Point,0)
	allPnts = append(allPnts, pnts[0])
	for i,pnt := range(pnts) {
		if i > 0 {
			line := GetLineOfPoints(pnts[i-1], pnt)
			for pI,p := range(line) {
				if pI > 0 {
					allPnts = append(allPnts, p)
				}
			}
		}
	}
	return allPnts
}

func GetLineOfPoints(pnt1, pnt2 *Utilities.Point) []*Utilities.Point {
	currentPoint := pnt1.Copy()
	line := make([]*Utilities.Point, 1)
	line[0] = currentPoint
	for !currentPoint.Equals(pnt2) {
		xdif, ydif := DetDirFrom2Ts(currentPoint, pnt2).GetXYDif()
		currentPoint = &Utilities.Point{currentPoint.X+int(xdif), currentPoint.Y+int(ydif)}
		line = append(line, currentPoint)
	}
	return line
}

func GetSnakeOfPolygon(ID int, pnts ...*Utilities.Point) *Snake {
	s := GetSnakeOfType(ID)
	s.tiles = &Utilities.SnakeTiles{}
	s.tiles.Ts = GetPolynomOfPoints(pnts...)
	s.currentDir = DetDirFrom2Ts(pnts[len(pnts)-2],pnts[len(pnts)-1])
	s.compAllPieces()
	return s
}

func GetSnakeOfType(playerID int) *Snake {
	SnkH,SnkB,SnkR,SnkL,SnkT := GetSnakeTilesFromPlayerSnakeBodies(playerID)
	s := CreateSnake(SnakeProInGame, SnkH,SnkB,SnkT,SnkR,SnkL)
	return s
}

type SnakePiece struct {
	Type Ptype
	dir *Direction
	parent *Snake
	img *Utilities.ImageObj
}
func (s *SnakePiece) Ptype() (Ptype) {
	return s.Type
}
func (s *SnakePiece) GetImg() *Utilities.ImageObj {
	s.img.Angle = s.dir.GetRotation()
	return s.img
}

type Snake struct {
	tiles *Utilities.SnakeTiles //Back == Kopf
	currentPieces []Piece
	currentDir, nextDir *Direction
	undestroyable int
	currentFett, Score, state int
	speed float64
	SnkH,SnkB,SnkT,SnkR,SnkL []*Utilities.ImageObj
	parentGame *InGame
}

func CreateSnake(g *InGame, SnkH,SnkB,SnkT,SnkR,SnkL []*Utilities.ImageObj) (s *Snake) {
	s = &Snake{}
	s.SnkH,s.SnkB,s.SnkT,s.SnkR,s.SnkL = SnkH,SnkB,SnkT,SnkR,SnkL
	s.parentGame = g
	return
}

func (s *Snake) Reset(speed float64, state, X, Y, lenght int, dir *Direction) {
	s.currentDir = dir; s.nextDir = nil;
	s.currentFett = lenght
	s.tiles = &Utilities.SnakeTiles{}
	s.tiles.Ts = make([]*Utilities.Point, 1)
	s.tiles.Ts[0] = &Utilities.Point{int(X),int(Y)}
	s.speed = speed; s.state = state
	s.undestroyable = 0
	
	s.compAllPieces()
}
func (s *Snake) compAllPieces() {
	Pcs := make([]Piece, len(s.tiles.Ts))
	Pcs[len(Pcs)-1] = &SnakePiece{P_SnkH, s.currentDir.Copy(), s, s.SnkH[int(s.state)]}
	for i := len(s.tiles.Ts)-2; i >= 0; i-- {
		var piece *SnakePiece
		if i > 0 {
			tL := s.tiles.Ts[i+1]; tN := s.tiles.Ts[i-1];t := s.tiles.Ts[i]
			TtoTN := DetDirFrom2Ts(t, tN)
			TtoTL := DetDirFrom2Ts(t, tL)
			var t_img *Utilities.ImageObj
			var t_dir *Direction
			if TtoTN.IsLeft(TtoTL) {
				t_img = s.SnkR[int(s.state)]
				t_dir = DetDirFrom2Ts(t, tL)
			}else if TtoTN.IsRight(TtoTL) {
				t_img = s.SnkL[int(s.state)]
				t_dir = DetDirFrom2Ts(tL, t)
			}else if TtoTN.IsInvers(TtoTL) || TtoTN.Equals(TtoTL) {
				t_img = s.SnkB[int(s.state)]
				t_dir = DetDirFrom2Ts(t, tL).Left()
			}
			piece = &SnakePiece{P_SnkB, t_dir, s, t_img}
		}else {
			t_dir := DetDirFrom2Ts(s.tiles.Ts[0], s.tiles.Ts[1]).Left()
			piece = &SnakePiece{P_SnkT, t_dir, s, s.SnkT[int(s.state)]}
		}
		Pcs[i] = piece
	}
	s.currentPieces = Pcs
}

func (s *Snake) Update(frame int) {
	if frame%int(FPS/s.speed) == 0{
		s.UpdateCurrDir()
		s.UpdateSnakeHead()
		s.UpdateSnakeTail()
		s.compAllPieces()
	}
}
func (s *Snake) UpdateCurrDir() {
	if s.nextDir != nil {
		s.currentDir = s.nextDir
		s.nextDir = nil
	}
}
func (s *Snake) UpdateSnakeHead() {
	xdif, ydif := s.currentDir.GetXYDif()
	s.tiles.AddBack(int(xdif), int(ydif), 1)
}
func (s *Snake) UpdateSnakeTail() {
	if s.currentFett >= 1 {
		s.currentFett --
	}else if s.currentFett == 0 {
		s.tiles.RemFront(1)
	}else if s.currentFett < 0 {
		for ;s.currentFett != 0; s.currentFett ++ {
			if len(s.tiles.Ts) > 1 {
				s.tiles.RemFront(1)
			}
		}
	}
}
func (s *Snake) GetAllPieces() ([]Piece, []*Utilities.Point) {
	return s.currentPieces, s.tiles.Ts
}
func (s *Snake) GetAllPoints() ([]*Utilities.Point) {
	return s.tiles.Ts
}
func (s *Snake) GetHead() (*Utilities.Point) {
	return s.tiles.GetLast()
}
func (s *Snake) EatFood(NutritionValue int) {
	s.currentFett += NutritionValue
}
func (s *Snake) Draw(screen *ebiten.Image, PntToScrn *map[Utilities.Point]Utilities.Point) {
	pcs, pnts := s.GetAllPieces()
	for i,pnt := range(pnts) {
		DrawPieceToScreen(pcs[i].GetImg(), screen, pnt, PntToScrn)
	}
}