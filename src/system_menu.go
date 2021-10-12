package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/theme"
	"net/url"
)

//system menu setup, this is the "external" system menu
func setupSystemMenu(w fyne.Window, a fyne.App) {
	newItem := fyne.NewMenuItem("New", nil)
	newItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") }),
		fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") }),
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
	sortItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Alpha ASC", func() { fmt.Println("File ASC sort") }),
		fyne.NewMenuItem("Alpha DESC", func() { fmt.Println("File DESC sort") }),
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
