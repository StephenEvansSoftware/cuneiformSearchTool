package gui

import (
	"fyne.io/fyne/v2"
)

// Simple interface for individual elements of the gui

type paneGetter interface {
	getPane() fyne.CanvasObject
}
