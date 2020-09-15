package Game

import (
	"marvin/SnakeProGo/Utilities"
)

type COLLISION int

var (
	AL_PL__N_APL = COLLISION(1)
	AL_PL__SELF = COLLISION(2)
	AL_PL__AL_PL = COLLISION(3)
	AL_PL__WALL = COLLISION(4)
	AL_PL__ANY_ITM = COLLISION(5)
	
	GH_PL__N_APL = COLLISION(7)
	GH_PL__ANY_ITM = COLLISION(8)
	GH_PL__WALL = COLLISION(9)
	
	//true if victim
	VICTORY = COLLISION(10)
	SKILL = COLLISION(11)
	LASER = COLLISION(12)
	BOMB = COLLISION(13)
	BOT = COLLISION(14)
	FART = COLLISION(15)
)

func (w *World) CheckAL_PL__N_APL(pl *Player) {
	Head := pl.GetHead()
	for i,apl := range(w.Apls) {
		if Head.CollidesWithPnt(apl.Pnt) {
			pl.EatFood(apl.NutritionValue)
			apl.OnCollision(P_Apl, P_SnkH, AL_PL__N_APL)
			pl.OnCollision(AL_PL__N_APL,0,0)
			if len(pl.GetTiles().Ts) == 0 {
				pl.SetDead(SnakeProInGame.frame)
			}
			if apl.Eaten {
				w.RemoveApple(i)
				w.SpawnApple()
				break
			}
		}
	}
}
func (w *World) CheckAL_PL__SELF(pl *Player) {
	PlTs := pl.GetAllPoints()
	if len(PlTs) > 1 {
		collPnt := pl.GetHead().CollidesWithPnts(PlTs[:len(PlTs)-1])
		if collPnt != nil {
			//pl.OnCollision(P_SnkH, P_SnkB, pl, collPnt)
			pl.OnCollision(AL_PL__SELF,0,1)
			pl.SetDead(w.frame)
		}
	}
}
func (w *World) CheckAL_PL__AL_PL(pl *Player, i int) {
	//TODO head collides with every other player or all points collide and player die if head is colliding
	if pl.GetState() == ALIVE {
	for i2 := i+1; i2 < len(w.Pls); i2++ {
		if w.Pls[i2].GetState() == ALIVE && w.Pls[i2] != pl {
			collPnts := Utilities.CollidePnts(pl.GetAllPoints(), w.Pls[i2].GetAllPoints())
			if len(collPnts) > 0 {
				p1c, p2c := false, false
				if pl.GetHead().IsContained(collPnts) {
					pl.OnCollision(AL_PL__AL_PL,0,1)
					w.Pls[i2].OnCollision(AL_PL__AL_PL,1,0)
					p1c = true
				}
				if w.Pls[i2].GetHead().IsContained(collPnts) {
					w.Pls[i2].OnCollision(AL_PL__AL_PL,0,1)
					pl.OnCollision(AL_PL__AL_PL,1,0)
					p2c = true
				}
				if p1c {
					pl.SetDead(w.frame)
				}
				if p2c {
					w.Pls[i2].SetDead(w.frame)
				}
			}
		}
	}
	}
}
func (w *World) CheckAL_PL__WALL(pl *Player) {
	if pl.GetHead().CollidesWithWall(XTiles,YTiles) {
		pl.OnCollision(AL_PL__WALL, 0, 1)
		pl.SetDead(w.frame)
	}
}
func (w *World) CheckAL_PL__ANY_ITM(pl *Player) {
	Head := pl.GetHead()
	for i,itm := range(w.Itms) {
		if Head.CollidesWithPnt(itm.Pos()) {
			SoundEffects["Item"].PR()
			pl.OnCollision(AL_PL__ANY_ITM, 0, 0)
			pl.CurrentItem = itm
			itm.OnCollision(P_Itm, P_SnkH, nil)
			w.RemoveItem(i)
			break
		}
	}
}
func (w *World) CheckGH_PL__N_APL(pl *Player) {
	Head := pl.GetHead()
	for _,apl := range(w.Apls) {
		if Head.CollidesWithPnt(apl.Pnt) {
			pl.OnCollision(GH_PL__N_APL, 0, 0)
			apl.OnCollision(P_Apl, P_SnkH, GH_PL__N_APL)
		}
	}
}
func (w *World) CheckGH_PL__ANY_ITM(pl *Player) {
	Head := pl.GetHead()
	for i,itm := range(w.Itms) {
		if Head.CollidesWithPnt(itm.Pos()) {
			rev, ok := itm.(*Revive)
			if ok {
				rev.Use(pl, SnakeProInGame)
			}
			pl.OnCollision(GH_PL__ANY_ITM, 0, 0)
			itm.OnCollision(P_Itm, P_SnkH, "GH")
			if ok {
				w.RemoveItem(i)
			}
			break
		}
	}
}
func (w *World) CheckGH_PL__WALL(pl *Player) {
	if pl.GetHead().CollidesWithWall(XTiles,YTiles) {
		pl.OnCollision(GH_PL__WALL, 0, 0)
		pl.Reset(pl.ID, pl.ConID, pl.GetSpeed(), GHOST)
	}
}