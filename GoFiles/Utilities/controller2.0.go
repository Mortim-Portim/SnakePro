package Utilities

import (
	"github.com/hajimehoshi/ebiten"
	"fmt"
	"log"
	"math"
	"errors"
	"encoding/json"
	"io/ioutil"
	"time"
)

type DirectionEventListener interface {
	GetAxis() (int,int,int,int)
	OnDirectionEvent(xdif, ydif float64)
}
type ButtonEventListener interface {
	GetButton() (int)
	OnButtonUp()
	OnButtonDown()
}

var (
	controllerThresh = 0.1
)

var (
	GamepadSDLIDs map[int]string
	SmartServIPAddr string
)

func UpdateControlls() {
	Ids := ebiten.GamepadIDs()
	GamepadSDLIDs = map[int]string{}
	for i, id := range(Ids) {
		GamepadSDLIDs[i] = ebiten.GamepadSDLID(id)
	}
}

func InitControlls(ControllerThresh float64, playerNum int) (error, *[]*Controller) {
	controllerThresh = controllerThresh
	if playerNum <= 0 {
		return errors.New(fmt.Sprintf("Not enough Player(%v) to Initialize Controlls", playerNum)), nil
	}
	Ids := ebiten.GamepadIDs()
	if len(Ids) <= 0 {
		log.Println("InitControlls: ebiten might not have been initialized")
	}
	cons := make([]*Controller, playerNum)
	cons[0] = &Controller{Id:"000", Mapper:Get1zu1Mapper(100)};cons[0].Buttons = make(map[int]float64, 0);cons[0].UpdateAll()
	fmt.Println("\nController Initialized:")
	fmt.Println(cons[0].GetInfos())
	for i := 1; i < playerNum; i++ {
		tmp_Con := &Controller{Id:"000", Mapper:Get1zu1Mapper(100)}
		if i <= len(Ids) {
			tmp_Con = &Controller{Id:ebiten.GamepadSDLID(Ids[i-1]), Mapper:Get1zu1Mapper(100)}
		}
		tmp_Con.Buttons = make(map[int]float64, 0)
		tmp_Con.UpdateAll()
		cons[i] = tmp_Con
		fmt.Println("\nController Initialized:")
		fmt.Println(tmp_Con.GetInfos())
	}
	
	IpToStruct = make(map[string]*SmartPhoneConn)
	_,ip := NewSmartPhoneServ(&cons)
	SmartServIPAddr = ip
	return nil, &cons
}

func Get1zu1Mapper(length int) map[int]int {
	mapper := make(map[int]int, length)
	for i := 0; i < length; i++ {
		mapper[i] = i
	}
	return mapper
}

func LoadMapper(path string) map[int]int {
	dat, err := ioutil.ReadFile(path)
   	if err != nil {
	   	return nil
   	}
	var newMapper map[int]int
	err2 := json.Unmarshal(dat, &newMapper)
	if err2 != nil {
	   	return nil
   	}
	return newMapper
}
func SaveMapper(path string, mapper map[int]int) {
	bytes, err := json.Marshal(mapper)
	CheckErr(err)
	err2 := ioutil.WriteFile(path, bytes, 0644)
    CheckErr(err2)
}


type Controller struct {
	Id, Name string
	Buttons map[int]float64
	Mapper map[int]int
	
	AxisStandard map[int]float64
	AxisValues []float64
	
	axisEvents []DirectionEventListener
	buttonEvents []ButtonEventListener
	
	Down []int
	
	SmartPhoneConn *SmartPhoneConn
}
func (c *Controller) SaveConfig(path string) {
	SaveMapper(fmt.Sprintf("%s/%s.txt", path, c.Id), c.Mapper)
}
func (c *Controller) LoadConfig(path string) {
	mapper := LoadMapper(fmt.Sprintf("%s/%s.txt", path, c.Id))
	if mapper != nil {
		c.Mapper = mapper
	}
}
func (c *Controller) ResetDirectionEventListener() {
	c.axisEvents = make([]DirectionEventListener, 0)
}
func (c *Controller) ResetButtonEventListener() {
	c.buttonEvents = make([]ButtonEventListener, 0)
}
func (c *Controller) RegisterDirectionEventListener(listener DirectionEventListener) {
	c.axisEvents = append(c.axisEvents, listener)
}
func (c *Controller) RegisterButtonEventListener(listener ButtonEventListener) {
	c.buttonEvents = append(c.buttonEvents, listener)
}


func (c *Controller) GetButtonValue(idx int) float64 {
	return c.Buttons[c.Mapper[idx]]
}
func (c *Controller) GetAxisValue(idx1, idx2 int) float64 {
	return c.Buttons[c.Mapper[idx2]]-c.Buttons[c.Mapper[idx1]]
}

func (c *Controller) SetMapper(oldidx, newidx int) {
	c.Mapper[oldidx] = newidx
}

func (c *Controller) GetInfos() (out string) {
	out = ""
	if c.Id == "000" {
		out += "Keyboard: "
	}else if c.Id == "SPS" {
		out += "Smartphone: "
	}else{
		out += "Gamepad: "
	}
	out += fmt.Sprintf("SDL ID: %s, Name: %s\n", c.Id, c.Name)
	out += fmt.Sprintf("  Down:	%v\n", c.Down)
	return
}

func (c *Controller) UseCurrentAxisValsAsStandard() {
	c.AxisStandard = make(map[int]float64)
	c.UpdateAll()
	for i,v := range(c.AxisValues) {
		if v != 0 {
			c.AxisStandard[i] = -v
		} 
	}
	c.UpdateAll()
}
func (c *Controller) UpdateSmartPhone() {
	c.Down = make([]int,0)
	for mapi := 0; mapi < 5; mapi ++ {
		i := c.Mapper[mapi]
		BtnState := 0.0
		if c.SmartPhoneConn.IsDown(i) {
			BtnState = 1.0
			c.Down = append(c.Down, i)
		}
		c.Buttons[i] = BtnState
	}
}
func (c *Controller) UpdateKeys(keys ...int) {
	c.Down = make([]int,0)
	for _,mapi := range(keys) {
		i := c.Mapper[mapi]
		k := AllKeys[i]
		KeyState := 0.0
		if ebiten.IsKeyPressed(k) {
			KeyState = 1.0
			c.Down = append(c.Down, i)
		}
		c.Buttons[i] = KeyState
	}
}
func (c *Controller) UpdateControllerButtons(buttons ...int) {
	var tmp_ID int
	for id, SDL := range(GamepadSDLIDs) {
		if SDL == c.Id {
			tmp_ID = id
		}
	}
	c.Name = ebiten.GamepadName(tmp_ID)
	axisNum := ebiten.GamepadAxisNum(tmp_ID)
	
	c.AxisValues = make([]float64, axisNum)
	for a := 0; a < axisNum; a++ {
		c.AxisValues[a] = 0
		val := ebiten.GamepadAxis(tmp_ID, a)
		if math.Abs(val+c.AxisStandard[a]) > controllerThresh {
			c.AxisValues[a] = val+c.AxisStandard[a]
		}
	}
	
	c.Buttons = make(map[int]float64, 0)
	c.Down = make([]int,0)
	for _,btnIdx := range(buttons) {
		b := ebiten.GamepadButton(btnIdx)
		c.Buttons[btnIdx] = 0.0
		if ebiten.IsGamepadButtonPressed(tmp_ID, b) {
			c.Buttons[btnIdx] = 1.0
			c.Down = append(c.Down, btnIdx)
		}
	}
	counter := int(ebiten.GamepadButton(ebiten.GamepadButtonNum(tmp_ID)))
	for _,v := range(c.AxisValues) {
		c.Buttons[counter] = 0.0
		c.Buttons[counter+1] = 0.0
		if v != 0 {
			if v < 0 {
				c.Buttons[counter] = math.Abs(v)
				c.Down = append(c.Down, counter)
			}else{
				c.Buttons[counter+1] = math.Abs(v)
				c.Down = append(c.Down, counter+1)
			}
		}
		counter += 2
	}
}
func (c *Controller) UpdateOnlyNeeded() (int,int) {
	start := time.Now()
	if c.Id != "000" && c.Id != "SPS" {
		c.UpdateController()
	}else if c.Id == "SPS" {
		c.UpdateSmartPhone()
	}else{
		c.Name = "Keyboard"
		c.UpdateKeyBoard()
	}
	BT := int(time.Since(start))
	start2 := time.Now()
	c.UpdateListeners()
	LT := int(time.Since(start2))
	return BT,LT
	//fmt.Println(c.JustDown, ":", c.JustUp)
}
func (c *Controller) UpdateAll() (int,int) {
	start := time.Now()
	if c.Id != "000" && c.Id != "SPS" {
		c.UpdateControllerAll()
	}else if c.Id == "SPS" {
		c.UpdateSmartPhone()
	}else{
		c.Name = "Keyboard"
		c.UpdateKeyBoardAll()
	}
	BT := int(time.Since(start))
	start2 := time.Now()
	c.UpdateListeners()
	LT := int(time.Since(start2))
	return BT,LT
	//fmt.Println(c.JustDown, ":", c.JustUp)
}
func (c *Controller) UpdateController() {
	c.UpdateControllerButtons(c.GetListenedBtns()...)
}
func (c *Controller) UpdateControllerAll() {
	var tmp_ID int
	for id, SDL := range(GamepadSDLIDs) {
		if SDL == c.Id {
			tmp_ID = id
		}
	}
	buttonNum := ebiten.GamepadButton(ebiten.GamepadButtonNum(tmp_ID))
	upBtn := make([]int, buttonNum)
	for i,_ := range(upBtn) {
		upBtn[i] = i
	}
	c.UpdateControllerButtons(upBtn...)
}
func (c *Controller) UpdateKeyBoard() {
	c.UpdateKeys(c.GetListenedBtns()...)
}
func (c *Controller) UpdateKeyBoardAll() {
	upKeys := make([]int, len(AllKeys))
	for i,_ := range(upKeys) {
		upKeys[i] = i
	}
	c.UpdateKeys(upKeys...)
}
func (c *Controller) UpdateListeners() {
	for _,axisListener := range(c.axisEvents) {
		x1,x2,y1,y2 := axisListener.GetAxis()
		if c.GetAxisValue(x1,x2) != 0 || c.GetAxisValue(y1,y2) != 0 {
			axisListener.OnDirectionEvent(c.GetAxisValue(x1,x2), c.GetAxisValue(y1,y2))
		}
	}
	for _,btnListener := range(c.buttonEvents) {
		btn := btnListener.GetButton()
		if contains(c.Down, c.Mapper[btn]) {
			btnListener.OnButtonDown()
		}else{
			btnListener.OnButtonUp()
		}
	}
}
func (c *Controller) GetListenedBtns() []int {
	upKeys := make([]int, 0)
	for _,axisListener := range(c.axisEvents) {
		x1,x2,y1,y2 := axisListener.GetAxis()
		upKeys = append(upKeys, c.Mapper[x1], c.Mapper[x2], c.Mapper[y1], c.Mapper[y2])
	}
	for _,btnListener := range(c.buttonEvents) {
		btn := btnListener.GetButton()
		upKeys = append(upKeys, c.Mapper[btn])
	}
	return upKeys
}
func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}