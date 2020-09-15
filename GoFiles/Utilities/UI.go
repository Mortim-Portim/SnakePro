package Utilities

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	//"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image/color"
	//"math"
)
const (
	TextSizeScale = 0.75
)

func GetWHForText(str string, size float64) (float64, float64) {
	mplusNormalFont := truetype.NewFace(TrueTypeFont, &truetype.Options{
		Size:    size,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	pnt := text.MeasureString(str, mplusNormalFont)
	return float64(pnt.X), float64(pnt.Y)
}

func GetSizeForText(str string, maxWidth, maxHeight int) (size float64) {
	mplusNormalFont := truetype.NewFace(TrueTypeFont, &truetype.Options{
		Size:    24,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	pnt := text.MeasureString(str, mplusNormalFont)
	factor := 1.0
	if float64(pnt.X-maxWidth) > float64(pnt.Y-maxHeight) {
		factor = float64(maxWidth)/float64(pnt.X)
	}else{
		factor = float64(maxHeight)/float64(pnt.Y)
	}
	return 24*factor
}

func GetButton(img *ImageObj, onPressLeft func(b *Button), onPressRight func(b *Button)) *Button {
	b := &Button{}
	b.Img = img; b.onPressLeft = onPressLeft; b.onPressRight = onPressRight; b.Active = true
	return b
}
func GetButtonWithText(img *ImageObj, str string, textCol color.RGBA, onPressLeft func(b *Button), onPressRight func(b *Button)) *Button {
	b := &Button{}
	b.onPressLeft = onPressLeft; b.onPressRight = onPressRight; b.Active = true
	
	w, h := img.Img.Size()
	mplusNormalFont := truetype.NewFace(TrueTypeFont, &truetype.Options{
		Size:    GetSizeForText(str, int(float64(w)*0.95), int(float64(h)*0.85)),
		DPI:     96,
		Hinting: font.HintingFull,
	})
	pnt := text.MeasureString(str, mplusNormalFont)
	
	xP, yP := 0, h/2+pnt.Y/4
	text.Draw(img.Img, str, mplusNormalFont, int(xP), int(yP), textCol)
	
	b.Img = img;
	return b
}
func GetButtonWithTextBack(str string, X, Y, W, H float64, textCol, backCol color.RGBA, onPressLeft func(b *Button), onPressRight func(b *Button)) *Button {
	imgObj := &ImageObj{}
	
	Back, _ := ebiten.NewImage(int(W), int(H), ebiten.FilterDefault)
	Back.Fill(backCol)
	imgObj.Img = Back
	imgObj.X = X;imgObj.Y = Y;imgObj.W = W;imgObj.H = H;imgObj.Angle = 0;
	
	return GetButtonWithText(imgObj, str, textCol, onPressLeft, onPressRight)
}

func GetButtonWithTextOnyl(str string, X, Y, W, H float64, textCol color.RGBA, onPressLeft func(b *Button), onPressRight func(b *Button)) *Button {
	return GetButtonWithTextBack(str, X,Y,W,H,textCol,color.RGBA{0,0,0,0},onPressLeft,onPressRight)
}

type Button struct {
	Img *ImageObj
		
	LPressed, RPressed, LastL, LastR, Active bool
	onPressLeft, onPressRight func(b *Button)
	Data interface{}
}
func (b *Button) Update() {
	if b.Active {
		b.LPressed = false
		b.RPressed = false
		x, y := ebiten.CursorPosition()
		if int(b.Img.X) <= x && x < int(b.Img.X+b.Img.W) && int(b.Img.Y) <= y && y < int(b.Img.Y+b.Img.H) {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				b.LPressed = true
			}
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
				b.RPressed = true
			}
		}
		if b.LPressed != b.LastL {
			if b.onPressLeft != nil {
				b.onPressLeft(b)
			}
		}
		b.LastL = b.LPressed
		
		if b.RPressed != b.LastR {
			if b.onPressRight != nil {
				b.onPressRight(b)
			}
		}
		b.LastR = b.RPressed
	}
}
func (b *Button) Draw(screen *ebiten.Image) {
	if b.Active {
		b.Img.DrawImageObj(screen)
	}
}
