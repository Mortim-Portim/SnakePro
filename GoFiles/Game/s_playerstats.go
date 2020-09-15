package Game

import (
	"encoding/json"
	"io/ioutil"
	//"errors"
	"fmt"
)

type LastPlayerNames struct {
	L1, L2, L3, L4 string
}
func (l *LastPlayerNames) Set(i int, name string) {
	if l != nil {
	if i == 0 {
		l.L1 = name
	}
	if i == 1 {
		l.L2 = name
	}
	if i == 2 {
		l.L3 = name
	}
	if i == 3 {
		l.L4 = name
	}
	}
}
func (l *LastPlayerNames) Get(i int) string {
	if l != nil {
	if i == 0 && len(l.L1) > 0 {
		return l.L1
	}
	if i == 1 && len(l.L2) > 0 {
		return l.L2
	}
	if i == 2 && len(l.L3) > 0 {
		return l.L3
	}
	if i == 3 && len(l.L4) > 0 {
		return l.L4
	}
	}
	return fmt.Sprintf("Player%v",i+1)
}
func (l *LastPlayerNames) Save(path string) {
	if l != nil {
		SaveStruct(path, *l)
	}
}
func LoadLastPlayerNames(path string) (*LastPlayerNames, error) {
	dat, err := ioutil.ReadFile(path)
   	if err != nil {
   		return &LastPlayerNames{}, err
   	}
   	l := &LastPlayerNames{}
	err2 := json.Unmarshal(dat, l)
	if err2 != nil {
		return &LastPlayerNames{}, err2
   	}
	return l, nil
}

type PlayerStats struct {
	Statistics
	name, txt, csv string
}
func GetNewPlayerStatistic(name string, strs []string) (s *PlayerStats) {
	s = &PlayerStats{}
	s.Init(strs)
	s.SetName(name)
	return
}

func (s *PlayerStats) SetName(name string) {
	s.name = name
	s.txt = fmt.Sprintf("%s%s/%s.txt", f_Resources, f_Statistics, s.name)
	s.csv = fmt.Sprintf("%s%s/%s.csv", f_Resources, f_Statistics, s.name)
}
func (s *PlayerStats) Name() string {
	return s.name
}
func (s *PlayerStats) Load() error {
	fmt.Println("Loading Stats")
	err := LoadStatistics(s.txt, &s.Statistics)
	return err
}
func (s *PlayerStats) Save() {
	fmt.Println("Saving Stats")
	SaveStruct(s.txt, s.Statistics)
	s.Statistics.WriteToCSV(s.csv)
}

func LoadStatistics(path string, obj *Statistics) error {
	dat, err := ioutil.ReadFile(path)
   	if err != nil {
   		return err
   	}
	err2 := json.Unmarshal(dat, &obj)
	if err2 != nil {
		return err2
   	}
	return nil
}
func SaveStruct(path string, obj interface{}) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	err2 := ioutil.WriteFile(path, bytes, 0644)
    if err2 != nil {
		panic(err2)
	}
}