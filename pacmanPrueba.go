package main

import (
	"bufio"
	_ "image/gif"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

const spriteSize = 20
const spriteSizeSmaller = 15

type coord struct {
	posX int
	posY int
}

type playerSprite struct {
	right *ebiten.Image
	left  *ebiten.Image
	down  *ebiten.Image
	up    *ebiten.Image
}

type Game struct {
	maze          [][]string
	score         int
	lives         int
	numDots       int
	walls         []*coord
	wallSprite    *ebiten.Image
	player        coord
	playerSprites playerSprite
	playerDir     string
	ghosts        []*coord
	ghostSprite   *ebiten.Image
	points        []*coord
	pointSprite   *ebiten.Image
}

func (g *Game) Update() error {

	for _, k := range inpututil.PressedKeys() {
		if k == ebiten.KeyRight {
			g.player.posX += 1
			g.playerDir = "Right"
		} else if k == ebiten.KeyLeft {
			g.player.posX -= 1
			g.playerDir = "Left"
		} else if k == ebiten.KeyUp {
			g.player.posY -= 1
			g.playerDir = "Up"
		} else if k == ebiten.KeyDown {
			g.player.posY += 1
			g.playerDir = "Down"
		} else if k == ebiten.KeyEscape {
			g.lives = 0
		}

		if relativePos(g.player.posX) > screenWidth {
			g.player.posX = -1
		}

		if relativePos(g.player.posX) < 0-playerWidth {
			g.player.posX = len(g.maze[0])
		}

		if g.player.posY <= 0 {
			g.player.posY = 0
		}

		if relativePos(g.player.posY) >= screenHeight-playerHeight {
			g.player.posY = len(g.maze) - 1
		}
	}

	g.moveGhosts()

	time.Sleep(200 * time.Millisecond)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//draw pacman
	opPlayer := &ebiten.DrawImageOptions{}
	playerH := spriteSize / playerHeight //CAMBIAR NOMBRE VARIABLE
	playerW := spriteSize / playerWidth  //CAMBIAR NOMBRE VARIABLE
	opPlayer.GeoM.Scale(playerW, playerH)
	opPlayer.GeoM.Translate(relativePos(g.player.posX), relativePos(g.player.posY))

	screen.DrawImage(g.playerSpriteDir(), opPlayer)

	//draw walls
	opWall := &ebiten.DrawImageOptions{}
	wallH := spriteSize / wallHeight //CAMBIAR NOMBRE VARIABLE
	wallW := spriteSize / wallWidth  //CAMBIAR NOMBRE VARIABLE
	//var wallPosX float64
	//var wallPosY float64
	//var prevWallPosX float64 = 0
	//var prevWallPosY float64 = 0

	for _, wall := range g.walls {
		opWall.GeoM.Scale(wallW, wallH)
		opWall.GeoM.Translate(relativePos(wall.posX), relativePos(wall.posY))
		screen.DrawImage(g.wallSprite, opWall)
		opWall.GeoM.Reset()

		//prevWallPosX = wall.posX
		//prevWallPosY = wall.posY
	}

	/*for i := 0; i < len(g.maze); i++ {
		//fmt.Printf("%s\n", g.maze[i])
		for j := 0; j < len(g.maze[0]); j++ {
			switch g.maze[i][j] {
			case "#":
				wallPosX = float64(j * spriteSize)
				wallPosY = float64(i * spriteSize)

				opWall.GeoM.Translate(wallPosX-prevWallPosX, wallPosY-prevWallPosY)
				screen.DrawImage(g.wallSprite, opWall)

				prevWallPosX = wallPosX
				prevWallPosY = wallPosY
			}
		}
	}*/

	//drawn ghosts
	opGhost := &ebiten.DrawImageOptions{}
	ghostH := spriteSize / ghostHeight //CAMBIAR NOMBRE VARIABLE
	ghostW := spriteSize / ghostWidth  //CAMBIAR NOMBRE VARIABLE

	//var prevGhostPosX float64 = 0
	//var prevGhostPosY float64 = 0

	for _, ghost := range g.ghosts {
		//fmt.Printf("%v, %v\n", ghost.posX, ghost.posY)
		opGhost.GeoM.Scale(ghostW, ghostH)
		opGhost.GeoM.Translate(relativePos(ghost.posX), relativePos(ghost.posY))
		screen.DrawImage(g.ghostSprite, opGhost)
		opGhost.GeoM.Reset()

		//prevGhostPosX = ghost.posX
		//prevGhostPosY = ghost.posY
	}

	//drawn points
	opPoint := &ebiten.DrawImageOptions{}
	pointH := spriteSize / pointSize //CAMBIAR NOMBRE VARIABLE
	pointW := spriteSize / pointSize //CAMBIAR NOMBRE VARIABLE

	for _, point := range g.points {
		opPoint.GeoM.Scale(pointW, pointH)
		opPoint.GeoM.Translate(relativePos(point.posX), relativePos(point.posY))
		screen.DrawImage(g.pointSprite, opPoint)
		opPoint.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(screenWidth), int(screenHeight)
}

func main() {
	playerRightImg, _, err := ebitenutil.NewImageFromFile("assets/pacman_derecha.png")
	checkError(err, "Load player image error")
	playerLeftImg, _, err := ebitenutil.NewImageFromFile("assets/pacman_izquierda.png")
	checkError(err, "Load player image error")
	playerUpImg, _, err := ebitenutil.NewImageFromFile("assets/pacman_arriba.png")
	checkError(err, "Load player image error")
	playerDownImg, _, err := ebitenutil.NewImageFromFile("assets/pacman_abajo.png")
	checkError(err, "Load player image error")
	wallImg, _, err := ebitenutil.NewImageFromFile("assets/wall.png")
	checkError(err, "Load wall image error")
	ghostImg, _, err := ebitenutil.NewImageFromFile("assets/200.gif")
	checkError(err, "Load ghost image error")
	pointImg, _, err := ebitenutil.NewImageFromFile("assets/punto.png")
	checkError(err, "Load point image error")

	g := &Game{}
	g.wallSprite = wallImg
	g.playerSprites.right = playerRightImg
	g.playerSprites.left = playerLeftImg
	g.playerSprites.up = playerUpImg
	g.playerSprites.down = playerDownImg
	g.ghostSprite = ghostImg
	g.pointSprite = pointImg

	g.playerDir = "Right"
	g.lives = 1

	g.readMaze("Maze.txt")

	screenWidth = float64(len(g.maze[0]) * spriteSize)
	screenHeight = float64(len(g.maze) * spriteSize)

	ebiten.SetWindowSize(int(screenWidth), int(screenHeight))
	ebiten.SetWindowTitle("Pacman")

	//g.playPosX = screenWidth/2 - playerWidth/2
	//g.playPosY = screenHeight/2 - playerHeight/2

	err = ebiten.RunGame(g)
	checkError(err, "Run game error")
}

//Funciones fuera del loop del juego

func (g *Game) readMaze(fileName string) {
	file, err := os.Open(fileName)
	checkError(err, "Load maze error")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linea := strings.Split(scanner.Text(), "")
		g.maze = append(g.maze, linea)
	}

	file.Close()

	for i := 0; i < len(g.maze); i++ {
		for j := 0; j < len(g.maze[0]); j++ {
			switch g.maze[i][j] {
			case "#":
				//wallPosX := j * spriteSize
				//wallPosY := i * spriteSize

				g.walls = append(g.walls, &coord{j, i})
			case "P":
				//playerPosX := j * spriteSize
				//playerPosY := i * spriteSize

				g.player = coord{j, i}
			case "G":
				//ghostPosX := j * spriteSize
				//ghostPosY := i * spriteSize

				g.ghosts = append(g.ghosts, &coord{j, i})
			case ".":
				//pointPosX := j * spriteSize
				//pointPosY := i * spriteSize

				g.points = append(g.points, &coord{j, i})
			}
		}
	}
	/*for _, wall := range g.walls {
		fmt.Printf("%v, %v\n", wall.posX, wall.posY)
	}*/
}

func relativePos(coordinate int) (coordRelative float64) {
	return float64(coordinate * spriteSize)
}

func (g *Game) makeMove(oldCol, oldRow int, dir string) (newCol, newRow int) {
	newRow, newCol = oldRow, oldCol

	switch dir {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(g.maze) - 1
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(g.maze) {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(g.maze[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(g.maze[0]) - 1
		}
	}

	if g.maze[newRow][newCol] == "#" {
		newRow = oldRow
		newCol = oldCol
	}

	return
}

func drawDirection() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return move[dir]
}

func (g *Game) moveGhosts() {
	for _, ghost := range g.ghosts {
		dir := drawDirection()
		ghost.posX, ghost.posY = g.makeMove(ghost.posX, ghost.posY, dir)
	}
}

func (g *Game) playerSpriteDir() *ebiten.Image {
	switch g.playerDir {
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

func checkError(err error, message string) {
	if err != nil {
		log.Printf("[%s]", message)
		panic(err)
	}
}
