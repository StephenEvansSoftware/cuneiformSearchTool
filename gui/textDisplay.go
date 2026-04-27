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
type TextDisplayPane interface {
	PaneGetter
	SetTitle(string)
	ParseMarkdown([]string)
	AppendMarkdown([]string)
}

type textDisplayPane struct {
	mainContainer *fyne.Container

	labelText *widget.Label

	richTextBox *widget.RichText
}

// TODO: tidy this up, move things into their own functions when needed, etc
func newTextDisplay() TextDisplayPane {
	tdp := &textDisplayPane{}
	tdp.labelText = widget.NewLabel(labelTextDisplayDefault)
	tdp.richTextBox = widget.NewRichTextFromMarkdown(richTextDisplayDefault)

	tdp.mainContainer = container.NewVBox(
		tdp.labelText,
		widget.NewSeparator(),
		tdp.richTextBox,
	)
	
	return tdp
}

func (tdp *textDisplayPane) GetPane() fyne.CanvasObject {
	return tdp.mainContainer
}

func (tdp *textDisplayPane) SetTitle(title string) {
	tdp.labelText.SetText(title)
}

// Simple passthrough without any formality
func (tdp *textDisplayPane) ParseMarkdown(content []string) {
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

func (tdp *textDisplayPane) AppendMarkdown(content []string) {
	for _, line := range content {
		tdp.richTextBox.AppendMarkdown(line)
	}
}
