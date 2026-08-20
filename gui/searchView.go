package gui

import (
	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/driver/desktop"
	//	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"fmt"
	//"time"
	//"cuneiformSearchTool/util"
)

const (
	tabItemNewSearchLabel = "New Search"
)

type searchView interface {
	// startNewTermsSearch([][]string, *util.SearchSettings, time.Time)
}

type svStruct struct {
	sidebar                   *fyne.Container // TODO TODO TODO: make this into a widget.Toolbar
	sidebarButtonAddNewSearch *widget.Button
	sidebarButtonRemoveSearch *widget.Button
	sidebarButtonSaveSearches *widget.Button
	sidebarButtonExportSearch *widget.Button

	searchPane *container.DocTabs

	gui guiInternal
}

func newSearchView(gui guiInternal) fyne.CanvasObject {
	sv := &svStruct{gui: gui}
	sv.setupSidebar()
	sv.setupSearchPane()
	return container.NewBorder(
		nil,
		nil,
		container.NewHBox(
			sv.sidebar,
			widget.NewSeparator(),
		),
		nil,
		sv.searchPane,
	)
}

/*
	type testButton struct {
		widget.Button
		helpText string
	}

	func (tb *testButton) MouseMoved(me *desktop.MouseEvent) {
		fmt.Println("mouse moved")
		fmt.Println(tb.helpText)
		fmt.Println(me.AbsolutePosition)
		fmt.Println(me.Position)
	}
*/

func (sv *svStruct) setupSidebar() {

	/*
		tb := &testButton{helpText: "help text here"}
		tb.Icon = theme.ContentAddIcon()
		tb.ExtendBaseWidget(tb)
	*/
	//	sv.sidebarButtonAddNewSearch = tb
	//sv.sidebarButtonAddNewSearch.Icon = theme.ContentAddIcon()

	sv.sidebarButtonAddNewSearch = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})
	sv.sidebarButtonRemoveSearch = widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {})
	sv.sidebarButtonSaveSearches = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {})
	sv.sidebarButtonExportSearch = widget.NewButtonWithIcon("", theme.DownloadIcon(), func() {})

	sv.sidebarButtonAddNewSearch.OnTapped = func() {
		sv.searchPane.Append(sv.createNewSearchTab())
	}

	sv.sidebar = container.NewVBox(
		sv.sidebarButtonAddNewSearch,
		widget.NewSeparator(),
		sv.sidebarButtonRemoveSearch,
		widget.NewSeparator(),
		sv.sidebarButtonSaveSearches,
		widget.NewSeparator(),
		sv.sidebarButtonExportSearch,
	)
}

func (sv *svStruct) setupSearchPane() {
	sv.searchPane = container.NewDocTabs(
		sv.createNewSearchTab(),
	)

	sv.searchPane.CreateTab = func() *container.TabItem {
		return sv.createNewSearchTab()
	}

	sv.searchPane.SetTabLocation(container.TabLocationLeading)

	/*
		sv.searchPane.CloseIntercept = func(ti *container.TabItem) {
			sv.searchPane.Append(sv.createNewSearchTab())
			fyne.Do(func() {
				sv.searchPane.Refresh()
			})
		}
	*/
}

func (sv *svStruct) createNewSearchTab() *container.TabItem {
	/*
		srd := createNewSearchResultDisplay()
		sr := createNewSearchResults(srd)
		scp := newSearchConfPane(sr)
	*/
	sr := newSearchResultsPane(sv.gui)
	scp := newSearchConfPane(sr)

	return container.NewTabItem(
		tabItemNewSearchLabel,
		container.NewBorder(
			nil, nil,
			container.NewHBox(
				scp.getPane(),
				widget.NewSeparator(),
			), nil,
			sr.getPane(),
		),
	)
}

/*
func (sv *svStruct) startNewTermsSearch([][]string, *util.SearchSettings, time.Time) {
	fmt.Println("sv terms search")
}
*/
