package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

var (
	AppleImg, AppleImg_Rotten *Utilities.ImageObj
)

func InitAppleImgs() {
	AppleImg = Utilities.LoadImgObj(f_Resources+res_dir+f_Images+f_Food+f_Apple, TileWidth, TileHeight, 0, 0, 0)
	AppleImg_Rotten = Utilities.LoadImgObj(f_Resources+res_dir+f_Images+f_Food+f_Apple_Rotten, TileWidth, TileHeight, 0, 0, 0)
}

type AppleTimer struct {
	Timer int
	Apple *Apple
}
func (o *AppleTimer) Update(g *InGame) {
	o.Timer --
	if o.Timer == 0 {
		o.Apple.Eaten = true
		g.World.RemoveAppleObj(o.Apple)
		g.World.SpawnApple()
	}
}
func (o *AppleTimer) Exists() bool {
	if o.Timer <= 0 {
		return false
	}
	return true
}
func (o *AppleTimer) Draw(screen *ebiten.Image, mapper *map[Utilities.Point]Utilities.Point){}

type Apple struct {
	X,Y int
	NutritionValue int
	Eaten bool
	Pnt *Utilities.Point
}
func GetRandomApple(W,H int) *Apple {
	XP,YP := Utilities.GetRandomPosInScreen(W,H)
	a := Apple{XP,YP,int(AppleNutrition),false,&Utilities.Point{0,0}}
	a.Pnt = &Utilities.Point{a.X, a.Y}
	return &a
}
func (a *Apple) Draw(screen *ebiten.Image, PntToScrn *map[Utilities.Point]Utilities.Point) {
	if a.NutritionValue < 0 {
		DrawPieceToScreen(AppleImg_Rotten, screen, a.Pnt, PntToScrn)
		
		
		
	}else{
		DrawPieceToScreen(AppleImg, screen, a.Pnt, PntToScrn)
	}
}
func (a *Apple) OnCollision(ownTile, collTile Ptype, collPcs interface{}) {
	if collTile == P_SnkH {
		colltype, ok := collPcs.(COLLISION)
		if ok {
			if colltype == AL_PL__N_APL {
				a.Eaten = true
				if a.NutritionValue < 0 {
					SoundEffects["Eating_Rotten"].PR()
				}else{
					SoundEffects["Eating"].PR()
				}
			}else if colltype == GH_PL__N_APL {
				if a.NutritionValue > 0 {
					timer := &AppleTimer{int(FPS*RottenAppleTime), a}
					SnakeProInGame.World.SObjs = append(SnakeProInGame.World.SObjs, timer)
				}
				a.NutritionValue = -int(AppleNutrition)
			}
		}
	}
}