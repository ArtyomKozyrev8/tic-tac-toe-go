package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ArtyomKozyrev8/tic-tac-toe-go/gameboard"
	"github.com/ArtyomKozyrev8/tic-tac-toe-go/ui"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // keep old style

	theBoard := board.Board{}
	theBoard.ShowBoardState()
	playWithAI := ui.DecideToPlayWithAI()

game:
	for {
		players := []board.FieldState{board.X, board.O}
		for _, player := range players {
			for {
				fmt.Println("Player " + theBoard.GetSymbol(player) + " make turn now!")

				if playWithAI && theBoard.GetSymbol(player) == theBoard.GetSymbol(board.O) {
					theBoard.AIMakeMove(player)
					theBoard.ShowBoardState()
					break
				} else {
					row := ui.GetUserInput(ui.RowBoardPartName)
					col := ui.GetUserInput(ui.ColumnBoardPartName)

					err := theBoard.MakeMove(player, row, col)

					if err != nil {
						fmt.Println(err)
					} else {
						theBoard.ShowBoardState()
						break
					}
				}
			}
			if theBoard.CheckIfWinningCondition() {
				fmt.Println(theBoard.GetSymbol(player) + " Win!")
				break game
			} else {
				// checks if all Board is full
				if theBoard.IsDraw() {
					fmt.Println("Nobody win... Play another round")
					theBoard.ClearBoard()
				}
			}
		}
	}
}
