package main

type Player struct {
	Name string
}

func NewPlayer(name string) Player {
	return Player{Name: name}
}
