# Architecture

## Description
For the graphical interface we used Ebiten, a framework for making videogames in Go.

We defined structure types for every element in our game: 
  -Player struct
  -Ghost struct
  -Dot Type
  -Wall type
  -Player Sprite
  -Ghost Sprite
  -Game struct
  
The game struct contains common information for easy access during the progress of the game. For example, it stores the maze structure, the sprites, and other counting data like 

## Diagram architecture

## Functions

main(): creates a game instance and calls various initiliazer methods, which read the command arguments, read the maze, read sprite files, set the Ebiten window and initialise ghosts. Starts the game in "Start" mode.

Update(): checks if the game has started or not. Updates booleans for ghosts, dots and powerDots (eaten/notEaten) as well as managing key buttons, 
          checking if Pacman can move or there is a wall in between. ALso responsible for reapearing Pacman on the opposite side when out of a border.

Draw(): draws the game layout using Ebiten, which include walls, dots, power dots, ghosts, player, score, start and end texts.

Layout()
readArg()
readFont()
readSprites()
readMaze()

setWindowConfig()

initialiseGhots()

makeGhostEatable(): when a powerDot is eaten, this function is called as a goroutine. It gives Pacman the ability to eat ghosts and lasts 10 seconds.

valueIsInSlice()

relativePos()

nextIsWall()

drawDirection()

ghostSpriteDir()

playerSpriteDir()

restart()

checkError()
