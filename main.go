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
		players := []board.FieldState{board.XSymbolInteger, board.OSymbolInteger}
		for _, player := range players {
			for {
				fmt.Println("Player " + board.Symbols[player] + " make turn now!")

				if playWithAI && board.Symbols[player] == board.Symbols[board.OSymbolInteger] {
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
				fmt.Println(board.Symbols[player] + " Win!")
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
