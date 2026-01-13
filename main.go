package main

import (
	"fmt"

	"github.com/ArtyomKozyrev8/tic-tac-toe-go/gameboard"
	"github.com/ArtyomKozyrev8/tic-tac-toe-go/ui"
)

func main() {
	theBoard := board.Board{}
	theBoard.ShowBoardState()
	totalSteps := 0 // checks if all Board is full
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

					if theBoard.IsEmpty(row, col) {
						theBoard.MakeMove(player, row, col)
						theBoard.ShowBoardState()
						break
					} else {
						fmt.Println("The place is already taken!")
					}
				}
			}
			if theBoard.CheckIfWinningCondition() {
				fmt.Println(theBoard.GetSymbol(player) + " Win!")
				break game
			} else {
				totalSteps += 1
				if totalSteps == board.ColumnsNumber*board.RowsNumber {
					fmt.Println("Nobody win... Play another round")
					totalSteps = 0
					theBoard.ClearBoard()
				}
			}
		}
	}
}
