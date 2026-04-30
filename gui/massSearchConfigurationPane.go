package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type massSearchConfigurationPane interface {
	paneGetter
}

type mscpStruct struct {
	mainContainer *fyne.Container

	gui guiInternal
}

func newMassSearchConfigurationPane(gui guiInternal) massSearchConfigurationPane {
	mscPane := &mscpStruct{gui: gui}

	mscPane.mainContainer = container.NewVBox(
		widget.NewLabel("mass search config pane"),
	)

	return mscPane
}

func (mscp *mscpStruct) getPane() fyne.CanvasObject {
	return mscp.mainContainer
}
