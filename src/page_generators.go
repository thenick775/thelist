package main

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func genAddForm(_ fyne.Window) fyne.CanvasObject {
	name := widget.NewEntry()
	name.SetPlaceHolder("Your Item Name")
	name.Validator = validation.NewRegexp(`^.+$`, "identifier must not be empty")

	rating := widget.NewEntry()
	rating.SetPlaceHolder("Item Rating (1-5)")
	rating.Validator = validation.NewRegexp(`^[1-5]{1}$`, "not a valid rating (1-5)")

	tagentry := NewSubmitEntry()
	tagentry.SetPlaceHolder("Enter Tags here")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name, HintText: "The Name of your Item"},
			{Text: "Rating", Widget: rating, HintText: "The Item's rating"},
			{Text: "Tags", Widget: tagentry, HintText: "Enter your tags here to add to list"},
		},
		OnCancel: func() {
			tree.OnSelected(state.currentMenuItem)
		},
		OnSubmit: func() {
			if state.noList || state.currentList == "" {
				dialog.ShowInformation("Information", "No list, no action taken", w)
				return
			}
			intVar, _ := strconv.Atoi(rating.Text)
			lists.Data[state.currentList] = append(lists.Data[state.currentList], ListItem{Name: name.Text, Rating: intVar, Tags: tagentry.Text})
			inquiry.FilterList += fmt.Sprintf("\n%s %s %s", name.Text, rating.Text, tagentry.Text)
			f := fmt.Sprintf("%s %s %s", name.Text, rating.Text, tagentry.Text)
			inquiry.SearchMap[f] = len(lists.Data[state.currentList]) - 1
			if inquiry.LinkageMap != nil { //refresh the linkage map/search map
				lists.RegexSearch(lists.SelectEntry.Text)
			} else { //append shown data to existing baselist
				lists.ShowData.strlist = append(lists.ShowData.strlist, name.Text)
				lists.ShowData.data.Reload()
			}
			name.SetText("")
			name.SetValidationError(nil)
			rating.SetText("")
			rating.SetValidationError(nil)
			tagentry.SetText("")
			tagentry.SetValidationError(nil)
			lists.ListModified = true
		},
	}
	tagentry.currFormFunc = form.OnSubmit

	title := widget.NewLabel("Add")
	intro := widget.NewLabel("Add items to your list here, use the enter key to submit\n")
	intro.Wrapping = fyne.TextWrapWord

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, container.NewPadded(form))
}

func genInquire(_ fyne.Window) fyne.CanvasObject {
	inquiry.InqTitle = widget.NewLabel("Inquire")
	if inquiry.LinkageMap == nil { //need to fix this here
		inquiry.InqIntro = widget.NewLabel("Type your regex query here,\nuse the enter key to filter your list:\n" + state.currentList + ", size: " + strconv.Itoa(len(lists.Data[state.currentList])))
	}
	inquiry.InqIntro.Wrapping = fyne.TextWrapWord

	copyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		strspl := strings.Split(inquiry.ExpandL1.Text, "\n")
		w.Clipboard().SetContent(strspl[1])
	})

	t := container.NewBorder(nil, nil, container.NewVBox(copyButton), nil, container.NewVBox(inquiry.ExpandL1, inquiry.ExpandL2, inquiry.ExpandL3))

	inquiry.InquiryTabs = container.NewAppTabs(
		container.NewTabItem("List", container.NewVScroll(lists.List)),
		container.NewTabItem("Item", t),
	)
	inquiry.InquiryTabs.SetTabLocation(container.TabLocationBottom)

	return container.NewBorder(
		container.NewVBox(inquiry.InqTitle, widget.NewSeparator(), inquiry.InqIntro, lists.SelectEntry), nil, nil, nil, inquiry.InquiryTabs)
}

func genRemove(w fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("Remove")
	intro := widget.NewLabel("Enter information of object to be removed\n")

	name := NewSubmitEntry()
	name.SetPlaceHolder("Your Item Name")
	name.Validator = validation.NewRegexp(`^.+$`, "identifier must not be empty")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name, HintText: "The Name of your Item"},
		},
		OnCancel: func() {
			tree.OnSelected(state.currentMenuItem)
		},
		OnSubmit: func() {
			if state.noList || state.currentList == "" {
				dialog.ShowInformation("Information", "No list, no action taken", w)
				return
			}
			cnf := dialog.NewConfirm("Confirmation", "Are you sure you want to delete?", func(response bool) {
				if response {
					ok := lists.RemoveElementByName(name.Text)
					if !ok {
						dialog.ShowError(fmt.Errorf("Error in list deletion,\nno action taken"), w)
					} else {
						//need to implement
						dialog.ShowInformation("Information", "List Item: "+name.Text+" deleted", w)
						lists.ListModified = true
						for i := range lists.ShowData.strlist { //remove from
							if lists.ShowData.strlist[i] == name.Text {
								lists.ShowData.strlist = append(lists.ShowData.strlist[:i], lists.ShowData.strlist[i+1:]...)
								break
							}
						}
						lists.ShowData.data.Reload()
						name.SetText("")
						name.SetValidationError(nil)
					}
				} else {
					fmt.Println("do not remove elem")
				}
			}, w)
			cnf.SetDismissText("No")
			cnf.SetConfirmText("Yes")
			cnf.Show()
		},
	}
	name.currFormFunc = form.OnSubmit

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, container.NewPadded(form))
}

func genEdit(_ fyne.Window) fyne.CanvasObject {
	var oldloc int
	var item ListItem
	if !state.noList || state.currentList == "" {
		if inquiry.LinkageMap == nil {
			item = lists.Data[state.currentList][lists.SelectEntry.list_loc]
			oldloc = lists.SelectEntry.list_loc
		} else { //use the linkage
			item = lists.Data[state.currentList][inquiry.LinkageMap[lists.SelectEntry.list_loc]]
			oldloc = inquiry.LinkageMap[lists.SelectEntry.list_loc]
		}
	}
	name := widget.NewEntry()
	name.SetPlaceHolder("Your Item Name")
	name.Validator = validation.NewRegexp(`^.+$`, "identifier must not be empty")
	name.SetText(item.Name)

	rating := widget.NewEntry()
	rating.SetPlaceHolder("Item Rating (1-5)")
	rating.Validator = validation.NewRegexp(`^[1-5]{1}$`, "not a valid rating (1-5)")
	rating.SetText(strconv.Itoa(item.Rating))

	tagentry := NewSubmitEntry()
	tagentry.SetPlaceHolder("Enter Tags here")
	tagentry.SetText(item.Tags)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name, HintText: "The Name of your Item"},
			{Text: "Rating", Widget: rating, HintText: "The Item's rating"},
			{Text: "Tags", Widget: tagentry, HintText: "Enter your tags here to add to list"},
		},
		OnCancel: func() {
			name.SetText(item.Name)
			rating.SetText(strconv.Itoa(item.Rating))
			tagentry.SetText(item.Tags)
		},
		OnSubmit: func() {
			if state.noList || state.currentList == "" {
				dialog.ShowInformation("Information", "No list, no action taken", w)
				return
			}
			intVar, _ := strconv.Atoi(rating.Text)
			lists.Data[state.currentList][oldloc] = ListItem{Name: name.Text, Rating: intVar, Tags: tagentry.Text}
			lists.ShowData.data.Reload()
			lists.ListModified = true
		},
	}
	tagentry.currFormFunc = form.OnSubmit

	title := widget.NewLabel("Edit")
	intro := widget.NewLabel("Edit items in your list here, use the enter key to submit\n")
	intro.Wrapping = fyne.TextWrapWord

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, container.NewPadded(form))
}

func genConfEdit(w fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("Edit Configuration")
	intro := widget.NewLabel("View and Edit your local Configuration\n")

	defaultList := widget.NewEntry()
	defaultList.SetText(conf["configuration"].(map[string]interface{})["default list"].(string))
	defaultSelected := widget.NewEntry()
	defaultSelected.SetText(conf["configuration"].(map[string]interface{})["default selected"].(string))
	defaultTheme := widget.NewSelectEntry([]string{"Light", "Dark"})
	defaultTheme.SetText(conf["configuration"].(map[string]interface{})["default theme"].(string))
	localItemFile := NewSubmitEntry()
	localItemFile.Validator = validation.NewRegexp(`^.+$`, "file path must not be empty")
	localItemFile.SetText(conf["configuration"].(map[string]interface{})["local item file"].(string))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Default List", Widget: defaultList, HintText: "Your default List that pulls up immediately"},
			{Text: "Default Open Menu Item", Widget: defaultSelected, HintText: "The default menu item selected"},
			{Text: "Default Theme", Widget: defaultTheme, HintText: "Light or Dark theme"},
			{Text: "Default Local Item File", Widget: localItemFile, HintText: "Absolute Location of your item list file"},
		},
		OnCancel: func() {
			defaultList.SetText(conf["configuration"].(map[string]interface{})["default list"].(string))
			defaultSelected.SetText(conf["configuration"].(map[string]interface{})["default selected"].(string))
			defaultTheme.SetText(conf["configuration"].(map[string]interface{})["default theme"].(string))
			localItemFile.SetText(conf["configuration"].(map[string]interface{})["local item file"].(string))
		},
		OnSubmit: func() {
			cnf := dialog.NewConfirm("Confirmation", "Are you sure you want to edit your configuration?", func(response bool) {
				if response {
					if lists.ListExists(defaultList.Text) {
						conf["configuration"].(map[string]interface{})["default list"] = defaultList.Text
					}

					if isMenuTreeLeaf(defaultSelected.Text) {
						conf["configuration"].(map[string]interface{})["default selected"] = defaultSelected.Text
					}

					if defaultTheme.Text != conf["configuration"].(map[string]interface{})["default theme"] {
						if strings.EqualFold(defaultTheme.Text, "light") {
							a.Settings().SetTheme(theme.LightTheme())
							conf["configuration"].(map[string]interface{})["default theme"] = defaultTheme.Text
							state.currentThemeAlias = defaultTheme.Text
						} else if strings.EqualFold(defaultTheme.Text, "dark") {
							a.Settings().SetTheme(theme.DarkTheme())
							conf["configuration"].(map[string]interface{})["default theme"] = defaultTheme.Text
							state.currentThemeAlias = defaultTheme.Text
						}
					}

					if localItemFile.Text != conf["configuration"].(map[string]interface{})["local item file"] {
						conf["configuration"].(map[string]interface{})["local item file"] = localItemFile.Text
						if state.noList && localItemFile.Text != "" {
							state.noList = false
						}
						//determine whether to load file
						if _, err := os.Stat(localItemFile.Text); err == nil {
							// file exists
							byteValue, err := os.ReadFile(localItemFile.Text)
							if err != nil {
								dialog.ShowError(fmt.Errorf("Failed to read new listing file"), w)
							} else {
								err = json.Unmarshal(byteValue, &lists.Data)
								if err != nil {
									dialog.ShowError(fmt.Errorf("Failed to load new listing file"), w)
								} else {
									if state.currentList == "" {
										if conf["configuration"].(map[string]interface{})["default list"].(string) != "" {
											state.currentList = conf["configuration"].(map[string]interface{})["default list"].(string)
										} else {
											listnames := lists.GetOrderedListNames()
											if len(listnames) > 0 {
												state.currentList = listnames[0]
											}
										}
									}
									lists.Initialize()
									dialog.ShowInformation("Information", "List successfully initialized from file", w)
								}
							}
						}
					}

					write_conf()
				} else {
					defaultList.SetText(conf["configuration"].(map[string]interface{})["default list"].(string))
					defaultSelected.SetText(conf["configuration"].(map[string]interface{})["default selected"].(string))
					defaultTheme.SetText(conf["configuration"].(map[string]interface{})["default theme"].(string))
					localItemFile.SetText(conf["configuration"].(map[string]interface{})["local item file"].(string))
				}
			}, w)
			cnf.SetDismissText("No")
			cnf.SetConfirmText("Yes")
			cnf.Show()
		},
	}
	localItemFile.currFormFunc = form.OnSubmit

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, container.NewPadded(form))
}

func genSwitchList(_ fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("Select List")
	intro := widget.NewLabel("Choose your active list\n")
	keys := lists.GetOrderedListNames()

	radiogr := widget.NewRadioGroup(keys, func(s string) {
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
		}
	})
	radiogr.Horizontal = false
	radiogr.Required = true
	radiogr.SetSelected(state.currentList)

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, container.NewVScroll(radiogr))
}

func genAddList(_ fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("Add List")
	intro := widget.NewLabel("Enter your new list name\n")

	newList := NewSubmitEntry()
	newList.Validator = validation.NewRegexp(`^.+$`, "identifier must not be empty")
	newList.SetPlaceHolder("Enter list to add")
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "New List", Widget: newList, HintText: "Your new list name"},
		},
		OnCancel: func() {
			tree.OnSelected(state.currentMenuItem)
		},
		OnSubmit: func() {
			if state.noList {
				dialog.ShowInformation("Information", "No list, no action taken", w)
				return
			}
			lists.Data[newList.Text] = []ListItem{}
			lists.ListModified = true
			newList.SetText("")
			newList.SetValidationError(nil)
		},
	}
	newList.currFormFunc = form.OnSubmit

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, container.NewPadded(form))
}

func genEditList(_ fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("Edit List")
	intro := widget.NewLabel("Edit current list name\n")

	newList := NewSubmitEntry()
	newList.Validator = validation.NewRegexp(`^.+$`, "identifier must not be empty")
	newList.SetPlaceHolder("Enter new name for list")
	newList.SetText(state.currentList)
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Edit List", Widget: newList, HintText: "Your new list name"},
		},
		OnCancel: func() {
			newList.SetText(state.currentList)
		},
		OnSubmit: func() {
			if state.noList || state.currentList == "" {
				dialog.ShowInformation("Information", "No list, no action taken", w)
				return
			}
			lists.Data[newList.Text] = lists.Data[state.currentList]
			delete(lists.Data, state.currentList)
			lists.ListModified = true
		},
	}
	newList.currFormFunc = form.OnSubmit

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, container.NewPadded(form))
}

func genDeleteList(_ fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("Delete List")
	intro := widget.NewLabel("Enter your list to delete\n")

	delList := widget.NewEntry()
	delList.Validator = validation.NewRegexp(`^.+$`, "identifier must not be empty")
	delList.SetPlaceHolder("Enter list to delete")
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Delete List", Widget: delList, HintText: "Your list to delete"},
		},
		OnCancel: func() {
			tree.OnSelected(state.currentMenuItem)
		},
		OnSubmit: func() {
			if state.noList {
				dialog.ShowInformation("Information", "No list, no action taken", w)
				return
			}
			cnf := dialog.NewConfirm("Confirmation", "Are you sure you want to delete a full list?", func(response bool) {
				if response {
					if _, ok := lists.Data[delList.Text]; ok {
						delete(lists.Data, delList.Text)
						lists.ListModified = true
						if state.currentList == delList.Text {
							key_zero := lists.GetOrderedListNames()[0] //need to work on case deleting last list
							lists.ShowData.strlist = lists.GenListFromMap(key_zero)
							lists.SelectEntry.list_loc = 0
							lists.List.Select(lists.SelectEntry.list_loc)
							lists.ShowData.data.Reload()
							state.currentList = key_zero
						}
						dialog.ShowInformation("Information", "List "+delList.Text+" deleted", w)
					} else {
						dialog.ShowError(fmt.Errorf("List name invalid,\nnothing to delete"), w)
					}
					delList.SetText("")
					delList.SetValidationError(nil)
				}
			}, w)
			cnf.SetDismissText("No")
			cnf.SetConfirmText("Yes")
			cnf.Show()
		},
	}

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, container.NewPadded(form))
}

func genWordCloud(_ fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("Word Cloud (Loading)")
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{1, 1}}) //dummy image

	image := canvas.NewImageFromImage(img)
	image.FillMode = canvas.ImageFillContain

	internalBreakdownList := []string{} //dummy data list
	breakdownShowList := binding.BindStringList(
		&internalBreakdownList,
	)

	dataBreakdown := widget.NewListWithData(breakdownShowList,
		func() fyne.CanvasObject {
			lb := widget.NewLabel("template")
			lb.Truncation = fyne.TextTruncateEllipsis
			return lb
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	exportButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		NewImgExportPop(image.Image)
	})

	imgHolder := container.NewBorder(nil, nil, container.NewVBox(exportButton), nil, image)

	wordCloudTabs := container.NewAppTabs(
		container.NewTabItem("Cloud", imgHolder),
		container.NewTabItem("Data", container.NewVScroll(dataBreakdown)),
	)

	go func() { //async image processing/rendering+data processing
		//gen image and retreive data
		g_img, g_data := genWordCloudImg()
		image.Image = g_img
		title.SetText("Word Cloud")
		title.Refresh()
		image.Refresh()

		sortedBreakdownKeys := make([]string, 0, len(g_data))
		for k := range g_data {
			sortedBreakdownKeys = append(sortedBreakdownKeys, k)
		}
		sort.SliceStable(sortedBreakdownKeys, func(i, j int) bool { //sort keys by value desc
			return g_data[sortedBreakdownKeys[i]] > g_data[sortedBreakdownKeys[j]]
		})

		for _, val := range sortedBreakdownKeys {
			internalBreakdownList = append(internalBreakdownList, "Word: "+val+", Count: "+strconv.Itoa(g_data[val]))
		}
		breakdownShowList.Reload()
	}()

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, wordCloudTabs)
}

func genStatistics(_ fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("Statistics")
	stats := genStats()
	totalLists := len(stats)

	//generate markdown
	mk_down := `## Total Lists: ` + strconv.Itoa(totalLists)
	for idx, stat := range stats {
		mk_down += `
---
## List #` + strconv.Itoa(idx+1) + `: ` + stat.Name + `

---
Total Items: ` + strconv.Itoa(stat.TotalContentCount) + `

Ratings Count:

`
		//display ratings in descending order
		sortedRatings := make([]int, 0, len(stat.ContentCountPerRating))
		for k := range stat.ContentCountPerRating {
			sortedRatings = append(sortedRatings, k)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(sortedRatings)))

		for _, key := range sortedRatings {
			mk_down += `
* Rating: ` + strconv.Itoa(key) + `, Count: ` + strconv.Itoa(stat.ContentCountPerRating[key])
		}
	}
	//end markdown generation
	rich := widget.NewRichTextFromMarkdown(mk_down)

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, container.NewVScroll(rich))
}
