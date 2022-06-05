package main

import (
	"bufio"
	"errors"
	"fmt"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var screenWidth float64
var screenHeight float64

const playerWidth float64 = 111
const playerHeight float64 = 105
const wallWidth float64 = 300
const wallHeight float64 = 300
const ghostWidth float64 = 112
const ghostHeight float64 = 106
const pointSize float64 = 100

var arcadeFont font.Face
var fontSize = 20.0

const spriteSize = 20

type playerType struct {
	posX      int
	posY      int
	direction string
	lives     int
}

type ghostType struct {
	posX      int
	posY      int
	direction string
	eaten     bool
}

type dotType struct {
	posX  int
	posY  int
	eaten bool
}

type wallType struct {
	posX int
	posY int
}

type playerSprite struct {
	right *ebiten.Image
	left  *ebiten.Image
	down  *ebiten.Image
	up    *ebiten.Image
	death *ebiten.Image
}

type Game struct {
	maze          [][]string
	score         int
	numDots       int
	numGhosts     int
	walls         []*wallType
	wallSprite    *ebiten.Image
	player        playerType
	playerSprites playerSprite
	ghosts        []*ghostType
	ghostSprite   *ebiten.Image
	dots          []*dotType
	dotSprite     *ebiten.Image
	powerDots     []*dotType
}

func (g *Game) Update() error {

	/*if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.posX += 1
		g.player.direction = "Right"
	}*/

	for _, k := range inpututil.PressedKeys() {
		if k == ebiten.KeyRight && !g.nextIsWall("Right") {
			g.player.posX += 1
			g.player.direction = "Right"
		} else if k == ebiten.KeyLeft && !g.nextIsWall("Left") {
			g.player.posX -= 1
			g.player.direction = "Left"
		} else if k == ebiten.KeyUp && !g.nextIsWall("Up") {
			g.player.posY -= 1
			g.player.direction = "Up"
		} else if k == ebiten.KeyDown && !g.nextIsWall("Down") {
			g.player.posY += 1
			g.player.direction = "Down"
		} else if k == ebiten.KeyEscape {
			g.player.lives = 0
		}

		if g.player.posX > len(g.maze[0]) {
			g.player.posX = 0
		}

		if g.player.posX < 0 {
			g.player.posX = len(g.maze[0])
		}
	}

	//g.moveGhosts()

	for _, ghost := range g.ghosts {
		if g.player.posX == ghost.posX && g.player.posY == ghost.posY {
			g.player.lives = 0
		}
	}

	for _, dot := range g.dots {
		if g.player.posX == dot.posX && g.player.posY == dot.posY && !dot.eaten {
			g.score += 10
			g.numDots -= 1
			dot.eaten = true
		}
	}

	for _, powerDot := range g.powerDots {
		if g.player.posX == powerDot.posX && g.player.posY == powerDot.posY && !powerDot.eaten {
			g.score += 50
			g.numDots -= 1
			powerDot.eaten = true
		}
	}

	time.Sleep(200 * time.Millisecond)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//draw walls
	opWall := &ebiten.DrawImageOptions{}
	wallHScale := spriteSize / wallHeight
	wallWScale := spriteSize / wallWidth

	for _, wall := range g.walls {
		opWall.GeoM.Scale(wallWScale, wallHScale)
		opWall.GeoM.Translate(relativePos(wall.posX), relativePos(wall.posY))
		screen.DrawImage(g.wallSprite, opWall)
		opWall.GeoM.Reset()
	}

	//drawn dots
	opDot := &ebiten.DrawImageOptions{}
	dotHScale := spriteSize / pointSize
	dotWScale := spriteSize / pointSize

	for _, dot := range g.dots {
		if dot.eaten {
			continue
		}
		opDot.GeoM.Scale(dotWScale, dotHScale)
		opDot.GeoM.Translate(relativePos(dot.posX), relativePos(dot.posY))
		screen.DrawImage(g.dotSprite, opDot)
		opDot.GeoM.Reset()
	}

	//drawn power dots
	opPowerDot := &ebiten.DrawImageOptions{}
	powerDotHScale := spriteSize / pointSize * 3
	powerDotWScale := spriteSize / pointSize * 3

	for _, powerDot := range g.powerDots {
		if powerDot.eaten {
			continue
		}
		opPowerDot.GeoM.Scale(powerDotWScale, powerDotHScale)
		opPowerDot.GeoM.Translate(relativePos(powerDot.posX)-spriteSize, relativePos(powerDot.posY)-spriteSize)
		screen.DrawImage(g.dotSprite, opPowerDot)
		opPowerDot.GeoM.Reset()
	}

	//drawn ghosts
	opGhost := &ebiten.DrawImageOptions{}
	ghostHScale := spriteSize / ghostHeight
	ghostWScale := spriteSize / ghostWidth

	for _, ghost := range g.ghosts {
		opGhost.GeoM.Scale(ghostWScale, ghostHScale)
		opGhost.GeoM.Translate(relativePos(ghost.posX), relativePos(ghost.posY))
		screen.DrawImage(g.ghostSprite, opGhost)
		opGhost.GeoM.Reset()
	}

	//draw player
	opPlayer := &ebiten.DrawImageOptions{}
	playerHScale := spriteSize / playerHeight
	playerWScale := spriteSize / playerWidth
	opPlayer.GeoM.Scale(playerHScale, playerWScale)
	opPlayer.GeoM.Translate(relativePos(g.player.posX), relativePos(g.player.posY))

	screen.DrawImage(g.playerSpriteDir(), opPlayer)

	//draw score
	str := "SCORE: " + strconv.Itoa(g.score)
	text.Draw(screen, str, arcadeFont, 0, len(g.maze)*spriteSize+spriteSize, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(screenWidth), int(screenHeight)
}

func main() {
	g := &Game{}

	g.readArg()
	g.readFont()
	g.readSprites()
	g.readMaze("Maze.txt")
	g.setWindowConfig()
	g.initialiseGhots()

	err := ebiten.RunGame(g)
	checkError(err, "Run game error")
}

func (g *Game) initialiseGhots() {
	for _, ghost := range g.ghosts {
		go g.ghostFuncionability(ghost)
	}
}

func (g *Game) ghostFuncionability(ghost *ghostType) {
	for {
		for {
			newY, newX := ghost.posY, ghost.posX

			switch ghost.direction {
			case "Up":
				newY = newY - 1

			case "Down":
				newY = newY + 1

			case "Right":
				newX = newX + 1
				if newX == len(g.maze[0]) {
					newX = 0
				}
			case "Left":
				newX = newX - 1
				if newX < 0 {
					newX = len(g.maze[0]) - 1
				}
			}

			if g.maze[newY][newX] == "#" {
				ghost.direction = drawDirection()
			} else {
				ghost.posX = newX
				ghost.posY = newY
				break
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}

//MAIN FUNCTIONS

func (g *Game) readArg() {
	if len(os.Args) == 3 {
		if os.Args[1] == "--enemies" {
			numGhosts, err := strconv.Atoi(os.Args[2])
			checkError(err, "Expected arguments: --enemies 3")
			if numGhosts < 1 || numGhosts > 9 {
				err = errors.New("Incorrect number of enemies")
				checkError(err, "Number of enemies can't be lower than 1 or higher than 9")
			} else {
				g.numGhosts = numGhosts
			}
		} else {
			err := errors.New("Incorrect arguments")
			checkError(err, "Expected arguments: --enemies 3")
		}
	} else {
		err := errors.New("Incorrect arguments")
		checkError(err, "Expected arguments: --enemies 3")
	}
}

func (g *Game) readFont() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	checkError(err, "Error loading font")

	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	checkError(err, "Error loading font")
}

func (g *Game) readSprites() {
	playerRightImg, _, err := ebitenutil.NewImageFromFile("assets/pacman_derecha.png")
	checkError(err, "Load player image error")
	playerLeftImg, _, err := ebitenutil.NewImageFromFile("assets/pacman_izquierda.png")
	checkError(err, "Load player image error")
	playerUpImg, _, err := ebitenutil.NewImageFromFile("assets/pacman_arriba.png")
	checkError(err, "Load player image error")
	playerDownImg, _, err := ebitenutil.NewImageFromFile("assets/pacman_abajo.png")
	checkError(err, "Load player image error")
	playerDeathImg, _, err := ebitenutil.NewImageFromFile("assets/wall.png")
	checkError(err, "Load player image error")
	wallImg, _, err := ebitenutil.NewImageFromFile("assets/wall.png")
	checkError(err, "Load wall image error")
	ghostImg, _, err := ebitenutil.NewImageFromFile("assets/ghost_derecha.png")
	checkError(err, "Load ghost image error")
	dotImg, _, err := ebitenutil.NewImageFromFile("assets/punto.png")
	checkError(err, "Load point image error")

	g.wallSprite = wallImg
	g.playerSprites.right = playerRightImg
	g.playerSprites.left = playerLeftImg
	g.playerSprites.up = playerUpImg
	g.playerSprites.down = playerDownImg
	g.playerSprites.death = playerDeathImg
	g.ghostSprite = ghostImg
	g.dotSprite = dotImg
}

func (g *Game) readMaze(fileName string) {
	file, err := os.Open(fileName)
	checkError(err, "Load maze error")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linea := strings.Split(scanner.Text(), "")
		g.maze = append(g.maze, linea)
	}

	file.Close()

	var ghostPositions []int
	for i := 0; i < g.numGhosts; i++ {
		rand.Seed(time.Now().UnixNano())
		ghostPos := rand.Intn(9)
		if valueIsInSlice(ghostPos, ghostPositions) {
			i -= 1
		} else {
			fmt.Print(ghostPos)
			ghostPositions = append(ghostPositions, ghostPos)
		}
	}
	indexGhosts := 0

	for i := 0; i < len(g.maze); i++ {
		for j := 0; j < len(g.maze[0]); j++ {
			switch g.maze[i][j] {
			case "#":
				g.walls = append(g.walls, &wallType{j, i})
			case "P":
				g.player = playerType{j, i, "Right", 1}
			case "G":
				if valueIsInSlice(indexGhosts, ghostPositions) {
					g.ghosts = append(g.ghosts, &ghostType{j, i, "Up", false})
				}
				indexGhosts += 1
			case ".":
				g.numDots += 1
				g.dots = append(g.dots, &dotType{j, i, false})
			case "O":
				g.numDots += 1
				g.powerDots = append(g.powerDots, &dotType{j, i, false})
			}
		}
	}
}

func (g *Game) setWindowConfig() {
	screenWidth = float64(len(g.maze[0]) * spriteSize)
	screenHeight = float64(len(g.maze)*spriteSize + spriteSize)

	ebiten.SetWindowSize(int(screenWidth), int(screenHeight))
	ebiten.SetWindowTitle("Pacman")
}

//GAMELOOP FUNCTIONS

func valueIsInSlice(x int, slice []int) bool {
	for _, y := range slice {
		if y == x {
			return true
		}
	}
	return false
}

func relativePos(coordinate int) (coordRelative float64) {
	return float64(coordinate * spriteSize)
}

func (g *Game) nextIsWall(dir string) bool {
	newY, newX := g.player.posY, g.player.posX

	switch dir {
	case "Up":
		newY = newY - 1
	case "Down":
		newY = newY + 1
	case "Right":
		newX = newX + 1
	case "Left":
		newX = newX - 1
	}

	if newX > len(g.maze[0]) || newX < 0 { //return false if entity entered the tunel
		return false
	} else if g.maze[newY][newX] == "#" {
		return true
	} else {
		return false
	}
}

func (g *Game) moveGhosts() {
	for _, ghost := range g.ghosts {
		g.makeMove(ghost)
	}
}

//func (g *Game) makeMove(oldCol, oldRow int, dir string) (newCol, newRow int) {
func (g *Game) makeMove(entity *ghostType) {
	newY, newX := entity.posY, entity.posX

	switch entity.direction {
	case "Up":
		newY = newY - 1
		/*if newRow < 0 {
			newRow = len(g.maze) - 1
		}*/
	case "Down":
		newY = newY + 1
		/*if newRow == len(g.maze) {
			newRow = 0
		}*/
	case "Right":
		newX = newX + 1
		if newX == len(g.maze[0]) {
			newX = 0
		}
	case "Left":
		newX = newX - 1
		if newX < 0 {
			newX = len(g.maze[0]) - 1
		}
	}

	if g.maze[newY][newX] == "#" {
		entity.direction = drawDirection()
		g.makeMove(entity)
	} else {
		entity.posX = newX
		entity.posY = newY
	}
}

func drawDirection() string {
	rand.Seed(time.Now().UnixNano())
	dir := rand.Intn(4)
	move := map[int]string{
		0: "Up",
		1: "Down",
		2: "Right",
		3: "Left",
	}
	return move[dir]
}

func (g *Game) playerSpriteDir() *ebiten.Image {
	if g.player.lives == 0 {
		return g.playerSprites.death
	} else {
		switch g.player.direction {
		case "Right":
			return g.playerSprites.right
		case "Left":
			return g.playerSprites.left
		case "Up":
			return g.playerSprites.up
		case "Down":
			return g.playerSprites.down
		default:
			return g.playerSprites.right
		}
	}
}

func checkError(err error, message string) {
	if err != nil {
		log.Printf("[%s]", message)
		panic(err)
	}
}
