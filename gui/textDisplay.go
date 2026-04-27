package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// This file creates a markdown text display specifically for showing cuneiform texts

const (
	labelTextDisplayDefault = "No file loaded."

	richTextDisplayDefault = "*No file loaded.*"
)

// TODO: better name?
type textDisplayPane interface {
	paneGetter
	SetTitle(string)
	ParseMarkdown([]string)
	AppendMarkdown([]string)
	// TODO: full method to set title and set markdown
}

type tdpStruct struct {
	mainContainer *fyne.Container

	labelText *widget.Label

	richTextBox *widget.RichText
}

// TODO: tidy this up, move things into their own functions when needed, etc
func newTextDisplay() textDisplayPane {
	tdp := &tdpStruct{}
	tdp.labelText = widget.NewLabel(labelTextDisplayDefault)
	tdp.richTextBox = widget.NewRichTextFromMarkdown(richTextDisplayDefault)


	tdp.mainContainer = container.NewBorder(
		container.NewVBox(
			tdp.labelText,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		container.NewScroll(tdp.richTextBox),
	)

	return tdp
}

func (tdp *tdpStruct) GetPane() fyne.CanvasObject {
	return tdp.mainContainer
}

func (tdp *tdpStruct) SetTitle(title string) {
	tdp.labelText.SetText(title)
}

// Simple passthrough without any formality
// TODO: change to string and optional additional strings
func (tdp *tdpStruct) ParseMarkdown(content []string) {
	if len(content) == 0 {
		// No text input, just parse empty string
		tdp.richTextBox.ParseMarkdown("")
		return
	}
	// Parse the first line
	tdp.richTextBox.ParseMarkdown(content[0])
	// Append all others
	tdp.AppendMarkdown(content[1:])
}

func (tdp *tdpStruct) AppendMarkdown(content []string) {
	for _, line := range content {
		tdp.richTextBox.AppendMarkdown(line)
	}
}
