package Game

import (
	"encoding/csv"
	"strings"
	"fmt"
	"log"
	"os"
)
type RuntimeStats struct {
	Statistics
}

type Statistics struct {
	Header []string
	Stats [][]string
	Head map[string]int
}
func (r *Statistics) Print() string {
	out := "{"
	for i,column := range(r.Stats) {
		out += "["+strings.Join(column, "; ") + "]"
		if i+1 < len(r.Stats) {
			out += "\n"
		}
	}
	return out+"}"
}
func (r *Statistics) Init(strs []string) {
	r.Header = strs
	r.Stats = make([][]string, 1)
	r.Stats[0] = strs
	
	r.ResetHeader()
}

func (r *Statistics) SetHead(head string, val int) {
	r.Head[head] = val
}
func (r *Statistics) AddHead(head string, val int) {
	r.Head[head] = r.Head[head]+val
}
func (r *Statistics) GetHead(head string) (int) {
	return r.Head[head]
}

func (r *Statistics) Append() {
	vals := make([]int, len(r.Header))
	for i,s := range(r.Header) {
	    vals[i] = r.Head[s]
	}
	strs := valsToString(vals)
	r.Stats = append(r.Stats, strs)
	r.ResetHeader()
}
func (r *Statistics) Set() {
	vals := make([]int, len(r.Header))
	for i,s := range(r.Header) {
	    vals[i] = r.Head[s]
	}
	strs := valsToString(vals)
	r.Stats = make([][]string, 2)
	r.Stats[0] = r.Header
	r.Stats[1] = strs
}

func (r *Statistics) ResetHeader() {
	r.Head = make(map[string]int, 0)
	for _,str := range(r.Header) {
		r.Head[str] = 0
	}
}

func valsToString(vals []int) []string {
	strs := make([]string, len(vals))
	for i,v := range(vals) {
		strs[i] = fmt.Sprintf("%v", v)
	}
	return strs
}

func (r *Statistics) WriteToCSV(path string) {
	f, e := os.Create(path)
	if e != nil {
		log.Fatal(e)
	}
	writer := csv.NewWriter(f)
	e = writer.WriteAll(r.Stats)
	if e != nil {
		log.Fatal(e)
	}
}