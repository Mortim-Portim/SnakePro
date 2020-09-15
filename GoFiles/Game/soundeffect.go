package Game

import (
	"marvin/SnakeProGo/Utilities"
	//"github.com/hajimehoshi/ebiten"
	"fmt"
	//"time"
	"math/rand"
)

type SoundEffect struct {
	sounds []*Utilities.AudioPlayer
}

func LoadSoundEffect(folder string) (s *SoundEffect) {
	s = &SoundEffect{}
	s.sounds = make([]*Utilities.AudioPlayer,0)
	for i := 0; i < 20; i++ {
		file := fmt.Sprintf("%s/%v.wav", folder, i+1)
		fmt.Println("Loading: ", file)
		player, err := Utilities.NewPlayer(file)
		if err != nil {
			break
		}
		s.sounds = append(s.sounds, player)
	}
	if len(s.sounds) == 0 {
		panic("Could not load SoundEffects")
	}
	return
}

func (s *SoundEffect) PR() {
	idx := rand.Intn(len(s.sounds))
	s.sounds[idx].Play()
}
