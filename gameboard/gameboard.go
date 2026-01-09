package board

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

const ColumnsNumber int = 3
const RowsNumber int = 3

type FieldState int // is used as enum for possible Board values

const EmptyFieldSymbolInteger FieldState = 0
const XSymbolInteger FieldState = 1 // first player symbol
const OSymbolInteger FieldState = 2 // second player or AI symbol

type Board struct {
	fields [ColumnsNumber][RowsNumber]FieldState
}

var Symbols = map[FieldState]string{
	EmptyFieldSymbolInteger: "*",
	XSymbolInteger:          "X",
	OSymbolInteger:          "O",
}

// ClearBoard makes Board absolutely clear (makes all fields' state equal to EmptyFieldSymbolInteger)
func (p *Board) ClearBoard() {
	for i := range p.fields {
		for j := range p.fields[i] {
			p.fields[i][j] = EmptyFieldSymbolInteger
		}
	}
}

// MakeMove adds symbol of to chosen field on Board
func (p *Board) MakeMove(symbol FieldState, row, col int) error {
	if p.fields[row][col] == EmptyFieldSymbolInteger {
		p.fields[row][col] = symbol
		return nil
	} else {
		return errors.New("this position on Board is already occupied")
	}
}

// ShowBoardState prints Board in console
func (p *Board) ShowBoardState() {
	fmt.Println("  0   1   2") // show columns' order index
	for rowIndex := range p.fields {
		line := ""
		for colIndex := range p.fields[rowIndex] {
			if colIndex < ColumnsNumber-1 {
				line += Symbols[p.fields[rowIndex][colIndex]] + " | "
			} else {
				line += Symbols[p.fields[rowIndex][colIndex]]
			}
		}
		fmt.Println(strconv.Itoa(rowIndex) + " " + line) // show row index and state of Board
	}

	fmt.Print("\n#############\n") // visually separates different states of Board
}

// CheckIfWinningCondition checks if any of the players won the game by analyzing Board state
func (p *Board) CheckIfWinningCondition() bool {

	// check all rows on Board if they are all occupied with the same symbols
	for i := range p.fields {
		if p.fields[i][0] == p.fields[i][1] && p.fields[i][0] == p.fields[i][2] {
			if p.fields[i][0] != EmptyFieldSymbolInteger {
				return true
			}
		}
	}

	// check all columns on Board if they are all occupied with the same symbols
	for i := range p.fields {
		if p.fields[0][i] == p.fields[1][i] && p.fields[1][i] == p.fields[2][i] {
			if p.fields[0][i] != EmptyFieldSymbolInteger {
				return true
			}
		}
	}

	// check first cross-line winning condition
	if p.fields[0][0] == p.fields[1][1] && p.fields[0][0] == p.fields[2][2] {
		if p.fields[0][0] != EmptyFieldSymbolInteger {
			return true
		}
	}

	// check second cross-line winning condition
	if p.fields[0][2] == p.fields[1][1] && p.fields[0][2] == p.fields[2][0] {
		if p.fields[0][2] != EmptyFieldSymbolInteger {
			return true
		}
	}

	return false
}

func (p *Board) occupyCrossLineIfCanLeadToVictory(aiSymbol FieldState) bool {
	winRow, winCol := -1, -1

	if p.fields[0][0] == p.fields[1][1] {
		winRow, winCol = 2, 2
	} else if p.fields[2][2] == p.fields[1][1] {
		winRow, winCol = 0, 0
	} else if p.fields[0][2] == p.fields[1][1] {
		winRow, winCol = 2, 0
	} else if p.fields[2][0] == p.fields[1][1] {
		winRow, winCol = 0, 2
	}

	if winRow != -1 {
		err := p.MakeMove(aiSymbol, winRow, winCol)
		if err == nil {
			return true
		}
	}

	return false
}

func (p *Board) occupyRowIfCanLeadToVictory(aiSymbol FieldState) bool {
	// check if any of rows or cols could give victory to any of the sides
	// means at least one identical symbols in one row
	for rowIndex := 0; rowIndex < RowsNumber; rowIndex++ {
		// contains state of one row
		enemySymbolsInRow := 0
		aiSymbolsInRow := 0
		freeColIndex := -1 // indicates that all positions are occupied, otherwise we get index of the position

		for colIndex := 0; colIndex < ColumnsNumber; colIndex++ {
			if p.fields[rowIndex][colIndex] == aiSymbol {
				aiSymbolsInRow++
			} else if p.fields[rowIndex][colIndex] == EmptyFieldSymbolInteger {
				freeColIndex = colIndex
			} else {
				enemySymbolsInRow++
			}
		}

		// check if AI can win first
		if aiSymbolsInRow == RowsNumber-1 && freeColIndex != -1 {
			err := p.MakeMove(aiSymbol, rowIndex, freeColIndex)
			if err == nil {
				return true
			}
		}

		// prevent player from winning
		if enemySymbolsInRow == RowsNumber-1 && freeColIndex != -1 {
			err := p.MakeMove(aiSymbol, rowIndex, freeColIndex)
			if err == nil {
				return true
			}
		}
	}

	return false
}

func (p *Board) occupyColumnIfCanLeadToVictory(aiSymbol FieldState) bool {
	// check if any of rows or cols could give victory to any of the sides
	// means at least one identical symbols in one column
	for colIndex := 0; colIndex < ColumnsNumber; colIndex++ {
		// contains state of one column
		enemySymbolsInCol := 0
		aiSymbolsInCol := 0
		freeRowIndex := -1 // indicates that all positions are occupied, otherwise we get index of the position

		for rowIndex := 0; rowIndex < RowsNumber; rowIndex++ {
			if p.fields[rowIndex][colIndex] == aiSymbol {
				aiSymbolsInCol++
			} else if p.fields[rowIndex][colIndex] == EmptyFieldSymbolInteger {
				freeRowIndex = rowIndex
			} else {
				enemySymbolsInCol++
			}
		}

		// check if AI can win
		if aiSymbolsInCol == ColumnsNumber-1 && freeRowIndex != -1 {
			err := p.MakeMove(aiSymbol, freeRowIndex, colIndex)
			if err == nil {
				return true
			}
		}

		// prevent player from winning
		if enemySymbolsInCol == ColumnsNumber-1 && freeRowIndex != -1 {
			err := p.MakeMove(aiSymbol, freeRowIndex, colIndex)
			if err == nil {
				return true
			}
		}
	}

	return false
}

func (p *Board) occupyCornersIfCanLeadToVictory(aiSymbol FieldState) bool {
	corners := [][2]int{
		{0, 0}, {0, ColumnsNumber - 1}, {RowsNumber - 1, 0}, {RowsNumber - 1, ColumnsNumber - 1},
	}

	// shuffle possible AI choices to make game look more different
	for i := range corners {
		newIndex := rand.Intn(len(corners))
		corners[i], corners[newIndex] = corners[newIndex], corners[i]
	}

	for i := range corners {
		err := p.MakeMove(aiSymbol, corners[i][0], corners[i][1])
		if err == nil {
			return true
		}
	}

	return false
}

func (p *Board) occupyAnyField(aiSymbol FieldState) {
	attempts := 100 // prevents infinite loop, though everything is checked in outer code
	for attempts > 0 {
		row := rand.Intn(RowsNumber)
		col := rand.Intn(ColumnsNumber)
		err := p.MakeMove(aiSymbol, row, col)
		if err == nil {
			return
		}
		attempts -= 1
	}
}

func (p *Board) occupyCenterOfTheBoard(aiSymbol FieldState) bool {
	if p.fields[1][1] == EmptyFieldSymbolInteger {
		err := p.MakeMove(aiSymbol, 1, 1)
		if err == nil {
			return true
		}
	}

	return false
}

// AIMakeMove contains logic for AI player to make moves
// AI symbol is always "O" symbol OSymbolInteger
func (p *Board) AIMakeMove(aiSymbol FieldState) {
	// always try to occupy center of the field first
	if p.occupyCenterOfTheBoard(aiSymbol) {
		return
	}

	// check if any row or line can lead of players to victory
	if p.occupyRowIfCanLeadToVictory(aiSymbol) {
		return
	}

	if p.occupyColumnIfCanLeadToVictory(aiSymbol) {
		return
	}

	// check if any cross-line can lead to victory
	if p.occupyCrossLineIfCanLeadToVictory(aiSymbol) {
		return
	}

	// no direct win move are found, now AI tries to take the best possible fields on Board

	// try to occupy corners of the field if they are free
	if p.occupyCornersIfCanLeadToVictory(aiSymbol) {
		return
	}

	// just fill all left spots of the field (actually the code is never reached, but it prevents possible problems
	p.occupyAnyField(aiSymbol) // always can work fine since we check in outer code that some empty fields exist
}
