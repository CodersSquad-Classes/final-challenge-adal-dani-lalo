# Architecture

## Description
This is a Pac-Man game imitation written in Go and built with Ebiten game engine.
It is worth to mention that we made use of goroutines in various ways in this project, like on a separated movement control for each ghost. 

## GUI
For the graphical interface we used Ebiten, a framework for making videogames in Go.
![ebiten](https://user-images.githubusercontent.com/78662124/172483724-8182bcc2-746f-4986-bf9a-ad3c887d58cd.png)

## Diagram architecture
The overall flow of the Pac-Man game consists of the following threads.

![pacman_architecture drawio](https://user-images.githubusercontent.com/78662124/172498190-7ef9e68b-fda0-467a-95dc-c58c9800ff6c.png)
There is a main thread, from where everything is called. Then, we have two Ebiten based thread for the functions Update and Draw.
Then, we use separate threads (Goroutines) for each one of the ghosts. That is why, it is multithreaded.

NOTE: we also use an extra goroutine for the special Power Dot that only lasts 10 seconds.


## Data structures
We defined structure types for every element in our game: 
       -Player struct
       -Ghost struct
       -Dot Type
       -Wall type
       -Player Sprite
       -Ghost Sprite
       -Game struct
  
The game struct contains common information for easy access during the progress of the game. For example, it stores the maze structure, the sprites, and other counting data like 

## Main functions

**main()**: creates a game instance and calls various initiliazer methods. Starts the game in "Start" mode. 
The sere the methods called on the game instance.

    -readArg():  reads command arguments
    -readFont(): loads the used font.
    -readSprites(): loads the sprites and assigns them to variables ans structs.
    -readMaze(): reads the txt defining the maze. Walls, dots and other elements are represented by symbols in the txt file.
    -setWindowConfig(): set up Ebiten display window.
    -initialiseGhots(): initializes each ghost's own thread (goroutine).
  
**Update()**: this is the gameloop, whichs is in charge of updating what the data structs and boolean frlags.

    -It checks if the game has started or not.
    -Updates booleans for ghosts, dots and powerDots (eaten/notEaten).
    -Managing key buttons, checking if Pacman can move or there is a wall in between.
    -Also responsible for reapearing Pacman on the opposite side when out of a border.

**Layout()**: this is an Ebiten based function which is in charge of setting up the window where the game will display.

**Draw()**: draws the game layout using Ebiten, which includes:

     -Walls
     -Dots
     -Power dots
     -Ghosts
     -Player
     -Score
     -Start and end texts.


#### Goroutine ghost behaviour
Instead of controlling all the ghosts in a single thread, we span as many goroutines as there are ghosts, so each ghost movement is managed in a separate goroutine. THe function that is called as a goroutine is the following:

       - ghostFunctionality()
      
 This goroutine function makes frequent calls to a helper function:
       
       - getRandomDirection(): get a random direction for the next move.
 
### Helper Functions
These are a group of functions which constantly check wheter a certain boolean flag is true or not.

       - valueIsInSlice(): check if a certain value is already in a given slice.
       - nextIsWall(): checks if the next position Pac-Man wants to move is a wall.
       - scaleCoord(): scale maze cell coordinate to screen pixel coordinate.
       - restart(): resets various player and ghosts parameters, as well as global variables like score.
       - checkError(): check wheter something throws an error. If so, the process is terminated.


### Game Loop Functions

       -makeGhostEatable(): when a powerDot is eaten, this function is called as a goroutine. It gives Pacman the ability to eat ghosts and lasts 10             seconds.

       -ghostSpriteDir(): returns the appropriate sprite based on the current ghost direction or state (in case it is vulnerable or normal).

       -playerSpriteDir(): returns the appropriate sprite based on the current player direction and state (in case it is alive or dead).
