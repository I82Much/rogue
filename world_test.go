package rogue

import (
	"testing"
)

func TestNewWorld(t *testing.T) {
	w := NewWorld(5, 6)
	if numRows := w.Rows(); numRows != 5 {
		t.Errorf("expected 5 rows got %d", numRows)
	}
	if numCols := w.Cols(); numCols != 6 {
		t.Errorf("expected 6 cols got %d", numCols)
	}

	if numRows := len(w.tiles); numRows != 5 {
		t.Errorf("expected 5 rows in data structure had %d", numRows)
	}
	if numCols := len(w.tiles[0]); numCols != 6 {
		t.Errorf("expected 6 cols in data structure had %d", numCols)
	}
}
