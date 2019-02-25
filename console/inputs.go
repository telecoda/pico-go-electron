package console

import (
	"github.com/hajimehoshi/ebiten"
)

type inputs struct {
	mouseDown bool
}

func newInputs() *inputs {
	return &inputs{}
}

func (i *inputs) Btn(id int) bool {
	return false
}

func (i *inputs) MouseClicked() bool {
	if i == nil {
		return false
	}

	currentState := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	var isClicked bool

	if currentState && !i.mouseDown {
		// already pressed
		isClicked = true
	}

	i.mouseDown = currentState

	return isClicked
}

func (i *inputs) MousePressed() bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

func (i *inputs) MousePosition() (int, int) {
	return ebiten.CursorPosition()
}
