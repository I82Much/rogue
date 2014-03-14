package player

import (
	"testing"
)

func TestDamageKillsPlayer(t *testing.T) {
	p := NewPlayer(100)
	if p.IsDead() {
		t.Errorf("should not be dead")
	}
	p.Damage(99)
	if p.Life != 1 {
		t.Errorf("expected one life left")
	}
	if p.IsDead() {
		t.Errorf("should not be dead after inflicting 99 damage")
	}
	p.Damage(1)
	if !p.IsDead() {
		t.Errorf("should be dead after inflicting 1 last point of damage")
	}
}
