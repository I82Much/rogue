package rogue

// The controller accepts input and converts it into commands for the model and view

type CombatController struct {
	model *CombatModel
	view  *CombatView
}
