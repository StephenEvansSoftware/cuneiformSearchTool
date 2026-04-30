package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type massSearchResultsPane interface {
	paneGetter
}

type msrpStruct struct {
	mainContainer *fyne.Container
}

func newMassSearchResultsPane() massSearchResultsPane {
	msrPane := &msrpStruct{}

	msrPane.mainContainer = container.NewVBox(
		widget.NewLabel("mass search results pane"),
	)

	return msrPane
}

func (msrp *msrpStruct) getPane() fyne.CanvasObject {
	return msrp.mainContainer
}
