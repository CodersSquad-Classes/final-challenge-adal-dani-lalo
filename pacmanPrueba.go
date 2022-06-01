package main

import (
	"bufio"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var screenWidth float64
var screenHeight float64

const playerWidth float64 = 50
const playerHeight float64 = 53
const wallWidth float64 = 300
const wallHeight float64 = 300
const ghostWidth float64 = 209
const ghostHeight float64 = 152
const pointSize float64 = 100

const spriteSize = 20
const spriteSizeSmaller = 15

type coord struct {
	posX float64
	posY float64
}

type Game struct {
	maze         [][]string
	walls        []*coord
	wallSprite   *ebiten.Image
	player       coord
	playerSprite *ebiten.Image
	playerDir    string
	ghosts       []*coord
	ghostSprite  *ebiten.Image
	points       []*coord
	pointSprite  *ebiten.Image
}

func (g *Game) Update() error {

	for _, k := range inpututil.PressedKeys() {
		if k == ebiten.KeyRight {
			g.player.posX += 3
			g.playerDir = "Right"
		} else if k == ebiten.KeyLeft {
			g.player.posX -= 3
			g.playerDir = "Left"
		} else if k == ebiten.KeyUp {
			g.player.posY -= 3
			g.playerDir = "Up"
		} else if k == ebiten.KeyDown {
			g.player.posY += 3
			g.playerDir = "Down"
		}

		if g.player.posX > screenWidth {
			g.player.posX = 0 - playerWidth
		}

		if g.player.posX < 0-playerWidth {
			g.player.posX = screenWidth
		}

		if g.player.posY <= 0 {
			g.player.posY = 0
		}

		if g.player.posY >= screenHeight-playerHeight {
			g.player.posY = screenHeight - playerHeight
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//draw pacman
	opPlayer := &ebiten.DrawImageOptions{}
	playerH := spriteSizeSmaller / playerHeight //CAMBIAR NOMBRE VARIABLE
	playerW := spriteSizeSmaller / playerWidth  //CAMBIAR NOMBRE VARIABLE
	opPlayer.GeoM.Scale(playerW, playerH)
	opPlayer.GeoM.Translate(g.player.posX, g.player.posY)
	screen.DrawImage(g.playerSprite, opPlayer)

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
		opWall.GeoM.Translate(wall.posX, wall.posY)
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
	ghostH := spriteSizeSmaller / ghostHeight //CAMBIAR NOMBRE VARIABLE
	ghostW := spriteSizeSmaller / ghostWidth  //CAMBIAR NOMBRE VARIABLE

	//var prevGhostPosX float64 = 0
	//var prevGhostPosY float64 = 0

	for _, ghost := range g.ghosts {
		//fmt.Printf("%v, %v\n", ghost.posX, ghost.posY)
		opGhost.GeoM.Scale(ghostW, ghostH)
		opGhost.GeoM.Translate(ghost.posX, ghost.posY)
		screen.DrawImage(g.ghostSprite, opGhost)
		opGhost.GeoM.Reset()

		//prevGhostPosX = ghost.posX
		//prevGhostPosY = ghost.posY
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(screenWidth), int(screenHeight)
}

func main() {
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/pacman.png")
	checkError(err, "Load player image error")
	wallImg, _, err := ebitenutil.NewImageFromFile("assets/wall.png")
	checkError(err, "Load wall image error")
	ghostImg, _, err := ebitenutil.NewImageFromFile("assets/Alien1.png")
	checkError(err, "Load ghost image error")
	pointImg, _, err := ebitenutil.NewImageFromFile("assets/punto.png")
	checkError(err, "Load point image error")

	g := &Game{}
	g.wallSprite = wallImg
	g.playerSprite = playerImg
	g.ghostSprite = ghostImg
	g.pointSprite = pointImg

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
				wallPosX := float64(j * spriteSize)
				wallPosY := float64(i * spriteSize)

				g.walls = append(g.walls, &coord{wallPosX, wallPosY})
			case "P":
				playerPosX := float64(j * spriteSize)
				playerPosY := float64(i * spriteSize)

				g.player = coord{playerPosX, playerPosY}
			case "G":
				ghostPosX := float64(j * spriteSize)
				ghostPosY := float64(i * spriteSize)

				g.ghosts = append(g.ghosts, &coord{ghostPosX, ghostPosY})
			}
		}
	}
	/*for _, wall := range g.walls {
		fmt.Printf("%v, %v\n", wall.posX, wall.posY)
	}*/
}

func checkError(err error, message string) {
	if err != nil {
		log.Printf("[%s]", message)
		panic(err)
	}
}
