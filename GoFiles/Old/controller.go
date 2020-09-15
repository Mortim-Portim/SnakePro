package Old
/**
import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"fmt"
	"log"
	"math"
	"errors"
)

const (
	controllerThresh = 0.01
)

var (
	GamepadSDLIDs map[int]string
)

func UpdateControlls() {
	Ids := ebiten.GamepadIDs()
	GamepadSDLIDs = map[int]string{}
	for i, id := range(Ids) {
		GamepadSDLIDs[i] = ebiten.GamepadSDLID(id)
	}
}

func InitControlls(sensi float64, playerNum int) (error, []*Controller) {
	if playerNum <= 0 {
		return errors.New(fmt.Sprintf("Not enough Player(%v) to Initialize Controlls", playerNum)), nil
	}
	Ids := ebiten.GamepadIDs()
	if len(Ids) <= 0 {
		log.Println("InitControlls: ebiten might not have been initialized")
	}
	cons := make([]*Controller, playerNum)
	cons[0] = &Controller{Id:"000", Sensitivity:sensi};cons[0].Update()
	fmt.Println("\nController Initialized:")
	fmt.Println(cons[0].GetInfos())
	for i := 1; i < playerNum; i++ {
		tmp_Con := &Controller{Id:"000", Sensitivity:sensi, Mapper:Get1zu1Mapper(100)}
		if i <= len(Ids) {
			tmp_Con = &Controller{Id:ebiten.GamepadSDLID(Ids[i-1]), Sensitivity:sensi}
		}
		tmp_Con.Update()
		cons[i] = tmp_Con
		fmt.Println("\nController Initialized:")
		fmt.Println(tmp_Con.GetInfos())
	}
	return nil, cons
}

func Get1zu1Mapper(length int) map[int]int {
	mapper := make(map[int]int, length)
	for i := 0; i < length; i++ {
		mapper[i] = i
	}
	return mapper
}

type Controller struct {
	Id, Name string
	AxisValues []float64
	buttons map[int]bool
	Mapper map[int]int
	
	JustDown []int
	JustUp []int
	
	Sensitivity float64
}

func (c *Controller) IsButtonDown(idx int) bool {
	return c.buttons[c.Mapper[idx]]
}

func (c *Controller) SetMapper(oldidx, newidx int) {
	c.Mapper[oldidx] = newidx
}

func (c *Controller) GetInfos() (out string) {
	out = ""
	if c.Id == "000" {
		out += "Keyboard: "
	}else{
		out += "Gamepad: "
	}
	out += fmt.Sprintf("SDL ID: %s, Name: %s\n", c.Id, c.Name)
	out += fmt.Sprintf("  Axes:	%f\n", c.AxisValues)
	out += fmt.Sprintf("  JustDown:	%v\n", c.JustDown)
	out += fmt.Sprintf("  Sensitivity: %v\n", c.Sensitivity)
	return
}

func (c *Controller) Update() {
	if c.Id != "000" {
		var tmp_ID int
		for id, SDL := range(GamepadSDLIDs) {
			if SDL == c.Id {
				tmp_ID = id
			}
		}
		c.Name = ebiten.GamepadName(tmp_ID)
		
		axisNum := ebiten.GamepadAxisNum(tmp_ID)
		buttonNum := ebiten.GamepadButton(ebiten.GamepadButtonNum(tmp_ID))
		
		c.AxisValues = make([]float64, axisNum)
		for a := 0; a < axisNum; a++ {
			c.AxisValues[a] = 0
			val := ebiten.GamepadAxis(tmp_ID, a)
			if math.Abs(val) > controllerThresh {
				c.AxisValues[a] = val
			}
		}
		
		c.buttons = make(map[int]bool, 0)
		c.JustDown = make([]int,0)
		c.JustUp = make([]int, 0)
		counter := 0
		for b := ebiten.GamepadButton(tmp_ID); b < buttonNum; b++ {
			c.buttons[counter] = ebiten.IsGamepadButtonPressed(tmp_ID, b)
			
			// Log button events.
			if inpututil.IsGamepadButtonJustPressed(tmp_ID, b) {
				c.JustDown = append(c.JustDown, counter)
			}
			if inpututil.IsGamepadButtonJustReleased(tmp_ID, b) {
				c.JustUp = append(c.JustUp, counter)
			}
			
			counter ++
		}
	}else{
		c.Name = "Keyboard"
		c.AxisValues = make([]float64, 2)
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			c.AxisValues[0] = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			c.AxisValues[0] = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			c.AxisValues[1] = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			c.AxisValues[1] = 1
		}
		
		c.JustDown = make([]int,0)
		c.JustUp = make([]int, 0)
		
		for i,k := range(AllKeys) {
			KeyState := ebiten.IsKeyPressed(k)
			if KeyState != c.buttons[i] {
				if KeyState {
					c.JustDown = append(c.JustDown, i)
				}else{
					c.JustUp = append(c.JustUp, i)
				}
			}
		}
		
		c.buttons = make(map[int]bool, 0)
		for i,k := range(AllKeys) {
			c.buttons[i] = ebiten.IsKeyPressed(k)
		}
	}
}
**/