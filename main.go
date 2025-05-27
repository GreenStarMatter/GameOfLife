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
const pixelWidth int = 100
const pixelHeight int = 100

const (
	blackPixel = iota
	whitePixel
	redPixel
	bluePixel
	greenPixel

)

var pixelData = map[int]life{
	blackPixel: {
		pixelValue: blackPixel,
		lifeType: "dead",
		aggression: 0,
		affinity: 0,
		color: color.RGBA{0x0, 0x0, 0x0, 0xff},
		strength: 0,
		defense: 0,
		surroundingLifeNeeded: 10,
		minLifeToSustain: 0,
		maxLifeToSustain: 0,
	},
	whitePixel: {
		pixelValue: whitePixel,
		lifeType: "neutral",
		aggression: 5,
		affinity: 8,
		color: color.RGBA{0xff, 0xff, 0xff, 0xff},
		strength: 5,
		defense: 5,
		surroundingLifeNeeded: 2,
		minLifeToSustain: 2,
		maxLifeToSustain: 3,
	},
	redPixel: {
		pixelValue: redPixel,
		lifeType: "aggressive",
		aggression: 8,
		affinity: 1,
		color: color.RGBA{0xa0, 0x0, 0x0, 0xff},
		strength: 5,
		defense: 5,
		surroundingLifeNeeded: 2,
		minLifeToSustain: 1,
		maxLifeToSustain: 2,
	},
	bluePixel: {
		pixelValue: bluePixel,
		lifeType: "defensive",
		aggression: 2,
		affinity: 2,
		color: color.RGBA{0x0, 0x0, 0xc0, 0xff},
		strength: 5,
		defense: 5,
		surroundingLifeNeeded: 4,
		minLifeToSustain: 3,
		maxLifeToSustain: 5,
	},
	greenPixel: {
		pixelValue: greenPixel,
		lifeType: "grass",
		aggression: 0,
		affinity: 3,
		color: color.RGBA{0x0, 0xa0, 0x0, 0xff},
		strength: 0,
		defense: 0,
		surroundingLifeNeeded: 1,
		minLifeToSustain: 0,
		maxLifeToSustain: 8,
	},
}


type Game struct{}

type life struct {
	pixelValue int
	lifeType string
	aggression int //1-10
	affinity int //1-10
	color color.RGBA
	strength int //1-10
	defense int //1-10
	surroundingLifeNeeded int //1-8
	minLifeToSustain int //1-8
	maxLifeToSustain int //1-8
	
}

func New(ident int) *life {
	newLife := pixelData[ident]
	return &newLife
}


type lifeState struct{
	lifeSeeded bool
	currentState [pixelWidth * pixelHeight]life
	prevState [pixelWidth * pixelHeight]life
}


func (l *lifeState) seedLife() {
	for i:=0; i < len(l.currentState); i++ {
		l.currentState[i] = *New(rand.Intn(len(pixelData)))
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
		lifeCountMap, incomingAttackStrength := l.getSurroundingLifeInMatrix(x,y)
		if (l.prevState[i].lifeType != "dead") { //TODO: change all of these rules to depend on life specific needs
			if (lifeCountMap[l.prevState[i].pixelValue] < l.prevState[i].minLifeToSustain) {
				l.currentState[i] = *New(0)
			} else if (l.prevState[i].minLifeToSustain<=lifeCountMap[l.prevState[i].pixelValue]) && (l.prevState[i].maxLifeToSustain>=lifeCountMap[l.prevState[i].pixelValue]){
				l.currentState[i] = l.prevState[i]
			} else if (l.prevState[i].maxLifeToSustain<lifeCountMap[l.prevState[i].pixelValue]) {
				l.currentState[i] = *New(0)
			}
		} else if (l.prevState[i].lifeType == "dead") || (l.prevState[i].lifeType == "grass") {

			filtered := filterBySurroundingLifeNeeded(lifeCountMap)
			if (len(filtered) != 0) {
				//maxLifeKey, _ := maxKeyVal(filtered)
				//l.currentState[i] = *New(maxLifeKey)
				randLifeKey := getRandomKey(filtered)
				l.currentState[i] = *New(randLifeKey)

			}
		}
		if (l.prevState[i].lifeType != "dead") && (l.prevState[i].defense < sumValues(incomingAttackStrength)) {
			l.currentState[i] = *New(0)
		}
		
	} 
}

func (l *lifeState) getSurroundingLifeInMatrix(x,y int) (lifeCountMap, attackSumMap map[int]int) {
	//for a dead cell
	lifeCount := 0
	lifeCountMap = make(map[int]int)
	incomingAttackStrength := 0
	attackSumMap = make(map[int]int)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (x+i<0) || (y+j<0) || (x+i>=pixelWidth) || (y+j>=pixelHeight) {
			} else{
				vLifePos := convertMatrixToVector(x+i,y+j)
				//lifeCountMap := make(map[int]int)
				if (l.prevState[vLifePos].lifeType == l.currentState[vLifePos].lifeType ) || (l.prevState[vLifePos].lifeType == "grass") { //sustain on current life
					lifeCountMap[l.prevState[vLifePos].pixelValue]++
					lifeCount = lifeCount + 1
				}
				provokeLevel := 0
				if l.prevState[vLifePos].aggression != 0 {
					provokeLevel = rand.Intn(l.prevState[vLifePos].aggression)
				}				
				if (l.prevState[vLifePos].affinity - provokeLevel) < 0 {
					strengthCheck := 0
					if l.prevState[vLifePos].aggression != 0 {
						strengthCheck = rand.Intn(l.prevState[vLifePos].aggression)
					}
					incomingAttackStrength = incomingAttackStrength + strengthCheck
					attackSumMap[l.prevState[vLifePos].pixelValue] += incomingAttackStrength
				}
			}
		}
	}
	return lifeCountMap, attackSumMap
}

func filterBySurroundingLifeNeeded(lifeCounts map[int]int) map[int]int {
	filtered := make(map[int]int)
	for pixelType, count := range lifeCounts {
		if count >= pixelData[pixelType].surroundingLifeNeeded {
			filtered[pixelType] = count
		}
	}
	return filtered
}

func getRandomKey(m map[int]int) int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return -1 // or handle empty map
	}
	return keys[rand.Intn(len(keys))]
}

func maxKeyVal(t map[int]int) (maxKey, maxVal int) {
	maxKey = -1
	maxVal = -1

	for k, v := range t {
		if v > maxVal {
			maxVal = v
			maxKey = k
		}
	}

	return maxKey, maxVal
}

func sumValues(m map[int]int) int {
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

func (g *Game) Draw(screen *ebiten.Image) {//TODO: see about removing everything except renderLife(screen) here
	if !gameLife.lifeSeeded {
		//screen.Fill(pixelMap[whitePixel])
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
			pixelColor := gameLife.currentState[i].color
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