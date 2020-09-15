package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"sort"
	"fmt"
)
type Point struct {
	X,Y float64
}
var (
	pokals = []string{"Wins", "Kills", "Bomb_K", "Laser_K", "Stupid_D", "Bot_K", "Skill_K", "Fart_K"}
	//HeadScale float64
	//HeadPos []*Point
	PokalScale float64
	PokalPos []*Point
	
	barChart *Utilities.ImageObj
)

func GetWinnerSnakes(first, second, third, x, y int) ([]*Snake) {
	S1 := GetSnakeOfPolygon(first, &Utilities.Point{x,y},&Utilities.Point{x,y-5})
	S2 := GetSnakeOfPolygon(second, &Utilities.Point{x-1,y},&Utilities.Point{x-1,y-4})
	S3 := GetSnakeOfPolygon(third, &Utilities.Point{x+1,y},&Utilities.Point{x+1,y-3})
	
	return []*Snake{S1, S2, S3}
}

type HallOfFame struct {
	Back *Utilities.ImageObj
	
	PokalBtns []*Utilities.Button
	PokalSnakes [][]*Snake
	//PlayerWinningHeads []*Utilities.ImageObj
	
	parentG *SnakePro
}
func GetPlayerStats(g *InGame, category string) []float64 {
	stat := make([]float64, len(g.Pls))
	for i,_ := range(g.Pls) {
		stat[i] = float64(g.Pls[i].Statistic.Statistics.Head[category])
	}
	return stat
}
func OnPokalClick(b *Utilities.Button) {
	if b.LPressed == false {	
		ID := b.Data.(int)
		if barChart == nil {
			barChart = Utilities.GetBarChart(XRes/2, YRes/2, XRes/4, YRes/4, color.RGBA{60,60,60,100}, SnakeColors, GetPlayerStats(SnakeProInGame, pokals[ID])...)
		}else{
			barChart = nil
		}
	}
}
/**
func DisplayHeadAtPokal(idx int, Head *Utilities.ImageObj, screen *ebiten.Image) {
	Head.X, Head.Y = HeadPos[idx].X, HeadPos[idx].Y
	Head.W, Head.H = HeadScale,HeadScale
	Head.Angle = 0
	Head.DrawImageObj(screen)
}
**/
func (h *HallOfFame) Init(g *SnakePro) {
	h.parentG = g
	UI_Path := f_Resources+res_dir+f_Images+f_UI
	h.Back = Utilities.LoadImgObj(UI_Path+f_HOF_Back, XRes, YRes, 0, 0, 0)
	
	PokalPos = make([]*Point, 8)
	PokalScale := XRes/8; distance := XRes/12
	btnX := XRes/8
	btnY := YRes*0.3
	
	h.PokalBtns = make([]*Utilities.Button, 8)
	for i,_ := range(h.PokalBtns) {
		xpos, ypos := btnX+float64(i)*(PokalScale+distance), btnY
		if i > 3 {
			ypos = btnY+PokalScale+distance*0.5
			xpos = btnX+float64(i-4)*(PokalScale+distance)
		}
		PokalPos[i] = &Point{xpos,ypos}
	}
	for i,pnt := range(PokalPos) {
		img := Utilities.LoadImgObj(UI_Path+fmt.Sprintf("%s%v.png", f_Pokals, i+1), PokalScale, PokalScale, pnt.X, pnt.Y, 0)
		btn := Utilities.GetButton(img, OnPokalClick, nil)
		btn.Data = i
		h.PokalBtns[i] = btn
	}
	/**
	HeadScale = XRes/16
	HeadPos = make([]*Point, 8)
	for i,pnt := range(PokalPos) {
		HeadPos[i] = &Point{pnt.X+PokalScale/2-HeadScale/2, pnt.Y-HeadScale*0.67}
	}**/
}
func getBestPlayerIdx(g *InGame, category string) int {
	highScore := 0
	bestIdx := -1
	for i,pl := range(g.Pls) {
		sc := pl.Statistic.Statistics.Head[category]
		if sc > highScore {
			highScore = sc
			bestIdx = i
		}
	}
	return bestIdx
}
type PlayerScoreAndID struct {ID, score int}
func GetStatisticsList(g *InGame, category string) []PlayerScoreAndID {
	stats := make([]PlayerScoreAndID, len(g.Pls))
	for i,pl := range(g.Pls) {
		stats[i] = PlayerScoreAndID{pl.ID,pl.Statistic.Statistics.Head[category]}
	}
	return stats
}
func GetBestOfStatistic(g *InGame, category string) []int {
	stats := GetStatisticsList(g, category)
	
	sort.SliceStable(stats, func(i, j int) bool { return stats[i].score > stats[j].score })
	if stats[0].score > 0 {
		IDs := make([]int, len(stats))	
		for i,_ := range(IDs) {
			IDs[i] = stats[i].ID
		}	
		return IDs
	}
	return nil
}
func (h *HallOfFame) Start(){
	/**
	h.PlayerWinningHeads = make([]*Utilities.ImageObj, 8)
	for i,_ := range(h.PlayerWinningHeads) {
		idx := getBestPlayerIdx(SnakeProInGame, pokals[i])
		h.PlayerWinningHeads[i] = nil
		if idx >= 0 {
			h.PlayerWinningHeads[i] = PlayerSnakeBodies[idx][0][0].Copy()
		}
	}
	**/
	h.PokalSnakes = make([][]*Snake, 8)
	for i,_ := range(h.PokalSnakes) {
		img := h.PokalBtns[i].Img
		xpos, ypos := (img.X+img.W/2)/TileWidth*SubImageScale, (img.Y+img.H/2)/TileHeight*SubImageScale-1
		best := GetBestOfStatistic(SnakeProInGame, pokals[i])
		if best != nil {
			h.PokalSnakes[i] = GetWinnerSnakes(best[0],best[1],best[2], int(xpos), int(ypos))
		}
	}
}
func (h *HallOfFame) Update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {h.parentG.CurrentState = 0}
	for _,b := range(h.PokalBtns) {
		b.Update()
	}
	return nil
}
func (h *HallOfFame) Draw(screen *ebiten.Image) {
	h.Back.DrawImageObj(screen)
	
	/**
	for i,hd := range(h.PlayerWinningHeads) {
		if hd != nil {
			DisplayHeadAtPokal(i, hd, screen)
		}
	}
	**/
	
	for _,snks := range(h.PokalSnakes) {
		for _,snk := range(snks) {
			if snk != nil {
				snk.Draw(screen, MapToScreen)
			}
		}
	}
	
	for _,b := range(h.PokalBtns) {
		b.Draw(screen)
	}
	
	if barChart != nil {
		barChart.DrawImageObj(screen)
	}
}
func (h *HallOfFame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(h.parentG.XRes), int(h.parentG.YRes)
}