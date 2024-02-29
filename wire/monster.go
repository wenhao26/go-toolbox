package main

type Monster struct {
	Name string
}

func NewMonster() Monster {
	return Monster{Name: "BigBigMan"}
}