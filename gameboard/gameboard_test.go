package board

import (
	"errors"
	"fmt"
	"testing"
)

func TestNotEmptyFieldError_Error(t *testing.T) {
	err := NotEmptyFieldError{row: 1, col: 2}
	expected := "Field row=1 col=2 is not empty"

	if err.Error() != expected {
		t.Errorf("got %q, want %q", err.Error(), expected)
	}
}

// Simple test for OutOfBoundsError
func TestOutOfBoundsError_Error(t *testing.T) {
	err := OutOfBoundsError{row: 5, col: 5}
	// Use RowsNumber-1 to match your code's dynamic logic
	expected := fmt.Sprintf("row=5 and col=5 should be from 0 to %d ", RowsNumber-1)

	if err.Error() != expected {
		t.Errorf("got %q, want %q", err.Error(), expected)
	}
}

func TestBoard_GetSymbol(t *testing.T) {
	tests := []struct {
		caseName   string
		fieldState FieldState
		want       string
	}{
		{"Test-1", Empty, "*"},
		{"Test-2", X, "X"},
		{"Test-3", O, "O"},
	}

	b := Board{}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			got := b.GetSymbol(tt.fieldState)
			if got != tt.want {
				t.Errorf("%s: expected %s, got %s", tt.caseName, tt.want, got)
			}
		})
	}
}

func TestBoard_ClearBoard(t *testing.T) {
	b := Board{}
	b.fields[0][0] = 1
	b.fields[1][1] = 1
	b.occupiedFieldsNumber = 2

	b.ClearBoard()

	if b.fields[0][0] != 0 || b.fields[1][1] != 0 {
		t.Errorf("ClearBoard: expected 0 for all fields")
	}

	if b.occupiedFieldsNumber != 0 {
		t.Errorf("ClearBoard: expected 0 for b.occupiedFieldsNumber")
	}
}

func TestBoard_MakeMove(t *testing.T) {
	b := Board{}
	b.fields[0][0] = 1

	err := b.MakeMove(X, 0, 0)
	if !errors.As(err, &NotEmptyFieldError{}) {
		t.Errorf("MakeMove: expected NotEmptyFieldError")
	}

	err = b.MakeMove(X, 9, 0)
	if !errors.As(err, &OutOfBoundsError{}) {
		t.Errorf("MakeMove: expected OutOfBoundsError")
	}

	err = b.MakeMove(X, 0, 19)
	if !errors.As(err, &OutOfBoundsError{}) {
		t.Errorf("MakeMove: expected OutOfBoundsError")
	}

	err = b.MakeMove(X, 1, 1)
	if err != nil {
		t.Errorf("MakeMove: unexpected error: %s", err)
	}

	if b.fields[1][1] != 1 {
		t.Errorf("MakeMove: expected %v, got %v", 1, b.fields[1][1])
	}

	if b.occupiedFieldsNumber != 1 {
		t.Errorf("MakeMove: expected %v, got %v", 1, b.occupiedFieldsNumber)
	}
}

func TestBoard_String(t *testing.T) {
	want := "  0   1   2\n"

	b := Board{}
	got := b.String()[0:12]
	if got != want {
		t.Errorf("String: expected %s, got %s", want, got)
	}
}

func TestBoard_CheckIfWinningCondition(t *testing.T) {
	tests := []struct {
		caseName string
		coords   [][2]int
		want     bool
	}{
		{"Test-1", [][2]int{{0, 0}, {0, 1}, {0, 2}}, true},
		{"Test-2", [][2]int{{1, 0}, {1, 1}, {1, 2}}, true},
		{"Test-3", [][2]int{{2, 0}, {2, 1}, {2, 2}}, true},
		{"Test-4", [][2]int{{0, 0}, {1, 0}, {2, 0}}, true},
		{"Test-5", [][2]int{{0, 1}, {1, 1}, {2, 1}}, true},
		{"Test-6", [][2]int{{0, 2}, {1, 2}, {2, 2}}, true},
		{"Test-7", [][2]int{{0, 0}, {1, 1}, {2, 2}}, true},
		{"Test-8", [][2]int{{0, 2}, {1, 1}, {2, 0}}, true},
		{"Test-9", [][2]int{{0, 2}, {1, 1}, {2, 1}}, false},
		{"Test-10", [][2]int{{0, 2}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			b := Board{}
			for _, v := range tt.coords {
				err := b.MakeMove(O, v[0], v[1])
				if err != nil {
					t.Errorf("%s MakeMove: unexpected error: %s", tt.caseName, err)
				}
			}

			got := b.CheckIfWinningCondition()
			if got != tt.want {
				t.Errorf("%s CheckIfWinningCondition: expected %v, got %v", tt.caseName, tt.want, got)
			}
		})
	}
}
