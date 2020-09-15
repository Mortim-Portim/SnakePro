package Utilities

import (
	"fmt"
	"marvin/GameServer/serv"
	"github.com/tidwall/evio"
)

const PORT = 8080
const INITSTATE = "00000\n"

var (
	SPS_serv *SmartPhoneServ
	IpToStruct map[string]*SmartPhoneConn
)

func OnNewConnection(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
	IpAddr := c.RemoteAddr().String()
	fmt.Printf("Smartphone conntected: %s\n",IpAddr)
	
	tmp_Con := &Controller{Id:"SPS", Mapper:Get1zu1Mapper(10), Name:IpAddr}
	tmp_Con.Buttons = make(map[int]float64, 0)
	tmp_Con.SmartPhoneConn = &SmartPhoneConn{Enc:INITSTATE}
	tmp_Con.SmartPhoneConn.EncToBool()
	tmp_Con.UpdateAll()
	
	*SPS_serv.Cons = append(*SPS_serv.Cons, tmp_Con)
	
	IpToStruct[IpAddr] = tmp_Con.SmartPhoneConn
	
	out = []byte("\n")
	
	return
}
func OnAcceptData(c evio.Conn, in []byte) (out []byte, action evio.Action) {
	IpAddr := c.RemoteAddr().String()
	fmt.Printf("%s: %s", IpAddr, string(in))
	
	IpToStruct[IpAddr].SetEnc(string(in))
	IpToStruct[IpAddr].EncToBool()
	
	
	return
}
func OnConnectionClosed(c evio.Conn, err error) (action evio.Action) {
	IpAddr := c.RemoteAddr().String()
	fmt.Printf("Smartphone disconnected: %s\n", IpAddr)
	delete(IpToStruct, IpAddr);
	return
}

type SmartPhoneConn struct {
	btns []bool
	Enc string
}
func (c *SmartPhoneConn) SetEnc(enc string) {
	c.Enc = enc
}
func (c *SmartPhoneConn) EncToBool() {
	if len(c.Enc) > 0 {
		c.btns = make([]bool, len(c.Enc))
		for i,s := range(c.Enc[:5]) {
			if string(s) == "0" {
				c.btns[i] = false
			}else{
				c.btns[i] = true
			}
		}
	}
}
func (c *SmartPhoneConn) IsDown(i int) bool {
	return c.btns[i]
}

type SmartPhoneServ struct {
	Cons *[]*Controller
}
func (s *SmartPhoneServ) Run() {
	tcp := serv.TCP{Port:PORT,OnNewConnection:OnNewConnection,OnAcceptData:OnAcceptData,OnConnectionClosed:OnConnectionClosed}
	tcp.Serve()
}
func NewSmartPhoneServ(cons *[]*Controller) (*SmartPhoneServ, string) {
	SPS_serv = &SmartPhoneServ{cons}
	SPS_serv.Run()
	return SPS_serv, serv.GetOutboundIP()
}