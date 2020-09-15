package Utilities

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Params struct {
	p map[string]float64
	strs map[string]string
}

func (p *Params) LoadFromFile(path string) error {
	f, err := os.Open(path)
	CheckErr(err)
	scanner := bufio.NewScanner(f)
	p.p = make(map[string]float64)
	p.strs = make(map[string]string)
    for scanner.Scan() {
    	line := scanner.Text()
    	line = strings.ReplaceAll(line, " ", "")
    	ps := strings.Split(line, ":")
    	if len(ps) >= 2 {
    		fl, err2 := strconv.ParseFloat(ps[1], 64)
    		if err2 != nil {
	    		p.strs[ps[0]] = ps[1]
    		}else{
	    		p.p[ps[0]] = fl
    		}
    	}
    }
	
	return nil
}

func (p *Params) GetS(key string) (string) {
	if val, ok := p.strs[key]; ok {
	    return val
	}
	return ""
}

func (p *Params) Get(key string) (float64) {
	if val, ok := p.p[key]; ok {
	    return val
	}
	return 0
}

func (p *Params) Print() string {
	out := ""
	for k,v := range(p.p) {
		out += fmt.Sprintf("%s : %0.4f\n", k, v)
	}
	return out
}