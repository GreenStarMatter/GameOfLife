package main

import (
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/gameoflife/game"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Fill")
	if err := ebiten.RunGame(&game.Game{}); err != nil {
		log.Fatal(err)
	}
}