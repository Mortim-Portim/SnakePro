package Game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"marvin/SnakeProGo/Utilities"
	"image/color"
	"time"
	"fmt"
	"log"
)

func EXITGAME() {
	RuntimeStatistics.WriteToCSV(f_Resources+"/runtimeStats.csv")
	
	for _,pl := range(SnakeProInGame.Pls) {
		pl.Statistic.Set()
		pl.Statistic.Save()
		
		lastPlayerNames.Set(pl.ID, pl.Statistic.name)
	}
	
	lastPlayerNames.Save(f_Resources+f_LastPlayerNames)
	
	log.Fatal("Regular Exit")
}


type GameState interface {
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
	Init(*SnakePro)
	Start()
}

type SnakePro struct {
	CurrentState, lastState int
	
	States []GameState
	
	Cons *[]*Utilities.Controller
	
	Toasts []*Toast
	
	XRes,YRes float64
	XTiles,YTiles float64
	XTileWidth, YTileWidth float64
	firstUp, initializing bool; introObj *Utilities.ImageObj
}
func (g *SnakePro) Init() {
	g.firstUp, g.initializing = true, true;
}

func (g *SnakePro) InitBeginning(screen *ebiten.Image) {
	fmt.Println("Initializing SnakePro")
	start := time.Now()
	
	g.lastState = g.CurrentState+1
	g.XRes, g.YRes = Ps.Get("XRes"), Ps.Get("YRes")
	g.XTiles,g.YTiles = Ps.Get("XTiles"), Ps.Get("YTiles")
	g.XTileWidth, g.YTileWidth = g.XRes/g.XTiles,g.YRes/g.YTiles
		
	fmt.Println("Initializing GameStates")
	men := &Menu{};men.Init(g);g.States = append(g.States, men)
	inG := &InGame{};inG.Init(g);g.States = append(g.States, inG);SnakeProInGame = inG
	hof := &HallOfFame{};hof.Init(g);g.States = append(g.States, hof)
		
	fmt.Println("Initializing Controlls")
	err, cons := Utilities.InitControlls(ControllerThresh, PlayerNum)
	fmt.Println("Server Running on:", Utilities.SmartServIPAddr,":",Utilities.PORT)
	Utilities.CheckErr(err)
	for i,_ := range(*cons) {
		(*cons)[i].LoadConfig(f_Resources+f_Controller)
	}
	g.Cons = cons
	for _,con := range(*g.Cons) {
		con.UseCurrentAxisValsAsStandard()
	}
	
	g.Toasts = make([]*Toast, 0)
	
	elapsed := time.Since(start)
    fmt.Printf("Initializing SnakePro took %s\n", elapsed)
    if int(elapsed) < 2000000000 {
	    time.Sleep(time.Duration(2000000000-int(elapsed)))
    }
    g.initializing = false
}
func (g *SnakePro) MakeToast(str string) {
	g.Toasts = append(g.Toasts, GetStandardToastShort(str))
}
func (g *SnakePro) UpdateToasts() {
	for i,_ := range(g.Toasts) {
		g.Toasts[i].RemainingTime --
	}
}
func (g *SnakePro) RemoveToasts() {
	remList := make([]int, 0)
	for i,t := range(g.Toasts) {
		if t.RemainingTime <= 0 {
			remList = append(remList, i)
		}
	}
	for i,rIdx := range(remList) {
		g.Toasts[rIdx] = g.Toasts[len(g.Toasts)-1-i]
	}
	g.Toasts = g.Toasts[:len(g.Toasts)-len(remList)]
}
func (g *SnakePro) DrawToasts(screen *ebiten.Image) {
	for i,t := range(g.Toasts) {
		t.Draw(screen, i)
	}
}

var CornerPoints []*Utilities.Point
var LoadingSnakes []*Snake
var LoadingPoints []*Utilities.Point
var UpButtonsTime, UpListenerTime []int
func (g *SnakePro) Update(screen *ebiten.Image) error {
	start := time.Now()
	
	if g.firstUp {
		go g.InitBeginning(screen)
		g.firstUp = false
		LoadingSnakes = make([]*Snake, 4)
		LoadingPoints = make([]*Utilities.Point, 0) //14*14
		XP, YP, W, H := int(XTiles/2-4), int(YTiles/3*2), 8, 8
		CornerPoints = []*Utilities.Point{{XP,YP},
											{XP+W,YP},
											{XP+W,YP+H},
											{XP,YP+W},
											{XP,YP}}
		for i,_ := range(LoadingSnakes) {
			LoadingSnakes[i] = GetSnakeOfType(i)
			pnt := CornerPoints[i]
			LoadingSnakes[i].Reset(FPS, 0, pnt.X, pnt.Y, 7, DetDirFrom2Ts(CornerPoints[i],CornerPoints[i+1]))
		}
		
		return nil
	}
	if g.initializing {
		screen.DrawImage(backGroundImg, nil)
		
		str := "Loading..."
		Utilities.MakePopUpAt(screen, str, int(XRes/2), int(YRes/8), Utilities.GetSizeForText(str, int(XRes/2), int(YRes/4)), color.RGBA{155,155,155,255}, color.RGBA{0,0,0,0})
		Loading_Intro.DrawImageObj(screen)
		
		for _,snk := range(LoadingSnakes) {
			snk.Update(1)
			if snk.GetHead().IsContained(CornerPoints) {
				snk.nextDir = snk.currentDir.Right()
			}
			snk.Draw(screen, MapToScreen)
		}
		
		
		
		return nil
	}
	Utilities.UpdateControlls()	
	RuntimeStatistics.SetHead("UpConsGen",int(time.Since(start)))
	
	UpButtonsTime, UpListenerTime = make([]int,0),make([]int,0)
	for i,_ := range(*g.Cons) {
		bt, lt := (*g.Cons)[i].UpdateOnlyNeeded()
		UpButtonsTime = append(UpButtonsTime, bt)
		UpListenerTime = append(UpListenerTime, lt)
	}	
	RuntimeStatistics.SetHead("UpButtons1",UpButtonsTime[0])
	RuntimeStatistics.SetHead("UpButtons2",UpButtonsTime[1])
	RuntimeStatistics.SetHead("UpButtons3",UpButtonsTime[2])
	RuntimeStatistics.SetHead("UpListener1",UpListenerTime[0])
	RuntimeStatistics.SetHead("UpListener2",UpListenerTime[1])
	RuntimeStatistics.SetHead("UpListener3",UpListenerTime[2])
	
	if g.CurrentState != g.lastState {
		fmt.Println("Starting GameState", g.CurrentState)
		g.States[g.CurrentState].Start()
		g.lastState = g.CurrentState
	}
	possErr := g.States[g.CurrentState].Update(screen)
	
	RuntimeStatistics.SetHead("Update",int(time.Since(start)))
	RuntimeStatistics.SetHead("FPS",int(ebiten.CurrentFPS()))
	RuntimeStatistics.SetHead("TPS",int(ebiten.CurrentTPS()))
	RuntimeStatistics.Append()
	
	g.UpdateToasts()
	g.RemoveToasts()
	
	g.DrawToScreen(screen)
	
    return possErr
}
func (g *SnakePro) DrawToScreen(screen *ebiten.Image) {
	start := time.Now()
	g.States[g.CurrentState].Draw(screen)
	
	//W,H := screen.Size()
	//ONLY FOR DEBUGGING <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	//msg := fmt.Sprintf(`FPS: %0.2f, TPS: %0.2f, W: %v, H: %v`, ebiten.CurrentFPS(), ebiten.CurrentTPS(), W, H)
	msg := fmt.Sprintf(`TPS: %0.2f`, ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
	//ONLY FOR DEBUGGING <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	RuntimeStatistics.SetHead("Draw",int(time.Since(start)))
	g.DrawToasts(screen)
}
//func (g *SnakePro) Draw(screen *ebiten.Image) {}
func (g *SnakePro) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(XRes), int(YRes)
}
func (g *SnakePro) SaveCons(ids []int) {
	path := f_Resources+f_Controller
	for _,i := range(ids) {
		(*g.Cons)[i].SaveConfig(path)
	}
}
