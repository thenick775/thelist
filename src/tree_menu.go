package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func menuTree(w fyne.Window, view *fyne.Container, defaultBranch string, defaultSelected string) *widget.Tree {
	listtree := map[string][]string{
		"":              {"Quick Actions", "My Lists", "Configuration"},
		"Quick Actions": {"Inquire", "Add", "Remove", "Edit"},
		"My Lists":      {"Switch List", "Add List", "Delete List", "Edit List"},
		"Configuration": {"Defaults"},
	}

	pagemap := map[string]MenuPageLink{ //for tree list page navigation/generation
		"Inquire":     MenuPageLink{View: genInquire},
		"Add":         MenuPageLink{View: genAddForm},
		"Remove":      MenuPageLink{View: genRemove},
		"Edit":        MenuPageLink{View: genEdit},
		"Defaults":    MenuPageLink{View: genConfEdit},
		"Switch List": MenuPageLink{View: genSwitchList},
		"Add List":    MenuPageLink{View: genAddList},
		"Edit List":   MenuPageLink{View: genEditList},
		"Delete List": MenuPageLink{View: genDeleteList},
	}

	tree := widget.NewTreeWithStrings(listtree)
	tree.OnSelected = func(uid string) { //here we switch between views
		if page, ok := pagemap[uid]; ok {
			state.currentMenuItem = uid
			view.Objects = []fyne.CanvasObject{page.View(w)}
			view.Refresh()
		} else {
			tree.ToggleBranch(uid)
			tree.Unselect(uid)
		}
	}

	tree.OpenBranch(defaultBranch)
	tree.Select(defaultSelected)

	return tree
}
