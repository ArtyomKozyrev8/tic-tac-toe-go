package ui

import (
	"errors"
	"fmt"
)

type BoardPartName string

// RowBoardPartName and ColumnBoardPartName are possible values for getUserInput
const RowBoardPartName BoardPartName = "Row"
const ColumnBoardPartName BoardPartName = "Column"

func getUserInputBasis(inputName BoardPartName) (int, error) {
	var index int
	fmt.Print("Enter " + inputName + ": ")
	_, errRow := fmt.Scan(&index)
	if errRow != nil {
		return 0, errRow
	}

	if index < 0 || index > 2 {
		return 0, errors.New(string(inputName) + "should be 0, 1, 2")
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
	playWithAI := true
	fmt.Print("Do you want to play with AI? (y/n): ")
	_, err := fmt.Scanln(&playWithAIInput)
	if err != nil {
		fmt.Println(err)
		playWithAI = true
	}

	if playWithAIInput == "y" || playWithAIInput == "Y" {
		playWithAI = true
	} else {
		playWithAI = false
	}

	return playWithAI
}
