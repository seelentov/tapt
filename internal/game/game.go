package game

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"tapt/pkg/api"
	"time"

	"github.com/eiannone/keyboard"
)

type Game struct {
	score     int
	lives     int
	level     int
	inputText string
}

func NewGame() *Game {
	return &Game{
		score: 0,
		lives: 10,
		level: 0,
	}
}

func (g *Game) clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (g *Game) printHeader() {
	fmt.Println("Press ESC to quit")
	fmt.Printf("Score: %d | Lives: %d | Level: %d\n\n", g.score, g.lives, g.level+1)
}

func (g *Game) printTextWithCursor(targetText string, cursorPos int) {
	g.clearScreen()
	g.printHeader()

	fmt.Println("Type the following text:")
	fmt.Println()

	fmt.Print("\033[32m")
	fmt.Print(targetText[:cursorPos])
	fmt.Print("\033[0m")

	fmt.Print(targetText[cursorPos:])
	fmt.Println()
}

func (g *Game) Run() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for g.lives > 0 {
		targetText, err := api.GetText()
		if err != nil {
			continue
		}

		cursorPos := 0
		g.inputText = ""

		for cursorPos < len(targetText) && g.lives > 0 {
			g.printTextWithCursor(targetText, cursorPos)

			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}

			if key == keyboard.KeyEsc {
				return
			}

			if targetText[cursorPos] == ' ' && key == keyboard.KeySpace {
				cursorPos++
				g.inputText += " "
				g.score++
				continue
			}

			if char > 0 && char == rune(targetText[cursorPos]) {
				cursorPos++
				g.inputText += string(char)
				g.score++
			} else if key != keyboard.KeySpace {
				g.lives--
			}
		}

		if cursorPos == len(targetText) {
			g.level++
			fmt.Println("\nCorrect! Level completed!")
			g.lives++
			if g.lives > 10 {
				g.lives = 10
			}
			time.Sleep(500 * time.Millisecond)
		}
	}

	g.clearScreen()
	if g.lives <= 0 {
		fmt.Println("Game Over! You've run out of lives.")
	} else {
		fmt.Println("Congratulations! You've completed all levels!")
	}
	fmt.Printf("Final Score: %d\n", g.score)
}
