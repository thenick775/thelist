package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"strconv"
)

// menu shortcuts
func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

// shared shortcuts below
func superFind(shortcut fyne.Shortcut) {
	if state.currentMenuItem == "Inquire" {
		w.Canvas().Focus(lists.SelectEntry) //maybe need unfocus somewhere?
	}
}

func superAdd(shortcut fyne.Shortcut) {
	tree.Select("Add")
}

func superEdit(shortcut fyne.Shortcut) {
	tree.Select("Edit")
}

func superInquire(shortcut fyne.Shortcut) {
	tree.Select("Inquire")
}

func superSwitchList(shortcut fyne.Shortcut) {
	tree.Select("Switch List")
}

func superSwitchListUp(shortcut fyne.Shortcut) {
	keys, s := lists.GetOrderedListNames(), ""

	for idx, val := range keys {
		if val == state.currentList {
			if idx-1 >= 0 {
				s = keys[idx-1]
				break
			} else {
				return
			}
		}
	}

	state.currentList = s
	lists.SelectEntry.SetText("")
	if state.alphasort.enabled {
		lists.GenListFromMap(s)
		lists.RegexSearch("")
	} else {
		lists.ShowData.strlist = lists.GenListFromMap(s)
		lists.ShowData.data.Reload()
		lists.SelectEntry.list_loc = 0
		lists.List.Select(lists.SelectEntry.list_loc)
		inquiryIndexAndExpand(0)
		inquiry.InqIntro.SetText("Type your regex query here,\nuse the enter key to filter your list:\n" + state.currentList + ", size: " + strconv.Itoa(len(lists.Data[state.currentList])))
	}
}

func superSwitchListDown(shortcut fyne.Shortcut) {
	keys, s := lists.GetOrderedListNames(), ""

	for idx, val := range keys {
		if val == state.currentList {
			if idx+1 < len(keys) {
				s = keys[idx+1]
				break
			} else {
				return
			}
		}
	}

	state.currentList = s
	lists.SelectEntry.SetText("")
	if state.alphasort.enabled {
		lists.GenListFromMap(s)
		lists.RegexSearch("")
	} else {
		lists.ShowData.strlist = lists.GenListFromMap(s)
		lists.ShowData.data.Reload()
		lists.SelectEntry.list_loc = 0
		lists.List.Select(lists.SelectEntry.list_loc)
		inquiryIndexAndExpand(0)
		inquiry.InqIntro.SetText("Type your regex query here,\nuse the enter key to filter your list:\n" + state.currentList + ", size: " + strconv.Itoa(len(lists.Data[state.currentList])))
	}
}

func superClearInquiry(shortcut fyne.Shortcut) {
	lists.SelectEntry.SetText("")
}

func setupDesktopShortcuts(w fyne.Window) {
	ctrlFind := desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlFind, superFind)
	ctrlAdd := desktop.CustomShortcut{KeyName: fyne.KeyG, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlAdd, superAdd)
	ctrlEdit := desktop.CustomShortcut{KeyName: fyne.KeyE, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlEdit, superEdit)
	ctrlInquire := desktop.CustomShortcut{KeyName: fyne.KeyI, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlInquire, superInquire)
	ctrlSwitchList := desktop.CustomShortcut{KeyName: fyne.KeyR, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlSwitchList, superSwitchList)
	ctrlSwitchListUp := desktop.CustomShortcut{KeyName: fyne.KeyUp, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlSwitchListUp, superSwitchListUp)
	ctrlSwitchListDown := desktop.CustomShortcut{KeyName: fyne.KeyDown, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlSwitchListDown, superSwitchListDown)
	ctrlClearInquiry := desktop.CustomShortcut{KeyName: fyne.KeyB, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlClearInquiry, superClearInquiry)
}
