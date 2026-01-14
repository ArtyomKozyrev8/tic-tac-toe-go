package board

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const ColumnsNumber int = 3
const RowsNumber int = 3

type FieldState int // is used as enum for possible Board values

type NotEmptyFieldError struct {
	row int
	col int
}

func (err NotEmptyFieldError) Error() string {
	return fmt.Sprintf("Field row=%d col=%d is not empty", err.row, err.col)
}

type OutOfBoundsError struct {
	row int
	col int
}

func (err OutOfBoundsError) Error() string {
	return fmt.Sprintf("row=%d and col=%d should be from 0 to %d ", err.row, err.col, RowsNumber-1)
}

const (
	Empty FieldState = iota
	X
	O // AI PLAYER ALWAYS (If Player decided to play with AI, otherwise it is second human player
)

type nextMove struct {
	rowIndex             int
	colIndex             int
	leadToVictory        bool
	preventsEnemyVictory bool
}

type Board struct {
	fields               [ColumnsNumber][RowsNumber]FieldState
	occupiedFieldsNumber int
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
	b.occupiedFieldsNumber = 0
}

// MakeMove adds symbol of to chosen field on Board
func (b *Board) MakeMove(symbol FieldState, row, col int) error {
	if row < 0 || col < 0 || row >= RowsNumber || col >= ColumnsNumber {
		err := OutOfBoundsError{row, col}
		return err
	}

	if b.fields[row][col] != Empty {
		err := NotEmptyFieldError{row, col}
		return err
	} else {
		b.fields[row][col] = symbol
		b.occupiedFieldsNumber += 1
	}
	return nil
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

func (b *Board) occupyCrossLineIfCanLeadToVictory(aiSymbol FieldState) *nextMove {
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
			if symbol == aiSymbol {
				return &nextMove{winRow, winCol, true, false}
			} else {
				return &nextMove{winRow, winCol, false, true}
			}
		}
	}

	return nil
}

// detectForkConditions - tries to catch at least some Win-Forks
func (b *Board) detectForkConditions() *nextMove {
	if b.fields[1][1] != O {
		return nil
	}

	if b.fields[0][0] == b.fields[2][2] && b.fields[2][2] == X {
		return &nextMove{1, 0, false, true}
	}

	if b.fields[0][2] == b.fields[2][0] && b.fields[2][0] == X {
		return &nextMove{1, 0, false, true}
	}

	return nil
}

// occupyLine checks if occupation of any Row or any Column can lead to AI victory or Enemy Player Victory
// aiSymbol - is a symbol which is used by AI player, in the program AI always use O
// ColCheck if false we check all Rows, if true we check all Columns
func (b *Board) occupyLine(aiSymbol FieldState, ColCheck bool) *nextMove {
	var preventEnemyWinChoice *nextMove = nil

	lineLen := RowsNumber // RowNumber = ColumnsNumber - we can take any

	// check if any of rows or cols could give victory to any of the sides
	// means at least one identical symbols in one row
	for i := 0; i < lineLen; i++ {
		enemySymbols, aiSymbols, freeIndex := 0, 0, -1

		for j := 0; j < lineLen; j++ {
			item := b.fields[i][j]
			if ColCheck {
				item = b.fields[j][i]
			}

			if item == aiSymbol {
				aiSymbols++
			} else if item == Empty {
				freeIndex = j
			} else {
				enemySymbols++
			}
		}

		// check if AI can win first
		if aiSymbols == lineLen-1 && freeIndex != -1 {
			// return if found victory choice
			if ColCheck {
				return &nextMove{freeIndex, i, true, false}
			}

			return &nextMove{i, freeIndex, true, false}
		}

		// prevent player from winning
		if enemySymbols == lineLen-1 && freeIndex != -1 {
			if ColCheck {
				preventEnemyWinChoice = &nextMove{freeIndex, i, false, true}
			} else {
				preventEnemyWinChoice = &nextMove{i, freeIndex, false, true}
			}
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

	if b.occupiedFieldsNumber == 3 {
		preventFork := b.detectForkConditions()
		if preventFork != nil {
			return preventFork
		}
	}

	var choices [3]*nextMove

	choices[0] = b.occupyCrossLineIfCanLeadToVictory(O)
	choices[1] = b.occupyLine(O, false) // check all rows
	choices[2] = b.occupyLine(O, true)  // check all columns

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

func (b *Board) IsDraw() bool {
	return b.occupiedFieldsNumber == ColumnsNumber*RowsNumber
}

// AIMakeMove contains logic for AI player to make moves
// AI symbol is always "O" symbol O
func (b *Board) AIMakeMove(aiSymbol FieldState) {
	move := b.chooseBestPossibleMove()
	if move != nil {
		err := b.MakeMove(aiSymbol, move.rowIndex, move.colIndex)
		if err != nil {
			fmt.Printf("Error making move: %v", err)
		}
	}
}
