package board

// TODO add tests to packages when all is improved!

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const ColumnsNumber int = 3
const RowsNumber int = 3

type FieldState int // is used as enum for possible Board values

const (
	Empty FieldState = iota
	X
	O
)

type nextMove struct {
	rowIndex             int
	colIndex             int
	leadToVictory        bool
	preventsEnemyVictory bool
}

type Board struct {
	fields [ColumnsNumber][RowsNumber]FieldState
}

func (b *Board) GetSymbol(state FieldState) string {
	switch state {
	case Empty:
		return "*"
	case X:
		return "X"
	case O:
		return "O"
	default:
		return "?"
	}
}

// ClearBoard makes Board absolutely clear (makes all fields' state equal to Empty)
func (b *Board) ClearBoard() {
	for i := range b.fields {
		for j := range b.fields[i] {
			b.fields[i][j] = Empty
		}
	}
}

// MakeMove adds symbol of to chosen field on Board
func (b *Board) MakeMove(symbol FieldState, row, col int) {
	b.fields[row][col] = symbol
}

func (b *Board) IsEmpty(row, col int) bool {
	return b.fields[row][col] == Empty
}

// ShowBoardState prints Board in console
func (b *Board) ShowBoardState() {
	strBuilder := strings.Builder{}
	strBuilder.Grow(100)
	strBuilder.WriteString("  0   1   2\n") // show columns' order index

	for rowIndex, row := range b.fields {
		strBuilder.WriteString(strconv.Itoa(rowIndex))
		strBuilder.WriteByte(' ')
		for colIndex, item := range row {
			strBuilder.WriteString(b.GetSymbol(item))
			if colIndex < ColumnsNumber-1 {
				strBuilder.WriteString(" | ")
			}
		}
		strBuilder.WriteByte('\n')
	}

	strBuilder.WriteString("\n#############\n")
	// strBuilder.String()

	fmt.Print(strBuilder.String())
}

// CheckIfWinningCondition checks if any of the players won the game by analyzing Board state
func (b *Board) CheckIfWinningCondition() bool {

	// check all rows on Board if they are all occupied with the same symbols
	for i := range b.fields {
		if b.fields[i][0] == b.fields[i][1] && b.fields[i][0] == b.fields[i][2] {
			if b.fields[i][0] != Empty {
				return true
			}
		}
	}

	// check all columns on Board if they are all occupied with the same symbols
	for i := range b.fields {
		if b.fields[0][i] == b.fields[1][i] && b.fields[1][i] == b.fields[2][i] {
			if b.fields[0][i] != Empty {
				return true
			}
		}
	}

	// check first cross-line winning condition
	if b.fields[0][0] == b.fields[1][1] && b.fields[0][0] == b.fields[2][2] {
		if b.fields[0][0] != Empty {
			return true
		}
	}

	// check second cross-line winning condition
	if b.fields[0][2] == b.fields[1][1] && b.fields[0][2] == b.fields[2][0] {
		if b.fields[0][2] != Empty {
			return true
		}
	}

	return false
}

func (b *Board) occupyCrossLineIfCanLeadToVictory() *nextMove {
	winRow, winCol := -1, -1

	// we compare with center, because AI always tries to occupy center at first
	if b.fields[0][0] == b.fields[1][1] {
		winRow, winCol = 2, 2
	} else if b.fields[2][2] == b.fields[1][1] {
		winRow, winCol = 0, 0
	} else if b.fields[0][2] == b.fields[1][1] {
		winRow, winCol = 2, 0
	} else if b.fields[2][0] == b.fields[1][1] {
		winRow, winCol = 0, 2
	}

	if winRow != -1 {
		if b.fields[winRow][winCol] == Empty {
			symbol := b.fields[1][1] // it is always occupied in the first AI move
			if symbol == X {
				return &nextMove{winRow, winCol, false, true}
			} else {
				return &nextMove{winRow, winCol, true, false}
			}
		}
	}

	return nil
}

func (b *Board) occupyRowIfCanLeadToVictory(aiSymbol FieldState) *nextMove {
	var preventEnemyWinChoice *nextMove = nil

	// check if any of rows or cols could give victory to any of the sides
	// means at least one identical symbols in one row
	for rowIndex := 0; rowIndex < RowsNumber; rowIndex++ {
		// contains state of one row
		enemySymbolsInRow := 0
		aiSymbolsInRow := 0
		freeColIndex := -1 // indicates that all positions are occupied, otherwise we get index of the position

		for colIndex := 0; colIndex < ColumnsNumber; colIndex++ {
			if b.fields[rowIndex][colIndex] == aiSymbol {
				aiSymbolsInRow++
			} else if b.fields[rowIndex][colIndex] == Empty {
				freeColIndex = colIndex
			} else {
				enemySymbolsInRow++
			}
		}

		// check if AI can win first
		if aiSymbolsInRow == RowsNumber-1 && freeColIndex != -1 {
			// return if found victory choice
			return &nextMove{rowIndex, freeColIndex, true, false}
		}

		// prevent player from winning
		if enemySymbolsInRow == RowsNumber-1 && freeColIndex != -1 {
			preventEnemyWinChoice = &nextMove{rowIndex, freeColIndex, false, true}
		}
	}

	return preventEnemyWinChoice
}

func (b *Board) occupyColumnIfCanLeadToVictory(aiSymbol FieldState) *nextMove {
	var preventEnemyWinChoice *nextMove = nil

	// check if any of rows or cols could give victory to any of the sides
	// means at least one identical symbols in one column
	for colIndex := 0; colIndex < ColumnsNumber; colIndex++ {
		// contains state of one column
		enemySymbolsInCol := 0
		aiSymbolsInCol := 0
		freeRowIndex := -1 // indicates that all positions are occupied, otherwise we get index of the position

		for rowIndex := 0; rowIndex < RowsNumber; rowIndex++ {
			if b.fields[rowIndex][colIndex] == aiSymbol {
				aiSymbolsInCol++
			} else if b.fields[rowIndex][colIndex] == Empty {
				freeRowIndex = rowIndex
			} else {
				enemySymbolsInCol++
			}
		}

		// check if AI can win
		if aiSymbolsInCol == ColumnsNumber-1 && freeRowIndex != -1 {
			// return if found victory choice
			return &nextMove{freeRowIndex, colIndex, true, false}
		}

		// prevent player from winning
		if enemySymbolsInCol == ColumnsNumber-1 && freeRowIndex != -1 {
			preventEnemyWinChoice = &nextMove{freeRowIndex, colIndex, false, true}
		}
	}

	return preventEnemyWinChoice
}

func (b *Board) occupyCorners() *nextMove {
	corners := [][2]int{
		{0, 0}, {0, ColumnsNumber - 1}, {RowsNumber - 1, 0}, {RowsNumber - 1, ColumnsNumber - 1},
	}

	// shuffle possible AI choices to make game look more different
	for i := range corners {
		newIndex := rand.Intn(len(corners))
		corners[i], corners[newIndex] = corners[newIndex], corners[i]
	}

	for i := range corners {
		if b.fields[corners[i][0]][corners[i][1]] == Empty {
			return &nextMove{corners[i][0], corners[i][1], false, false}
		}
	}

	return nil
}

func (b *Board) occupyAnyField() *nextMove {
	for i := range b.fields {
		for j := range b.fields[i] {
			if b.fields[i][j] == Empty {
				return &nextMove{i, j, false, false}
			}
		}
	}
	return nil
}

func (b *Board) occupyCenterOfBoard() *nextMove {
	if b.fields[1][1] == Empty {
		return &nextMove{1, 1, false, false}
	}

	return nil
}

func (b *Board) chooseBestPossibleMove() *nextMove {
	// always try to occupy center of the field first
	moveCenter := b.occupyCenterOfBoard()
	if moveCenter != nil {
		return moveCenter
	}

	var choices [3]*nextMove

	choices[0] = b.occupyCrossLineIfCanLeadToVictory()
	choices[1] = b.occupyRowIfCanLeadToVictory(O)
	choices[2] = b.occupyColumnIfCanLeadToVictory(O)

	// check if any move can lead to AI Victory
	for _, choice := range choices {
		if choice != nil {
			if choice.leadToVictory {
				return choice
			}
		}
	}

	// check if any move can prevent User Victory
	for _, choice := range choices {
		if choice != nil {
			if choice.preventsEnemyVictory {
				return choice
			}
		}
	}

	moveCorner := b.occupyCorners()
	if moveCorner != nil {
		return moveCorner
	}

	return b.occupyAnyField()
}

// AIMakeMove contains logic for AI player to make moves
// AI symbol is always "O" symbol O
func (b *Board) AIMakeMove(aiSymbol FieldState) {
	move := b.chooseBestPossibleMove()
	if move != nil {
		b.MakeMove(aiSymbol, move.rowIndex, move.colIndex)
	}
}
