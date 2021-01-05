/*
 @Desc

 @Date 2020-04-02 13:10
 @Author yinjk
*/
package main

import (
	"fmt"
	"time"
)

var (
	greenBg   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	whiteBg   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellowBg  = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	redBg     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blueBg    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magentaBg = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyanBg    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	green     = string([]byte{27, 91, 51, 50, 109})
	white     = string([]byte{27, 91, 51, 55, 109})
	yellow    = string([]byte{27, 91, 51, 51, 109})
	red       = string([]byte{27, 91, 51, 49, 109})
	blue      = string([]byte{27, 91, 51, 52, 109})
	magenta   = string([]byte{27, 91, 51, 53, 109})
	cyan      = string([]byte{27, 91, 51, 54, 109})
	reset     = string([]byte{27, 91, 48, 109})
)

func main() {
	board := [][]int{
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0},
		{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0},
		{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
		{0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	ticker := time.NewTicker(time.Millisecond * (1000 / 24))
	printFrame(len(board), len(board[0]), "   ")
	for {
		select {
		case <-ticker.C:
			gameOfLife(board)
			prettyPrint(board)
		}
	}
}

func prettyPrint(board [][]int) {
	left(1000)
	up(1000)
	down(1)
	for _, v := range board {
		right(1)
		for i, a := range v {
			if i%2 == 0 {
				fmt.Print(cyanBg)
			} else {
				fmt.Print(whiteBg)
			}
			if a == 1 {
				fmt.Print(red, " o ")
			} else {
				fmt.Print(white, "   ")
			}
		}
		fmt.Println()
	}
}
func printFrame(r, col int, fill string) {
	up(1000)
	left(1000)
	fillLen := len(fill)
	for i := 0; i < r+2; i++ {
		for j := 0; j < col+2; j++ {
			fmt.Print(green)
			if i == 0 && j == 0 {
				fmt.Print("┏")
			} else if i == 0 && j == col+1 {
				fmt.Print("┓")
			} else if i == r+1 && j == 0 {
				fmt.Print("┗")
			} else if i == r+1 && j == col+1 {
				fmt.Print("┛")
			} else if i == 0 || i == r+1 {
				for i := 0; i < fillLen; i++ {
					fmt.Print("━")
				}
			} else if j == 0 || j == col+1 {
				fmt.Print("┃")
			} else {
				if j%2 == 0 {
					fmt.Print(cyanBg)
				} else {
					fmt.Print(whiteBg)
				}
				fmt.Print(fill, reset)
			}
		}
		fmt.Println()
	}
	fmt.Print(reset)
}

func gameOfLife(board [][]int) {
	fmt.Print(cyanBg)
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			live := 0
			for i1 := -1; i1 <= 1; i1++ {
				for j1 := -1; j1 <= 1; j1++ {
					if i1 == 0 && j1 == 0 {
						continue
					}
					if i+i1 >= 0 && i+i1 < len(board) && j+j1 >= 0 && j+j1 < len(board[i]) && (board[i+i1][j+j1]&1 == 1) {
						live++
					}
				}
			}
			//从活到死
			if (live < 2 || live > 3) && board[i][j] == 1 {
				board[i][j] = -1
				//从死到活
			} else if live == 3 && board[i][j] == 0 {
				board[i][j] = 2
			}
		}
	}
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] > 0 {
				board[i][j] = 1
			} else {
				board[i][j] = 0
			}

		}
	}
}

func up(d int) {
	fmt.Print(fmt.Sprintf("\u001b[%dA", d))
}
func down(d int) {
	fmt.Print(fmt.Sprintf("\u001b[%dB", d))
}
func left(d int) {
	fmt.Print(fmt.Sprintf("\u001b[%dD", d))
}
func right(d int) {
	fmt.Print(fmt.Sprintf("\u001b[%dC", d))
}
