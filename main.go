package main

import (
	"image/color"
	"log"
	"math/rand"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

var gameLife lifeState
var stepSize int = 15
var frameCounter int = 0
const pixelWidth int = 300
const pixelHeight int = 300

const (
	blackPixel = iota
	whitePixel
	redPixel
	bluePixel
	greenPixel

)

var pixelMap = map[int]color.RGBA{
	whitePixel: color.RGBA{0xff, 0xff, 0xff, 0xff},
	blackPixel: color.RGBA{0x0, 0x0, 0x0, 0xff},
	redPixel: color.RGBA{0xff, 0x0, 0x0, 0xff},
	bluePixel: color.RGBA{0x0, 0x0, 0xff, 0xff},
	greenPixel: color.RGBA{0x0, 0xff, 0x0, 0xff},

}
type Game struct{}

type lifeState struct{
	lifeSeeded bool
	currentState [pixelWidth * pixelHeight]int
	prevState [pixelWidth * pixelHeight]int
}


func (l *lifeState) seedLife() {
	for i:=0; i < len(l.currentState); i++ {
		l.currentState[i] = rand.Intn(2)
	}

}

func (g *Game) Update() error {
	
	gameLife.prevState = gameLife.currentState
	if !gameLife.lifeSeeded {
		fmt.Println("Seeding")
		gameLife.seedLife()
	} else {
		if frameCounter > stepSize {
			gameLife.updateLife()
			frameCounter = 0
		} else {
			frameCounter = frameCounter + 1
		}
	}
	return nil
}

func (l *lifeState) updateLife() {
	for i:=0; i < len(l.currentState); i++ {
		x,y := convertVectorToMatrix(i, pixelWidth * pixelHeight)
		lifeCount := l.getSurroundingLifeInMatrix(x,y)
		if (l.prevState[i] == 1) {
			if (lifeCount<2) {
				l.currentState[i] = 0///Live cell, Fewer than 2 Live Neighbors dies
			} else if (lifeCount==2) || (lifeCount==3){
				l.currentState[i] = 1///Live cell, 2 or 3 live on to next generation
			} else if (lifeCount>3) {
				l.currentState[i] = 0///Live Cell, More than 3 neighbors dies
			}
		} else if l.prevState[i] == 0 {
			if (lifeCount==3) {
				l.currentState[i] = 1///Dead cell, Exactly 3 live neighbors becomes live cell
			}
		}
		
	} 
}

func (l *lifeState) getSurroundingLifeInMatrix(x,y int) int {
	lifeCount := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (x+i<0) || (y+j<0) || (x+i>=pixelWidth) || (y+j>=pixelHeight) {
			} else {
				vLifePos := convertMatrixToVector(x+i,y+j)
				lifeCount = lifeCount + l.prevState[vLifePos]
			}
		}
	}
	return lifeCount
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !gameLife.lifeSeeded {
		screen.Fill(pixelMap[whitePixel])
		gameLife.lifeSeeded = true
		fmt.Println("Finished Seeding")
	} else {
		renderLife(screen)
	}
}

func renderLife(screen *ebiten.Image) {
	for i:=0; i<len(gameLife.currentState); i++ {
		//if gameLife.currentState[i] != gameLife.prevState[i] {
			x,y := convertVectorToMatrix(i, pixelWidth * pixelHeight)
			pixelColor := pixelMap[gameLife.currentState[i]]
			screen.Set(x,y,pixelColor)
		//}
	}
}

func convertVectorToMatrix(pos int, vlen int)  (x,y int) {
	y = pos / pixelWidth //position divided by x width
	x = pos % pixelWidth //position modulo all the widths
	return x,y
}

func convertMatrixToVector(x,y int)  int {
	return y * pixelWidth + x //assuming y is the wrapping dimension
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return pixelWidth, pixelHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Fill")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}