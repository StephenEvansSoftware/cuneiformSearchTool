package gui

import (
	"fmt"
	"slices"

	"fyne.io/fyne/v2"

	"cuneiformSearchTool/util"
)

const (
	menuTextOptions = "Options"

	menuTextPrintSettings             = "Print Settings"
	menuItemTextPrintSettingsMetaInfo = "      Meta Information"
	menuItemTextPrintSettingsUsePos   = "      Use POS over EPOS"
	menuItemTextPrintSettingsPosEpos  = "      Show POS / EPOS"
	menuItemTextPrintSettingsForm     = "      Form over frag"
	menuItemTextPrintSettingsCf       = "      Cf"

	menuTextView           = "View"
	menuItemTextViewSearch = "Search"
	menuItemTextViewCorpus = "Corpus"
)

type mainMenu interface {
	getMainMenu() *fyne.MainMenu
	getPrintSettings() *util.PrintSettings
}

type mmStruct struct {
	mainMenu *fyne.MainMenu

	// Menus
	optionsMenu *fyne.Menu

	printOptionsMenu     *fyne.Menu
	printOptionsMetaInfo *fyne.MenuItem
	printOptionsUsePos   *fyne.MenuItem
	printOptionsPosEpos  *fyne.MenuItem
	printOptionsForm     *fyne.MenuItem
	printOptionsCf       *fyne.MenuItem

	viewMenu   *fyne.Menu
	viewSearch *fyne.MenuItem
	viewCorpus *fyne.MenuItem

	// Reference to parent to pass information back and forth
	gui guiInternal
}

func newMainMenu(gui guiInternal) mainMenu {
	mm := &mmStruct{gui: gui}

	mm.setupOptionsMenu()
	mm.setupPrintSettingsMenu()
	mm.setupViewMenu()

	mm.mainMenu = fyne.NewMainMenu(
		mm.optionsMenu,
		mm.printOptionsMenu,
		mm.viewMenu,
	)

	return mm
}

func (mm *mmStruct) setupOptionsMenu() {
	// TODO: quick hacky way to print out the corpus meta info, fix
	optionsPrintMetaInfo := fyne.NewMenuItem("Print Corpus Meta Info", func() {
		corpusOrder := []string{}
		corpusMetaInfo := map[string]*util.CorpusMetaInfo{}
		for _, ci := range util.ListCorpusAndTexts() {
			corpusOrder = append(corpusOrder, ci.CorpusName)
			cmi, err := util.GetCorpusMetaInfo(ci.CorpusName, false, false)
			if err != nil {
				panic(err.Error()) // TODO: handle better
			}
			corpusMetaInfo[ci.CorpusName] = cmi
		}

		for _, co := range corpusOrder {
			fmt.Printf(",%s", co)
		}
		fmt.Printf("\nNumberOfTexts")
		for _, co := range corpusOrder {
			fmt.Printf(",%d", corpusMetaInfo[co].NumberOfTexts)
		}
		fmt.Printf("\nTotalFragsInAllTexts")
		for _, co := range corpusOrder {
			fmt.Printf(",%d", corpusMetaInfo[co].TotalFragsInAllTexts)
		}
		fmt.Printf("\nTotalFragsInAllTextsWithoutX")
		for _, co := range corpusOrder {
			fmt.Printf(",%d", corpusMetaInfo[co].TotalFragsInAllTextsWithoutX)
		}
		fmt.Printf("\nTotalLinesInAllTexts")
		for _, co := range corpusOrder {
			fmt.Printf(",%d", corpusMetaInfo[co].TotalLinesInAllTexts)
		}
		fmt.Printf("\nAverageLineLength")
		for _, co := range corpusOrder {
			fmt.Printf(",%.2f", corpusMetaInfo[co].AverageLineLength)
		}
		fmt.Printf("\nAverageLineLengthWithoutBlanks")
		for _, co := range corpusOrder {
			fmt.Printf(",%.2f", corpusMetaInfo[co].AverageLineLengthWithoutBlanks)
		}
		fmt.Println()
	})

	// TODO: another quick hacky way to get the data fast, fix later
	optionsPrintMetaInfoCf := fyne.NewMenuItem("Print Corpus CF Info", func() {
		corpusOrder := []string{}
		corpusMetaInfo := map[string]*util.CorpusMetaInfo{}
		for _, ci := range util.ListCorpusAndTexts() {
			corpusOrder = append(corpusOrder, ci.CorpusName)
			cmi, err := util.GetCorpusMetaInfo(ci.CorpusName, false, false)
			if err != nil {
				panic(err.Error()) // TODO: handle better
			}
			corpusMetaInfo[ci.CorpusName] = cmi
		}
		cfOrder := []string{}
		for _, co := range corpusOrder {
			fmt.Printf(",%s", co)
			for key := range corpusMetaInfo[co].CfCounts {
				cfOrder = append(cfOrder, key)
			}
		}
		slices.Sort(cfOrder)
		cfOrder = slices.Compact(cfOrder)

		for _, cf := range cfOrder {
			fmt.Printf("\n%s", cf)
			for _, co := range corpusOrder {
				fmt.Printf(",%d", corpusMetaInfo[co].CfCounts[cf])
			}
		}
		fmt.Println()
	})

	// TODO: same again, just getting data quickly in a gross way
	optionsPrintMetaInfoCfSansPN := fyne.NewMenuItem("Print Corpus CF Info Sans PN", func() {
		corpusOrder := []string{}
		corpusMetaInfo := map[string]*util.CorpusMetaInfo{}
		for _, ci := range util.ListCorpusAndTexts() {
			corpusOrder = append(corpusOrder, ci.CorpusName)
			cmi, err := util.GetCorpusMetaInfo(ci.CorpusName, true, false)
			if err != nil {
				panic(err.Error()) // TODO: handle better
			}
			corpusMetaInfo[ci.CorpusName] = cmi
		}
		cfOrder := []string{}
		for _, co := range corpusOrder {
			fmt.Printf(",%s", co)
			for key := range corpusMetaInfo[co].CfCounts {
				cfOrder = append(cfOrder, key)
			}
		}
		slices.Sort(cfOrder)
		cfOrder = slices.Compact(cfOrder)

		for _, cf := range cfOrder {
			fmt.Printf("\n%s", cf)
			for _, co := range corpusOrder {
				fmt.Printf(",%d", corpusMetaInfo[co].CfCounts[cf])
			}
		}
		fmt.Println()
	})

	// TODO: same again, just getting data quickly in a gross way
	optionsPrintMetaInfoCfSansPNWithType := fyne.NewMenuItem("Print Corpus CF Info Sans PN With Type", func() {
		corpusOrder := []string{}
		corpusMetaInfo := map[string]*util.CorpusMetaInfo{}
		cfTypes := map[string][]string{}
		for _, ci := range util.ListCorpusAndTexts() {
			corpusOrder = append(corpusOrder, ci.CorpusName)
			cmi, err := util.GetCorpusMetaInfo(ci.CorpusName, true, true)
			if err != nil {
				panic(err.Error()) // TODO: handle better
			}
			corpusMetaInfo[ci.CorpusName] = cmi
			for sk, sv := range cmi.CfTypes {
				cfTypes[sk] = append(cfTypes[sk], sv...)
			}
		}
		cfOrder := []string{}
		fmt.Printf(",Type")
		for _, co := range corpusOrder {
			fmt.Printf(",%s", co)
			for key := range corpusMetaInfo[co].CfCounts {
				cfOrder = append(cfOrder, key)
			}
		}
		slices.Sort(cfOrder)
		cfOrder = slices.Compact(cfOrder)

		for _, cf := range cfOrder {
			slices.Sort(cfTypes[cf])
			typeCounts := map[string]int{}
			for _, s := range cfTypes[cf] {
				typeCounts[s] += 1
			}
			cfTypes[cf] = slices.Compact(cfTypes[cf])
			cfTypeCountString := "["
			for _, cft := range cfTypes[cf] {
				cfTypeCountString = fmt.Sprintf("%s %s:%d", cfTypeCountString, cft, typeCounts[cft])
			}
			cfTypeCountString = cfTypeCountString + " ]"
			fmt.Printf("\n%s,%s", cf, cfTypeCountString)
			for _, co := range corpusOrder {
				fmt.Printf(",%d", corpusMetaInfo[co].CfCounts[cf])
			}
		}
		fmt.Println()
	})

	mm.optionsMenu = fyne.NewMenu(
		menuTextOptions,
		optionsPrintMetaInfo,
		optionsPrintMetaInfoCf,
		optionsPrintMetaInfoCfSansPN,
		optionsPrintMetaInfoCfSansPNWithType,
	)
}

func (mm *mmStruct) setupPrintSettingsMenu() {
	mm.printOptionsMetaInfo = fyne.NewMenuItem(menuItemTextPrintSettingsMetaInfo, func() {
		mm.printOptionsMetaInfo.Checked = !mm.printOptionsMetaInfo.Checked
	})
	mm.printOptionsUsePos = fyne.NewMenuItem(menuItemTextPrintSettingsUsePos, func() {
		mm.printOptionsUsePos.Checked = !mm.printOptionsUsePos.Checked
	})
	mm.printOptionsPosEpos = fyne.NewMenuItem(menuItemTextPrintSettingsPosEpos, func() {
		mm.printOptionsPosEpos.Checked = !mm.printOptionsPosEpos.Checked
	})
	mm.printOptionsForm = fyne.NewMenuItem(menuItemTextPrintSettingsForm, func() {
		mm.printOptionsForm.Checked = !mm.printOptionsForm.Checked
	})
	mm.printOptionsCf = fyne.NewMenuItem(menuItemTextPrintSettingsCf, func() {
		mm.printOptionsCf.Checked = !mm.printOptionsCf.Checked
	})

	mm.printOptionsMenu = fyne.NewMenu(
		menuTextPrintSettings,
		mm.printOptionsMetaInfo,
		mm.printOptionsUsePos,
		mm.printOptionsPosEpos,
		mm.printOptionsForm,
		mm.printOptionsCf,
	)
}

func (mm *mmStruct) setupViewMenu() {
	// TODO: disable the menu item for views once the view is loaded
	mm.viewSearch = fyne.NewMenuItem(menuItemTextViewSearch, func() { mm.gui.setViewSearch() })
	mm.viewCorpus = fyne.NewMenuItem(menuItemTextViewCorpus, func() { mm.gui.toggleCorpusPane() })

	mm.viewMenu = fyne.NewMenu(
		menuTextView,
		mm.viewSearch,
		mm.viewCorpus,
	)
}

func (mm *mmStruct) getMainMenu() *fyne.MainMenu {
	return mm.mainMenu
}

func (mm *mmStruct) getPrintSettings() *util.PrintSettings {
	return &util.PrintSettings{
		PrintTextMetaInformation: mm.printOptionsMetaInfo.Checked,
		UsePos:                   mm.printOptionsUsePos.Checked,
		PrintPosEpos:             mm.printOptionsPosEpos.Checked,
		PrintForm:                mm.printOptionsForm.Checked,
		PrintCf:                  mm.printOptionsCf.Checked,
	}
}
