package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"os"
)

//system menu setup, this is the "external" system menu
func setupSystemMenu(w fyne.Window, a fyne.App) {
	newItem := fyne.NewMenuItem("New", nil)
	newItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("File", func() {
			dirpath := widget.NewEntry()
			dirpath.Validator = validation.NewRegexp(`^.+$`, "file name must not be empty")
			items := []*widget.FormItem{
				widget.NewFormItem("File name", dirpath),
			}
			dialog.ShowForm("New File", "Submit", "Cancel", items, func(b bool) {
				if b {
					empty, err := os.Create(dirpath.Text)
					if err != nil {
						dialog.ShowError(fmt.Errorf("Failed to create new file"), w)
					} else {
						dialog.ShowInformation("Information", "File successfully created", w)
						empty.Close()
					}
				}
			}, w)
		}),
		fyne.NewMenuItem("Directory", func() {
			dirpath := widget.NewEntry()
			dirpath.Validator = validation.NewRegexp(`^.+$`, "dir path must not be empty")
			items := []*widget.FormItem{
				widget.NewFormItem("Directory path", dirpath),
			}
			dialog.ShowForm("New Directory", "Submit", "Cancel", items, func(b bool) {
				if b {
					err := os.MkdirAll(dirpath.Text, os.ModePerm)
					if err != nil {
						dialog.ShowError(fmt.Errorf("Failed to create Directory"), w)
					} else {
						dialog.ShowInformation("Information", "Directory successfully created", w)
					}
				}
			}, w)
		}),
	)
	settingsItem := fyne.NewMenuItem("Settings", func() {
		w := a.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	})

	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(&fyne.ShortcutCut{
			Clipboard: w.Clipboard(),
		}, w)
	})
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(&fyne.ShortcutCopy{
			Clipboard: w.Clipboard(),
		}, w)
	})
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(&fyne.ShortcutPaste{
			Clipboard: w.Clipboard(),
		}, w)
	})
	findItem := fyne.NewMenuItem("Find", func() { fmt.Println("Menu Find") })

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://github.com/thenick775/thelist")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("Shortcut Keys", func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Support", func() {}),
		fyne.NewMenuItem("Sponsor", func() {}))

	themeMenu := fyne.NewMenu("Theme",
		fyne.NewMenuItem("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		fyne.NewMenuItem("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}))

	sortItem := fyne.NewMenuItem("Sort", nil)
	sortItem.ChildMenu = fyne.NewMenu("", fyne.NewMenuItem("Alpha ASC", func() {
		if !state.alphasort.enabled {
			state.alphasort.enabled = true
		}
		state.alphasort.order = 0
		lists.RegexSearch(lists.SelectEntry.Text)
	}),
		fyne.NewMenuItem("Alpha DESC", func() {
			if !state.alphasort.enabled {
				state.alphasort.enabled = true
			}
			state.alphasort.order = 1
			lists.RegexSearch(lists.SelectEntry.Text)
		}),
		fyne.NewMenuItem("Enable/Disable", func() {
			state.alphasort.enabled = !state.alphasort.enabled
			lists.RegexSearch(lists.SelectEntry.Text)
		}),
	)
	dataMenu := fyne.NewMenu("Data", sortItem)

	file := fyne.NewMenu("File", newItem)
	if !fyne.CurrentDevice().IsMobile() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		themeMenu,
		dataMenu,
		helpMenu,
	)
	w.SetMainMenu(mainMenu)
	w.SetMaster()
}
