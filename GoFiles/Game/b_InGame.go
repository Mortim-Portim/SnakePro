package Game

import (
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	//"marvin/SnakeProGo/Utilities"
	"time"
	"fmt"
)

type InGame struct {
	World *World
	
	Pls []*Player
	
	frame int
	ReturnTime int
	
	parentG *SnakePro
}
func (g *InGame) Init(parentG *SnakePro) {
	g.parentG = parentG
	
	InitAppleImgs()
	
	g.Pls = make([]*Player, PlayerNum)
	for i,_ := range(g.Pls) {
		g.Pls[i] = &Player{}
		g.Pls[i].Init(g, i, PlayerSpeed)
		g.Pls[i].Statistic.Load()
	}
	
	g.World = &World{}
	g.World.InitBeginning()
	
	fmt.Println("InGame initialized")
}
func (g *InGame) Start() {
	g.frame = 0
	g.ReturnTime = -1
	for i,_ := range(g.Pls) {
		conId := ControllerIds[i]
		g.Pls[i].Reset(i, conId, PlayerSpeed, ALIVE)
		g.Pls[i].SetUndestroyable(int(FPS*ReviveImmortalTime))
		(*g.parentG.Cons)[conId].RegisterDirectionEventListener(&PlayerDirectionEventListener{0,1,2,3,g.Pls[i]})
		(*g.parentG.Cons)[conId].RegisterButtonEventListener(&PlayerButtonEventListener{4,g.Pls[i],false})
	}
	
	g.World.Start(g.Pls)
	
	fmt.Println("InGame started")
}

func (g *InGame) Update(screen *ebiten.Image) error {
	start2 := time.Now()
	for i,_ := range(g.Pls) {
		g.Pls[i].Update(g.frame)
	}
	RuntimeStatistics.SetHead("PlayUp",int(time.Since(start2)))
	
	start3 := time.Now()
	g.World.Update(g.frame)
	g.CheckIfAllDead()
	RuntimeStatistics.SetHead("WorldUp",int(time.Since(start3)))
	
	g.frame ++
	return nil
}
func (g *InGame) Draw(screen *ebiten.Image) {
	screen.DrawImage(backGroundImg, nil)
	if g.World != nil {
		g.World.Draw(screen)
	}
}
func (g *InGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.parentG.XRes), int(g.parentG.YRes)
}

func (g *InGame) CheckIfAllDead() {
	if g.ReturnTime < 0 {
		alivePls := make([]int, 0)
		for i,_ := range(g.Pls) {
			if g.Pls[i].GetState() != GHOST {
				alivePls = append(alivePls, i)
			}
		}
		if len(alivePls) <= 1 {
			g.ReturnTime = g.frame+PlayerAliveAloneFrames
		}
		if len(alivePls) == 1 {
			g.Pls[alivePls[0]].OnCollision(VICTORY,0,0)
		}
	}else if g.ReturnTime == g.frame {
		g.parentG.CurrentState = 0
		SoundEffects["Win"].PR()
		fmt.Println("AllSnakes Dead!!")
	}
}