package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/arrufat/papago"
)

func makeTranslatePage() ui.Control {
	vbox := ui.NewVerticalBox()
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	// source language group
	group := ui.NewGroup("Source Language")
	group.SetMargined(true)

	sourceLangCBox := ui.NewCombobox()
	targetLangCBox := ui.NewCombobox()
	for _, lang := range papago.SupportedLanguages() {
		sourceLangCBox.Append(lang.String())
		targetLangCBox.Append(lang.String())
	}
	sourceLangCBox.SetSelected(1)
	targetLangCBox.SetSelected(0)
	smle := ui.NewMultilineEntry()
	group.SetChild(sourceLangCBox)
	vbox.Append(group, false)
	vbox.Append(smle, true)

	hbox.Append(vbox, true)

	// swap button
	swapBtn := ui.NewButton(" ⇄ ")
	grid := ui.NewGrid()
	grid.SetPadded(true)
	grid.Append(swapBtn, 0, 0, 2, 1, false, ui.AlignCenter, false, ui.AlignCenter)
	translateBtn := ui.NewButton("Translate")
	grid.Append(translateBtn, 0, 1, 2, 1, false, ui.AlignCenter, false, ui.AlignCenter)
	hbox.Append(grid, false)

	// target language group
	vbox = ui.NewVerticalBox()
	group = ui.NewGroup("Target Language")
	group.SetMargined(true)

	tmle := ui.NewMultilineEntry()
	group.SetChild(targetLangCBox)
	vbox.Append(group, false)
	vbox.Append(tmle, true)

	hbox.Append(vbox, true)

	translateCallback := func() {
		if smle.Text() == "" {
			return
		}
		var sourceLang, targetLang papago.Language
		sourceLangIdx := sourceLangCBox.Selected()
		if sourceLangIdx == -1 {
			detLang, err := papago.Detect(smle.Text())
			if err != nil {
				log.Fatal(err)
			}
			sourceLang = detLang
		} else {
			sourceLang = papago.SupportedLanguages()[sourceLangIdx]
		}
		targetLangIdx := targetLangCBox.Selected()
		if targetLangIdx == -1 {
			log.Println("Target language not selected, skipping...")
			return
		}
		targetLang = papago.SupportedLanguages()[targetLangIdx]
		if sourceLang == targetLang {
			log.Println("Source and target languages are the same, skipping...")
			return
		}
		text := strings.Replace(smle.Text(), "\n", "\\n", -1)
		fmt.Printf("%q\n", text)
		translation, err := papago.Translate(text, sourceLang, targetLang, papago.TranslateOptions{})
		if err != nil {
			log.Fatal(err)
		}
		tmle.SetText(translation)
	}

	targetLangCBox.OnSelected(func(*ui.Combobox) {
		translateCallback()
	})

	translateBtn.OnClicked(func(*ui.Button) {
		translateCallback()
	})

	smle.OnChanged(func(*ui.MultilineEntry) {
		// if smle.Text() == "" || strings.HasSuffix(smle.Text(), "\n") {
		// translateCallback()
		// }
	})

	swapBtn.OnClicked(func(*ui.Button) {
		auxLangIdx := sourceLangCBox.Selected()
		auxText := smle.Text()
		sourceLangCBox.SetSelected(targetLangCBox.Selected())
		targetLangCBox.SetSelected(auxLangIdx)
		smle.SetText(tmle.Text())
		tmle.SetText(auxText)
	})

	return hbox
}

// func makeSettingsPage() ui.Control {

// }

func setupUI() {
	mainwin := ui.NewWindow("NAVER Papago Translate", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)
	tab.Append("Translation", makeTranslatePage())
	tab.SetMargined(0, true)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
