package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

var PlayFieldSymbols = map[int]string{
	0: "*",
	1: "X",
	2: "O",
}

type PlayField struct {
	fields [3][3]int
}

func (p *PlayField) Clear() {
	for i := range p.fields {
		for j := range p.fields[i] {
			p.fields[i][j] = 0
		}
	}
}

func (p *PlayField) MakeStep(symbol, row, col int) error {
	if p.fields[row][col] == 0 {
		p.fields[row][col] = symbol
		return nil
	} else {
		return errors.New("The position is already occupied!")
	}
}

func (p *PlayField) CheckIfWin() bool {
	for i := range p.fields {
		if p.fields[i][0] == p.fields[i][1] && p.fields[i][0] == p.fields[i][2] {
			if p.fields[i][0] != 0 {
				return true
			}
		}
	}

	for i := range p.fields {
		if p.fields[0][i] == p.fields[1][i] && p.fields[1][i] == p.fields[2][i] {
			if p.fields[0][i] != 0 {
				return true
			}
		}
	}

	if p.fields[0][0] == p.fields[1][1] && p.fields[0][0] == p.fields[2][2] {
		if p.fields[0][0] != 0 {
			return true
		}
	}

	if p.fields[0][2] == p.fields[1][1] && p.fields[0][2] == p.fields[2][0] {
		if p.fields[0][2] != 0 {
			return true
		}
	}

	return false
}

func (p *PlayField) AiMakeMove(aiSymbol int) {
	// always try to occupy center of the field
	if p.fields[1][1] == 0 {
		err := p.MakeStep(aiSymbol, 1, 1)
		if err == nil {
			return
		}
	}

	// check if any of rows or cols could give victory to any of the sides
	for rowIndex := 0; rowIndex < 3; rowIndex++ {
		enemySymbolsInRow := 0
		aiSymbolsInRow := 0
		freeColIndex := -1

		enemySymbolsInCol := 0
		aiSymbolsInCol := 0
		freeRowIndex := -1

		for colIndex := 0; colIndex < 3; colIndex++ {
			if p.fields[rowIndex][colIndex] == aiSymbol {
				aiSymbolsInRow++
			} else if p.fields[rowIndex][colIndex] == 0 {
				freeColIndex = colIndex
			} else {
				enemySymbolsInRow++
			}

			localRowIndex := colIndex
			localColIndex := rowIndex
			if p.fields[localRowIndex][localColIndex] == aiSymbol {
				aiSymbolsInCol++
			} else if p.fields[localRowIndex][localColIndex] == 0 {
				freeRowIndex = localRowIndex
			} else {
				enemySymbolsInCol++
			}
		}

		if aiSymbolsInRow == 2 && freeColIndex != -1 {
			err := p.MakeStep(aiSymbol, rowIndex, freeColIndex)
			if err == nil {
				return
			}
		}

		if aiSymbolsInCol == 2 && freeRowIndex != -1 {
			err := p.MakeStep(aiSymbol, freeRowIndex, rowIndex)
			if err == nil {
				return
			}
		}

		if enemySymbolsInRow == 2 && freeColIndex != -1 {
			err := p.MakeStep(aiSymbol, rowIndex, freeColIndex)
			if err == nil {
				return
			}
		}

		if enemySymbolsInCol == 2 && freeRowIndex != -1 {
			err := p.MakeStep(aiSymbol, freeRowIndex, rowIndex)
			if err == nil {
				return
			}
		}
	}

	// try to occupy corners of the field if they are free
	corners := [][2]int{
		{0, 0}, {0, 2}, {2, 0}, {2, 2},
	}

	for i := range corners {
		newIndex := rand.Intn(len(corners))
		corners[i], corners[newIndex] = corners[newIndex], corners[i]
	}

	for i := range corners {
		err := p.MakeStep(aiSymbol, corners[i][0], corners[i][1])
		if err == nil {
			return
		}
	}

	// just fill all left spots of the field
	for {
		row := rand.Intn(3)
		col := rand.Intn(3)
		err := p.MakeStep(aiSymbol, row, col)
		if err == nil {
			return
		}
	}

}

func (p *PlayField) Print() {
	fmt.Println("  0   1   2")
	for i := range p.fields {
		line := ""
		for j := range p.fields[i] {
			if j < 2 {
				line += PlayFieldSymbols[p.fields[i][j]] + " | "
			} else {
				line += PlayFieldSymbols[p.fields[i][j]]
			}
		}
		fmt.Println(strconv.Itoa(i) + " " + line)
	}

	fmt.Println("\n****************\n")
}

func getUserInput(inputName string) (int, error) {
	var index int
	fmt.Print("Enter " + inputName + ": ")
	_, errRow := fmt.Scan(&index)
	if errRow != nil {
		return 0, errRow
	}

	if index < 0 || index > 2 {
		return 0, errors.New(inputName + "should be 0, 1, 2")
	}

	return index, nil
}

func getUserInputCircle(inputName string) int {
	for {
		index, err := getUserInput(inputName)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			return index
		}
	}
}

func main() {
	p := PlayField{}
	p.Print()
	totalSteps := 0
game:
	for {
		players := []int{1, 2}
		for _, player := range players {
			for {
				fmt.Println("Player " + PlayFieldSymbols[player] + " go!")

				if PlayFieldSymbols[player] == "O" {
					p.AiMakeMove(player)
					p.Print()
					break
				} else {
					row := getUserInputCircle("Row")
					col := getUserInputCircle("Column")

					err := p.MakeStep(player, row, col)
					if err != nil {
						fmt.Println(err)
					} else {
						p.Print()
						break
					}
				}
			}
			if p.CheckIfWin() {
				fmt.Println(PlayFieldSymbols[player] + " Win!")
				break game
			} else {
				totalSteps += 1
				fmt.Println(totalSteps)
				if totalSteps == 9 {
					fmt.Println("Nobody win... Play another round")
					totalSteps = 0
					p.Clear()
				}
			}
		}
	}
}
