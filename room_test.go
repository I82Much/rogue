package rogue

import (
	"testing"
)

func TestNewRoom(t *testing.T) {
	w := NewRoom(5, 6)
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

func TestMove(t *testing.T) {
	w := NewRoom(10, 10)
	w.Spawn(5, 5)
	if got := w.MovePlayer(-1, 0); got != Move {
		t.Errorf("couldn't move to unobstructed spot: %v", got)
	}
}

func TestMoveRemovesPlayer(t *testing.T) {
	w := NewRoom(10, 10)
	w.Spawn(5, 5)
	if w.CreatureAt(Loc(5, 5)) != PlayerCreature {
		t.Errorf("player didn't spawn")
	}
	w.MovePlayer(-1, 0)
	if w.CreatureAt(Loc(5, 5)) != None {
		t.Errorf("player wasn't removed")
	}
	if w.CreatureAt(Loc(4, 5)) != PlayerCreature {
		t.Errorf("player did't move")
	}
}

func TestInBounds(t *testing.T) {
	w := NewRoom(10, 10)
	w.Spawn(5, 5)
	tests := []struct {
		loc  Location
		want bool
	}{
		{
			loc:  Loc(4, 5),
			want: true,
		},
		{
			loc:  Loc(0, 0),
			want: true,
		},
		{
			loc:  Loc(10, 0),
			want: false,
		},
	}
	for _, test := range tests {
		if got := w.InBounds(test.loc); got != test.want {
			t.Errorf("loc %v got %v want %v", test.loc, got, test.want)
		}
	}
}
