package Game

import (
	"github.com/hajimehoshi/ebiten"
	"math/rand"
	//"fmt"
	//"time"
)

type World struct {
	ItemFactories []ItemFactory
	Itms []Item
	SObjs []StaticObj
	Pls []*Player
	Apls []*Apple
	frame int
}
func (w *World) SpawnApple() {
	//TODO --------------------------------------------
	appl := GetRandomApple(int(XTiles),int(YTiles))
	w.Apls = append(w.Apls, appl)
}
func (w *World) RemoveApple(i int) {
	w.Apls[i] = w.Apls[len(w.Apls)-1]
    w.Apls = w.Apls[:len(w.Apls)-1]
}
func (w *World) RemoveAppleObj(a *Apple) {
	for i,apl := range(w.Apls) {
		if apl.X == a.X && apl.Y == a.Y {
			w.RemoveApple(i)
			return
		}
	}
}
func (w *World) RemoveObjs(is []int) {
	if len(is) > 0 && len(w.SObjs) > 0 {
		for c,i := range(is) {
			newVal := w.SObjs[len(w.SObjs)-1-c]
			w.SObjs[i] = newVal
		}
		w.SObjs = w.SObjs[:len(w.SObjs)-len(is)]
	}
}
func (w *World) RemoveItem(i int) {
	w.Itms[i] = w.Itms[len(w.Itms)-1]
    w.Itms = w.Itms[:len(w.Itms)-1]
}
func (w *World) SpawnItem() {
	idx := rand.Intn(len(w.ItemFactories))
	prob := rand.Float64()
	if prob < w.ItemFactories[idx].GetProbability() {
		w.Itms = append(w.Itms, w.ItemFactories[idx].GetRandomItemInBounds(int(XTiles),int(YTiles)))
	}
}

func (w *World) InitBeginning() {
	w.ItemFactories = make([]ItemFactory, 6)
	
	spF := &SpeedFactory{}; spF.Init()
	w.ItemFactories[0] = spF	
	revF := &ReviveFactory{}; revF.Init()
	w.ItemFactories[1] = revF
	lasF := &LaserFactory{}; lasF.Init()
	w.ItemFactories[2] = lasF
	bmbF := &BombFactory{}; bmbF.Init()
	w.ItemFactories[3] = bmbF
	botF := &BotFactory{}; botF.Init()
	w.ItemFactories[4] = botF
	fatF := &FartFactory{}; fatF.Init()
	w.ItemFactories[5] = fatF
}

func (w *World) Start(pls []*Player) {
	w.Apls = make([]*Apple, 0)
	for i := 0; i < int(BeginApples); i++ {
		w.SpawnApple()
	}
	w.Itms = make([]Item, 0)
	for i := 0; i < int(BeginItems); i++ {
		w.SpawnItem()
	}
	w.SObjs = make([]StaticObj,0)
	w.Pls = pls
}

func (w *World) Update(frame int) {
	w.frame = frame
	rems := make([]int, 0)
	for i,_ := range(w.SObjs) {
		w.SObjs[i].Update(SnakeProInGame)
		if !w.SObjs[i].Exists() {
			rems = append(rems, i)
		}
	}
	w.RemoveObjs(rems)
	w.CheckCollision()
	
	if frame%int(ItemSpawnPeriod) == 0 {
		w.SpawnItem()
	}
}
func (w *World) CheckCollision() {
	for i,pl := range(w.Pls) {
		switch pl.GetState() {
		case WAITING:
			
		case ALIVE:
			w.CheckAL_PL__N_APL(pl)
			w.CheckAL_PL__SELF(pl)
			w.CheckAL_PL__AL_PL(pl, i)
			w.CheckAL_PL__WALL(pl)
			w.CheckAL_PL__ANY_ITM(pl)
		case GHOST:
			w.CheckGH_PL__N_APL(pl)
			w.CheckGH_PL__ANY_ITM(pl)
			w.CheckGH_PL__WALL(pl)
		default:
			
		}
	}
}
func (w *World) CheckCollisionPlayer(pl *Player) {
	switch pl.GetState() {
		case WAITING:
		case ALIVE:
			w.CheckAL_PL__N_APL(pl)
			w.CheckAL_PL__SELF(pl)
			w.CheckAL_PL__AL_PL(pl, 0)
			w.CheckAL_PL__WALL(pl)
			w.CheckAL_PL__ANY_ITM(pl)
		case GHOST:
			w.CheckGH_PL__N_APL(pl)
			w.CheckGH_PL__ANY_ITM(pl)
			w.CheckGH_PL__WALL(pl)
		default:
	}
}

func (w *World) Draw(screen *ebiten.Image) {
	for _,obj := range(w.SObjs) {
		obj.Draw(screen, MapToScreen)
	}
	for _,itm := range(w.Itms) {
		itm.Draw(screen, MapToScreen)
	}
	for _,apl := range(w.Apls) {
		apl.Draw(screen, MapToScreen)
	}
	for _,pl := range(w.Pls) {
		pl.Draw(screen, MapToScreen)
	}
}

func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}