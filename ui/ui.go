package ui

import (
	"fmt"
	"strings"

	"github.com/ArtyomKozyrev8/tic-tac-toe-go/gameboard"
)

type BoardPartName string

// RowBoardPartName and ColumnBoardPartName are possible values for getUserInput
const (
	RowBoardPartName    BoardPartName = "Row"
	ColumnBoardPartName BoardPartName = "Column"
)

func getUserInputBasis(inputName BoardPartName) (int, error) {
	var index int
	fmt.Print("Enter " + inputName + ": ")
	_, errRow := fmt.Scanln(&index)
	if errRow != nil {
		return 0, errRow
	}

	if index < 0 || index > board.ColumnsNumber {
		return 0, fmt.Errorf("%s should be > 0 and < %d", inputName, board.ColumnsNumber)
	}

	return index, nil
}

func GetUserInput(inputName BoardPartName) int {
	for {
		index, err := getUserInputBasis(inputName)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			return index
		}
	}
}

func DecideToPlayWithAI() bool {
	var playWithAIInput string
	fmt.Print("Do you want to play with AI? (y/n): ")

	_, err := fmt.Scanln(&playWithAIInput)
	if err != nil {
		fmt.Println(err)
		return true
	}

	return strings.ToLower(playWithAIInput) == "y" || strings.ToLower(playWithAIInput) == "yes"
}
