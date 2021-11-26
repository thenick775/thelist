package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func genKeyBoardShortcutPopup() {
	helppop := a.NewWindow("Help")
	helppop.SetContent(container.NewVScroll(container.NewVBox(
		widget.NewLabel("Keyboard Shortcuts"),
		widget.NewSeparator(),
		widget.NewLabel("Super+F -> Focus on Inquiry Input Field"),
		widget.NewLabel("Super+B -> Clear the Inquiry Input Field"),
		widget.NewLabel("Super+G -> Open Add View"),
		widget.NewLabel("Super+E -> Open Edit View"),
		widget.NewLabel("Super+I -> Open Inquiry View"),
		widget.NewLabel("Super+R -> Open Switch List View"),
		widget.NewLabel("Super+Up -> Switch List (alphabetically up)"),
		widget.NewLabel("Super+Down -> Switch List (alphabetically down)"),
	)))
	helppop.Resize(fyne.NewSize(210, 275))

	if deskCanvas, ok := helppop.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(key *fyne.KeyEvent) {
			if key.Name == fyne.KeyEscape {
				helppop.Close()
			}
		})
	}
	helppop.Show()
}

//used for export to csv and to json
func NewExportPop(filetype string) {
	exp := a.NewWindow("Export " + filetype)
	fullexport := false
	cancel := widget.NewButton("Cancel", func() {
		exp.Close()
	})
	fname := widget.NewEntry()
	fname.SetPlaceHolder("File name")
	submit := widget.NewButton("Submit", func() {
		switch filetype {
		case "CSV":
			write_csv(fullexport, fname.Text)
		case "JSON":
			write_json(fullexport, fname.Text)
		}
	})
	submit.Disable()
	fname.Validator = validation.NewRegexp(`^.+$`, "file name cannot be empty")
	fname.SetOnValidationChanged(func(err error) {
		if err != nil {
			submit.Disable()
		} else {
			submit.Enable()
		}
	})
	check := widget.NewCheck("Full Data", func(value bool) {
		fullexport = value
	})

	exp.SetContent(container.NewVScroll(container.NewVBox(
		check,
		fname,
		container.NewHBox(submit, cancel),
	)))
	exp.Resize(fyne.NewSize(400, 150))

	if deskCanvas, ok := exp.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(key *fyne.KeyEvent) {
			if key.Name == fyne.KeyEscape {
				exp.Close()
			}
		})
	}
	exp.Show()
}
