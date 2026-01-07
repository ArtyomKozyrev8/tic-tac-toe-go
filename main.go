package main

import (
	"errors"
	"fmt"
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
