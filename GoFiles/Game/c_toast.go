package Game

import (
	"github.com/hajimehoshi/ebiten"
	"marvin/SnakeProGo/Utilities"
	"image/color"
)

type Toast struct {
	btn *Utilities.Button
	RemainingTime int
}
func (t *Toast) Draw(screen *ebiten.Image, idx int) {
	t.btn.Img.Y += t.btn.Img.H*float64(idx)
	t.btn.Draw(screen)
	t.btn.Img.Y -= t.btn.Img.H*float64(idx)
}

func GetToast(str string, time float64, textCol, backCol color.RGBA, X,Y,W,H float64) *Toast {
	newToast := &Toast{}
	newToast.RemainingTime = int(time*FPS)
	newToast.btn = Utilities.GetButtonWithTextBack(str, X, Y, W, H, textCol, backCol, nil, nil)
	return newToast
}
func GetToastUp(str string, time float64, textCol, backCol color.RGBA) *Toast {
	W, H := Utilities.GetWHForText(str, 20)
	return GetToast(str, time, textCol, backCol, XRes/2-W/2, 0, W, H)
}
func GetStandardToastShort(str string) *Toast {
	return GetToastUp(str, 2, color.RGBA{255,255,255,255}, color.RGBA{0,0,0,190})
}