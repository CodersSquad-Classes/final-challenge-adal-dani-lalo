[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-c66648af7eb3fe8bc4f294546bfd86ef473780cde1dea487d3c4ff354943c9ae.svg)](https://classroom.github.com/online_ide?assignment_repo_id=7915600&assignment_repo_type=AssignmentRepo)

Link to video:
https://youtu.be/5gHEM-aDI9Q

##Team:

  -Adalberto Rodríguez
  -Daniel Díaz
  -Eduardo Legorreta

Multithreaded Pacman Game - (single-node)
=========================================
![Screenshot from 2022-06-07 22-17-18](https://user-images.githubusercontent.com/78662124/172524239-8bc170c7-e9a5-4e16-889a-cb967cf1dde9.png)


## Description
This is a Pac-Man game imitation written in Go and built with Ebiten game engine.
It is worth to mention that we made use of goroutines in various ways in this project, like on a separated movement control for each ghost.

## How to run
You can use Makefile to build the program. Just do:

  make
  
And then just do:

  make test
  
The game will start with 4 enemies.

#### Specify enemies amount
You need to run the executable. For example, if you want 7 enemies, just do:

  ./pacman --enemies 7
  
In order to clean, just do:

  make clean
