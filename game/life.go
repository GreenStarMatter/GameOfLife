package game

import "image/color"

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
		surroundingLifeNeeded: 3,
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


type lifeState struct{
	lifeSeeded bool
	currentState [pixelWidth * pixelHeight]life
	prevState [pixelWidth * pixelHeight]life
}