package Game

import (
	"github.com/hajimehoshi/ebiten"
	"marvin/SnakeProGo/Utilities"
	"image/color"
	"fmt"
	"math"
	//"time"
)
const MaxScoreTiles = 18

var (
	ControllSettings = []string{"Left","Right","Up","Down","Item"}
	SettingControll int
	ChangingConID int
	SettingNameID int
	NEWNAMEFORID string
	PLAYPLAYPLAY bool
)

func OnPlayNameLClick(b *Utilities.Button) {
	if b.RPressed == false {
		SettingNameID = b.Data.(int)
		NEWNAMEFORID = ""
	}
}
func OnPlayNameRClick(b *Utilities.Button) {
	if b.LPressed == false {
		ChangingConID = b.Data.(int)
	}
}

func OnExitClick(b *Utilities.Button) {
	if b.LPressed == false {
		EXITGAME()
	}
}
func OnControllerClick(b *Utilities.Button) {
	if b.LPressed == false {
		SettingControll = 1
	}
}
func PlayClick(b *Utilities.Button) {
	if b.LPressed == false {
		PLAYPLAYPLAY = true
	}
}
func HOFClick(b *Utilities.Button) {
	if b.LPressed == false {
		SnakeProGame.CurrentState = 2
	}
}

type Menu struct {
	Buttons, NameBtns []*Utilities.Button
	
	BackGroundImg *Utilities.ImageObj
	
	PlayerScores []*Snake
	
	ChangingCons []int
	
	parentG *SnakePro
}
func (g *Menu) Init(parentG *SnakePro) {
	g.parentG = parentG
	
	btns_si := XRes/20
	exitImg := Utilities.LoadImgObj(f_Resources+res_dir+f_Images+f_UI+f_Quit, btns_si, btns_si, btns_si*Exit_X, btns_si*Exit_Y, 0)
	exitButton := Utilities.GetButton(exitImg, OnExitClick, nil)
	controllerImg := Utilities.LoadImgObj(f_Resources+res_dir+f_Images+f_UI+f_ControllerImg, btns_si, btns_si, btns_si*Con_X, btns_si*Con_Y, 0)
	ControllerButton := Utilities.GetButton(controllerImg, OnControllerClick, nil)
	playImg := Utilities.LoadImgObj(f_Resources+res_dir+f_Images+f_UI+f_Play, btns_si, btns_si, btns_si*Play_X, btns_si*Play_Y, 0)
	playButton := Utilities.GetButton(playImg, PlayClick, nil)
	HOFImg := Utilities.LoadImgObj(f_Resources+res_dir+f_Images+f_UI+f_HOF, btns_si, btns_si, btns_si*HOF_X, btns_si*HOF_Y, 0)
	HOFButton := Utilities.GetButton(HOFImg, HOFClick, nil)	
	
	g.Buttons = append(g.Buttons, exitButton);g.Buttons = append(g.Buttons, ControllerButton)
	g.Buttons = append(g.Buttons, playButton);g.Buttons = append(g.Buttons, HOFButton)
	
	g.BackGroundImg = Utilities.LoadImgObj(f_Resources+res_dir+f_Images+f_UI+f_MenuBack, XRes, YRes, 0,0,0)
	
	fmt.Println("Menu initialized")
}
func (g *Menu) UpdateDrawedPlayerName(i int) {
	xp, yp := PLN_XPOS*XRes, PLN_YPOS*YRes
	w, h := PLN_WIDTH*XRes, PLN_HEIGHT*YRes/float64(PlayerNum)
	str := fmt.Sprintf("%s : %v",SnakeProInGame.Pls[i].Statistic.Name(), SnakeProInGame.Pls[i].GetScore())
	btn := Utilities.GetButtonWithTextOnyl(str, xp,yp+h*float64(i),w,h, SnakeColors[i],OnPlayNameLClick, OnPlayNameRClick)
	btn.Data = i
	g.NameBtns[i] = btn
}
func (g *Menu) Start() {
	fmt.Println("Menu started")
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
	PLAYPLAYPLAY = false
	for _,con := range(*g.parentG.Cons) {
		con.ResetDirectionEventListener()
		con.ResetButtonEventListener()
	}
	g.NameBtns = make([]*Utilities.Button,PlayerNum)
	for i,_ := range(g.NameBtns) {
		g.UpdateDrawedPlayerName(i)
	}
	
	highScore := 0
	for _,pl := range(SnakeProInGame.Pls) {
		if pl.GetScore() > highScore {
			highScore = pl.GetScore()
		} 
	}
	g.PlayerScores = make([]*Snake, PlayerNum)
	for i,pl := range(SnakeProInGame.Pls) {
		score := pl.GetScore()
		X := int(MaxScoreTiles*float64(score)/float64(highScore))
		//fmt.Printf("score: %v, highScore: %v, X: %v\n", score, highScore, X)
		if X > 0 {
			g.PlayerScores[i] = GetSnakeOfPolygon(i, &Utilities.Point{0,i*2+1},&Utilities.Point{X,i*2+1})
		}
	}
	
	ChangingConID = -1
	SettingNameID = -1
}

var KeyBackspaceState bool
func (g *Menu) Update(screen *ebiten.Image) error {
	for _,b := range(g.Buttons) {
		b.Update()
	}
	for _,b := range(g.NameBtns) {
		b.Update()
	}
	if PLAYPLAYPLAY {
		g.parentG.CurrentState = 1
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	}
	if SettingControll > 0 {
		g.ChangingCons = make([]int, 0)
		g.SetController()
		return nil
	}
	if ChangingConID >= 0 {
		conId := -1
		for i,con := range(*g.parentG.Cons) {
			con.UpdateAll()
			if len(con.Down) > 0 {
				conId = i
				break
			}
		}
		if conId != -1 {
			ControllerIds[ChangingConID] = conId
			g.parentG.MakeToast(fmt.Sprintf("%s now uses %s",SnakeProInGame.Pls[ChangingConID].Statistic.Name(), (*SnakeProGame.Cons)[ControllerIds[ChangingConID]].Name))
			ChangingConID = -1
		}
	}
	if SettingNameID >= 0 {
		NEWNAMEFORID += string(ebiten.InputChars())
		SnakeProInGame.Pls[SettingNameID].Statistic.SetName(NEWNAMEFORID)
		g.UpdateDrawedPlayerName(SettingNameID)
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			plSetting := SnakeProInGame.Pls[SettingNameID]
			SnakeProGame.MakeToast(fmt.Sprintf("Name for Player %v is now %s", SettingNameID, plSetting.Statistic.Name()))
			err := plSetting.Statistic.Load()
			if err != nil {
				SnakeProGame.MakeToast(fmt.Sprintf("%v",err))
			}
			
			plSetting.Statistic.Save()
			SettingNameID = -1
		}
		
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			if KeyBackspaceState == false {
				NEWNAMEFORID = NEWNAMEFORID[:len(NEWNAMEFORID)-1]
			}
			KeyBackspaceState = true
		}else{
			KeyBackspaceState = false
		}
	}
	return nil
}
func (g *Menu) Draw(screen *ebiten.Image) {
	g.BackGroundImg.DrawImageObj(screen)
	for _,b := range(g.Buttons) {
		b.Draw(screen)
	}
	for _,b := range(g.NameBtns) {
		b.Draw(screen)
	}
	for _,snk := range(g.PlayerScores) {
		if snk != nil {
			snk.Draw(screen, MapToScreen)
		}
	}
	if SettingControll > 0 {
		Utilities.MakePopUpAt(screen,
			fmt.Sprintf("Press any Button to assign to '%s'", ControllSettings[SettingControll-1]),
			int(XRes/2), int(YRes/2), 24,
			color.RGBA{255,255,255,255}, color.RGBA{255,0,0,255})
		//Utilities.PrintTextAt(screen, fmt.Sprintf("Press any Button to assign to '%s'", ControllSettings[g.SettingControll-1]), int(XRes/2), int(YRes/2), 255,255,255, 24)
	}
	if ChangingConID >= 0 {
		Utilities.MakePopUpAt(screen,
			fmt.Sprintf("Press any Button on any Controller to assign to %s", SnakeProInGame.Pls[ChangingConID].Statistic.Name()),
			int(XRes/2), int(YRes/2), 24,
			color.RGBA{255,255,255,255}, color.RGBA{255,0,0,255})
	}
}
func (g *Menu) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.parentG.XRes), int(g.parentG.YRes)
}

func (g *Menu) SetController() {
	buttonDownId := -1
	conId := -1
	value := -1.0
	for i,_ := range(*g.parentG.Cons) {
		(*g.parentG.Cons)[i].UpdateAll()
		if len((*g.parentG.Cons)[i].Down) > 0 {
			conId = i
			buttonDownId = (*g.parentG.Cons)[i].Down[0]
			value = (*g.parentG.Cons)[i].Buttons[buttonDownId]
			break
		}
	}
	if conId != -1 && math.Abs(value) > 0.5 {
		if SettingControll >= 2 {
			if (*g.parentG.Cons)[conId].Mapper[SettingControll-2] == buttonDownId {
				return
			}
		}
		(*g.parentG.Cons)[conId].SetMapper(SettingControll-1, buttonDownId)
		//time.Sleep(1*time.Second)
		SettingControll ++
		g.ChangingCons = append(g.ChangingCons, conId)
	}
	if SettingControll > 5 {
		g.parentG.SaveCons(g.ChangingCons)
		SettingControll = 0
		g.parentG.MakeToast("Controlls set")
	}
}