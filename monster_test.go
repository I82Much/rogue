package rogue

import (
	"testing"
)

func TestDamageKills(t *testing.T) {
	m := NewMonster(100)
	if m.IsDead() {
		t.Errorf("should not be dead")
	}
	m.Damage(99)
	if m.Life != 1 {
		t.Errorf("expected one life left")
	}
	if m.IsDead() {
		t.Errorf("should not be dead after inflicting 99 damage")
	}
	m.Damage(1)
	if !m.IsDead() {
		t.Errorf("should be dead after inflicting 1 last point of damage")
	}
}