package Game

import (
	"marvin/SnakeProGo/Utilities"
	"github.com/hajimehoshi/ebiten"
	"image"
	"math/rand"
	"time"
	"image/color"
	"log"
	"fmt"
)

const (
	SubImageScale = 2
	FPS = 30
)

var (
	SnakeProGame *SnakePro
	SnakeProInGame *InGame
	
	//Loaded Parameter Map
	Ps, IMGPOSITIONS *Utilities.Params
	lastPlayerNames *LastPlayerNames
	
	//Resource Directory
	res_dir string
	
	//Screen attributes
	XRes,YRes float64
	XTiles,YTiles float64
	MapToScreen *map[Utilities.Point]Utilities.Point
	
	//Loading Screen Intro
	Loading_Intro *Utilities.ImageObj
	backGroundImg *ebiten.Image
	
	//Player attributes
	PlayerNum, PlayerStartLenght int
	TileWidth, TileHeight float64
	PlayerSpeed, ControllerThresh float64
	ControllerIds []int
	SnakeColors = []color.RGBA{{0,0,255,255},{255,0,60,255},{255,255,0,255},{239,49,165,255}}
	
	//World attributes
	BeginItems, BeginApples, ItemSpawnPeriod, AppleNutrition uint
	PlayerAliveAloneFrames int
	RottenAppleTime float64
	
	//Stats
	RuntimeStatistics *RuntimeStats
	
	//Menu
	Exit_X, Exit_Y, Play_X, Play_Y, Con_X, Con_Y, HOF_X, HOF_Y float64
	PLN_XPOS, PLN_YPOS, PLN_WIDTH, PLN_HEIGHT float64
	Win_Head_Pos_X, Win_Head_Rot,Win_Head_S float64
	
	//Sounds
	SoundEffects map[string]*SoundEffect
	
	//Items
	LaserLength, BotLength int
	ReviveImmortalTime, LaserTime, BombTime, BotSpeed, SpeedLength, SpeedStrength, FartTime, FartStrength float64
	BombLength, BombRadius, FartRadius int
)

func InitSounds() {
	Utilities.InitAudioContext()
	SoundEffectsPath := f_Resources+res_dir+f_Sounds
	SoundEffects = make(map[string]*SoundEffect)
	SoundEffects["Dead"] = LoadSoundEffect(SoundEffectsPath+				f_Dead)
	SoundEffects["Eating"] = LoadSoundEffect(SoundEffectsPath+				f_Eating)
	SoundEffects["Eating_Rotten"] = LoadSoundEffect(SoundEffectsPath+		f_Eating_Rotten)
	SoundEffects["Explosion"] = LoadSoundEffect(SoundEffectsPath+			f_Explosion)
	SoundEffects["Farting"] = LoadSoundEffect(SoundEffectsPath+				f_Farting)
	SoundEffects["Item"] = LoadSoundEffect(SoundEffectsPath+				f_Item)
	SoundEffects["Revive"] = LoadSoundEffect(SoundEffectsPath+				f_Revive)
	SoundEffects["Laser"] = LoadSoundEffect(SoundEffectsPath+			f_Shooting)
	SoundEffects["Speed"] = LoadSoundEffect(SoundEffectsPath+				f_Speed)
	SoundEffects["Win"] = LoadSoundEffect(SoundEffectsPath+					f_Win)
}

func StartGame(game ebiten.Game, name string, icons []image.Image) {
	ebiten.SetFullscreen(true)
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetWindowTitle(name)
	ebiten.SetWindowIcon(icons)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetMaxTPS(FPS)
	if err := ebiten.RunGame(game); err != nil {
		EXITGAME()
		log.Fatal(err)
	}
	EXITGAME()
}

type Teststruct struct {
	Fett map[string]int
	Fettheit2 [][]string
}
type Teststruct2 struct {
	Teststruct
	Fettheit [][]string
}

func Start() {
	rand.Seed(time.Now().UnixNano())
	Ps = &Utilities.Params{}
	Ps.LoadFromFile(f_Resources+f_Params)
	W, H := ebiten.ScreenSizeInFullscreen();XRes,YRes = float64(W),float64(H)
	XTiles,YTiles = Ps.Get("XTiles"), Ps.Get("YTiles")
	TileWidth, TileHeight = (XRes/XTiles)*(SubImageScale),(YRes/YTiles)*(SubImageScale)
	res_dir = Ps.GetS("ResourceLoc")
	PlayerNum = int(Ps.Get("Player"))
	PlayerSpeed = Ps.Get("PlayerSpeed")
	BeginItems, BeginApples, ItemSpawnPeriod = uint(Ps.Get("BeginItems")),uint(Ps.Get("BeginApples")),uint(Ps.Get("ItemSpawnPeriod"))*FPS
	PlayerStartLenght = int(Ps.Get("PlayerStartLenght"))
	ControllerThresh = Ps.Get("ControllerThresh")
	AppleNutrition = uint(Ps.Get("AppleNutrition"))
	PlayerAliveAloneFrames = int(Ps.Get("PlayerAliveAloneFrames")*FPS)
	SpeedLength = Ps.Get("SpeedLength")
	SpeedStrength = Ps.Get("SpeedStrength")
	LaserLength	 = int(Ps.Get("LaserLength"))
	ReviveImmortalTime = Ps.Get("ReviveImmortalTime")
	LaserTime = Ps.Get("LaserTime")
	BombTime = Ps.Get("BombTime")
	BombLength = int(Ps.Get("BombLength"))
	BombRadius = int(Ps.Get("BombRadius"))
	BotLength = int(Ps.Get("BotLength"))
	BotSpeed = Ps.Get("BotSpeed")
	FartRadius = int(Ps.Get("FartRadius"))
	FartTime = Ps.Get("FartTime")
	FartStrength = Ps.Get("FartStrength")
	RottenAppleTime = Ps.Get("RottenAppleTime")
	for i := 0; i < PlayerNum; i++ {
		ControllerIds = append(ControllerIds, i)
	}
	InitSounds()
	lastPlayerNames, _ = LoadLastPlayerNames(f_Resources+f_LastPlayerNames)
	
	fmt.Println(Ps.Print())
	
	IMGPOSITIONS = &Utilities.Params{}
	IMGPOSITIONS.LoadFromFile(f_Resources+res_dir+f_Positions)
	Exit_X, Exit_Y, Play_X, Play_Y, Con_X, Con_Y = IMGPOSITIONS.Get("Exit_X"),IMGPOSITIONS.Get("Exit_Y"),IMGPOSITIONS.Get("Play_X"),IMGPOSITIONS.Get("Play_Y"),IMGPOSITIONS.Get("Con_X"),IMGPOSITIONS.Get("Con_Y")
	PLN_XPOS, PLN_YPOS, PLN_WIDTH, PLN_HEIGHT = IMGPOSITIONS.Get("PLN_XPOS"),IMGPOSITIONS.Get("PLN_YPOS"),IMGPOSITIONS.Get("PLN_WIDTH"),IMGPOSITIONS.Get("PLN_HEIGHT")
	Win_Head_Pos_X, Win_Head_Rot,Win_Head_S = IMGPOSITIONS.Get("Win_Head_Pos_X"),IMGPOSITIONS.Get("Win_Head_Rot"),IMGPOSITIONS.Get("Win_Head_S")
	HOF_X, HOF_Y = IMGPOSITIONS.Get("HOF_X"),IMGPOSITIONS.Get("HOF_Y")
	fmt.Println(IMGPOSITIONS.Print())
	
	Utilities.Init(f_Resources+res_dir+f_Font+f_Comic)
	
	PlayerSnakeBodies = make([][][]*Utilities.ImageObj, 4)
	for i := 0; i<4; i++ {
		PlayerSnakeBodies[i] = GetPlayerTiles(i)
	}
	
	MapToScreenP := make(map[Utilities.Point]Utilities.Point)
	for x := 0; x < int(XTiles); x++ {
		for y := 0; y < int(YTiles); y++ {
			nX := float64(x)*(TileWidth/SubImageScale)-(TileWidth/SubImageScale)*0.5
			nY := float64(y)*(TileHeight/SubImageScale)-(TileHeight/SubImageScale)*0.5
			MapToScreenP[Utilities.Point{x,y}] = Utilities.Point{int(nX),int(nY)}
		}
	}
	MapToScreen = &MapToScreenP
	
	RuntimeStatistics = &RuntimeStats{}
	RuntimeStatistics.Init([]string{"Update", "UpConsGen", "UpButtons1", "UpButtons2", "UpButtons3", "UpListener1", "UpListener2", "UpListener3", "PlayUp", "WorldUp", "Draw", "FPS", "TPS"})
	
	SnakeProGame = &SnakePro{}
	SnakeProGame.Init()
	
	err, icons := Utilities.InitIcons(f_Resources+res_dir+f_Icons, []int{16,32,48,64,128,256}, f_IconFormat)
	Utilities.CheckErr(err)
	
	Loading_Intro = Utilities.LoadImgObj(f_Resources+res_dir+f_Icons+f_Intro, XRes/4, 0, XRes/4*1.5, YRes/4, 0)
	Loading_Intro.H = float64((*Loading_Intro.OriginalImg).Bounds().Max.Y)
	fmt.Println(Loading_Intro.Print())
	
	backGroundImg, _ = ebiten.NewImage(int(XRes), int(YRes), ebiten.FilterDefault)
	backGroundImg.Fill(color.RGBA{60,60,60,255})
	
	StartGame(SnakeProGame, "SnakePro", icons)
}