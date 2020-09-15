package Utilities

import (
	"os"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"github.com/nfnt/resize"
	"image/color"
	"fmt"
	"math"
)
var TrueTypeFont *truetype.Font

func Init(path string) {
	font, err := ioutil.ReadFile(path)
   	if err != nil {
   		panic(err)
   	}
	tt, err := truetype.Parse(font)
	CheckErr(err)
	TrueTypeFont = tt
}

func PrintTextAt(screen *ebiten.Image, str string, x, y int, r, g, b uint8, size float64) {
	//ebitenutil.DebugPrintAt(screen, str, x, y)
	mplusNormalFont := truetype.NewFace(TrueTypeFont, &truetype.Options{
		Size:    size,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	pnt := text.MeasureString(str, mplusNormalFont)
	text.Draw(screen, str, mplusNormalFont, x-pnt.X/2, y-pnt.Y/2, color.RGBA{r,g,b,255})
}
func MakePopUpAt(screen *ebiten.Image, str string, x, y int, size float64, textCol, backCol color.RGBA) {
	w,_ := screen.Size()
	mplusNormalFont := truetype.NewFace(TrueTypeFont, &truetype.Options{
		Size:    size*(float64(w)/1920.0),
		DPI:     96,
		Hinting: font.HintingFull,
	})
	pnt := text.MeasureString(str, mplusNormalFont)
	
	w, h := int(float64(pnt.X)*1.2), int(float64(pnt.Y)*1.5)
	popUpBack, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	popUpBack.Fill(backCol)
	xP, yP := (w-pnt.X)/2, h-(h-pnt.Y)
	text.Draw(popUpBack, str, mplusNormalFont, int(xP), int(yP), textCol)
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x)-float64(w)/2.0, float64(y)-float64(h)/2)
	screen.DrawImage(popUpBack, op)
}

func LoadImgObj(path string, width, height, x, y, angle float64) *ImageObj {
	_, img := LoadImg(path)
	scaledImg := resize.Resize(uint(width), uint(height), *img, resize.NearestNeighbor)
	return &ImageObj{ImgToEbitenImg(&scaledImg), &scaledImg, width, height, x, y, angle}
}
type ImageObj struct {
	Img *ebiten.Image; OriginalImg *image.Image
	W, H, X, Y, Angle float64
}
func (obj *ImageObj) Print() string {
	return fmt.Sprintf("OW: %v, OH: %v, W: %f, H: %f, X: %f, Y: %f, Angle: %f", (*obj.OriginalImg).Bounds().Max.X, (*obj.OriginalImg).Bounds().Max.Y, obj.W, obj.H, obj.X, obj.Y, obj.Angle)
}
func (obj *ImageObj) Copy() *ImageObj {
	return &ImageObj{obj.Img, obj.OriginalImg, obj.W, obj.H, obj.X, obj.Y, obj.Angle}
}
func (obj *ImageObj) DrawImageObj(screen *ebiten.Image) {
	w, h := obj.Img.Size()
	op := &ebiten.DrawImageOptions{}
	xScale := obj.W/(float64)(w)
	yScale := obj.H/(float64)(h)

	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(float64(obj.Angle) * 2 * math.Pi / 360)
	op.GeoM.Scale(xScale, yScale)

	op.GeoM.Translate(obj.X+obj.W/2, obj.Y+obj.H/2)
	screen.DrawImage(obj.Img, op)
}
func (obj *ImageObj) DrawImageObjAlpha(screen *ebiten.Image, alpha float64) {
	w, h := obj.Img.Size()
	op := &ebiten.DrawImageOptions{}
	xScale := obj.W/(float64)(w)
	yScale := obj.H/(float64)(h)

	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(float64(obj.Angle) * 2 * math.Pi / 360)
	op.GeoM.Scale(xScale, yScale)

	op.GeoM.Translate(obj.X+obj.W/2, obj.Y+obj.H/2)
	
	op.ColorM.Scale(1, 1, 1, alpha)
	screen.DrawImage(obj.Img, op)
}


func LoadImg(path string) (error, *image.Image) {
	f, err := os.Open(path)
	if err != nil {
		return err, nil
	}
	img, _, err2 := image.Decode(f)
	if err2 != nil {
		return err, nil
	}
	return nil, &img
}
func ImgToEbitenImg(img *image.Image) (*ebiten.Image) {	
	gophersImage, err3 := ebiten.NewImageFromImage(*img, ebiten.FilterDefault)
	if err3 != nil {
		panic(fmt.Sprintf("Cannot Create Ebiten Image: %e",err3))
	}
	return gophersImage
}
func LoadEbitenImg(path string) (*ebiten.Image) {	
	err, img := LoadImg(path)
	if err != nil {
		panic(fmt.Sprintf("Cannot Load Image: %e",err))
	}
	return ImgToEbitenImg(img)
}

func InitIcons(path string, sizes []int, fileformat string) (error, []image.Image) {
	imgs := make([]image.Image, len(sizes))
	for i,_ := range(imgs) {
		err, img := LoadImg(fmt.Sprintf("%s/%v.%s", path, sizes[i], fileformat))
		if err != nil {
			return err, nil
		}
		imgs[i] = *img
	}
	return nil, imgs
}