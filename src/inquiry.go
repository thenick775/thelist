package main

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func (i *Inquiry) Initialize() {
	i.ExpandL1 = widget.NewLabel("Name: ")
	i.ExpandL1.Wrapping = fyne.TextWrapWord
	i.ExpandL2 = widget.NewLabel("Rating: ")
	i.ExpandL2.Wrapping = fyne.TextWrapWord
	i.ExpandL3 = widget.NewLabel("Tags: ")
	i.ExpandL3.Wrapping = fyne.TextWrapWord
}

func (l *userList) Initialize() {
	l.ListModified = false
	if state.alphasort.enabled {
		l.ShowData = listData{strlist: []string{}}
		l.ShowData.data = binding.BindStringList(
			&l.ShowData.strlist,
		)
		lists.GenListFromMap(state.currentList)
	} else {
		l.ShowData = listData{strlist: lists.GenListFromMap(state.currentList)} //gaurd this?
		l.ShowData.data = binding.BindStringList(
			&l.ShowData.strlist,
		)
	}

	l.SelectEntry = newInquiryEntry()
	l.SelectEntry.PlaceHolder = "Type your regular expression"

	l.List = widget.NewListWithData(l.ShowData.data,
		func() fyne.CanvasObject {
			lb := widget.NewLabel("template")
			lb.Truncation = fyne.TextTruncateEllipsis
			return lb
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	l.List.OnSelected = inquiryIndexAndExpand

	if state.alphasort.enabled {
		inquiry.InqIntro = widget.NewLabel("Type your regex query here,\nuse the enter key to filter your list")
		inquiry.InqIntro.Wrapping = fyne.TextWrapWord
		l.RegexSearch("")
	}
}

func newInquiryEntry() *inquiryEntry {
	entry := &inquiryEntry{}
	entry.ExtendBaseWidget(entry)
	entry.list_loc = 0
	return entry
}

// inquiry specific key handlers
func (i *inquiryEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		if i.Text == "" && !state.alphasort.enabled {
			lists.ShowData.strlist = lists.GenListFromMap(state.currentList)
			lists.SelectEntry.list_loc = 0
			lists.ShowData.data.Reload()
			inquiry.LinkageMap = nil
			lists.List.Select(0)
			inquiryIndexAndExpand(0)
			inquiry.InqIntro.SetText("Type your regex query here,\nuse the enter key to filter your list:\n" + state.currentList + ", size: " + strconv.Itoa(len(lists.Data[state.currentList])))
		} else {
			lists.RegexSearch(i.Text)
		}
	case fyne.KeyDown: //for inquiry list navigation
		i.Entry.KeyDown(key)
		lists.List.Select(i.list_loc + 1)
		inquiry.InquiryScrollStop = true
		go inquiryScroll(*key, i.list_loc+1)
	case fyne.KeyUp: //for inquiry list navigation
		i.Entry.KeyUp(key)
		lists.List.Select(i.list_loc - 1)
		inquiry.InquiryScrollStop = true
		go inquiryScroll(*key, i.list_loc-1)
	case fyne.KeyLeft: //for inquiry list
		inquiry.InquiryTabs.SelectIndex(0)
	case fyne.KeyRight: //for inquiry detail
		inquiry.InquiryTabs.SelectIndex(1)
	case fyne.KeyPageDown:
		lists.List.ScrollToBottom()
	case fyne.KeyPageUp:
		lists.List.ScrollToTop()
	case fyne.KeyEscape: //for inquiry escape focused
		w.Close()
	}
}

func (m *inquiryEntry) TypedShortcut(s fyne.Shortcut) {
	if _, ok := s.(*desktop.CustomShortcut); !ok {
		m.Entry.TypedShortcut(s)
		return
	} else if ok {
		t := s.(*desktop.CustomShortcut)
		fmt.Println("shortcut name:", s.ShortcutName(), s.(*desktop.CustomShortcut).KeyName, s.(*desktop.CustomShortcut).Modifier)
		fmt.Println(desktop.SuperModifier)
		if t.Modifier == desktop.SuperModifier {
			switch t.KeyName {
			case fyne.KeyG:
				superAdd(s)
			case fyne.KeyE:
				superEdit(s)
			case fyne.KeyI:
				superInquire(s)
			case fyne.KeyR:
				superSwitchList(s)
			case fyne.KeyB:
				superClearInquiry(s)
			case fyne.KeyUp:
				superSwitchListUp(s)
			case fyne.KeyDown:
				superSwitchListDown(s)
			}
		}
	}
}

func (e *inquiryEntry) KeyUp(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyDown: //for inquiry stop scroll
		fallthrough
	case fyne.KeyUp:
		inquiry.InquiryScrollStop = false
	}
}

// these are global keyhandlers attatched to the desktop window
// they work in conjunction with the inquiry specific key handlers
func deskdown(key *fyne.KeyEvent) {
	if state.currentMenuItem == "Inquire" { //for inquiry
		switch key.Name {
		case fyne.KeyDown: //for inquiry list navigation
			lists.List.Select(lists.SelectEntry.list_loc + 1)
			inquiry.InquiryScrollStop = true
			go inquiryScroll(*key, lists.SelectEntry.list_loc)
		case fyne.KeyUp: //for inquiry list navigation
			lists.List.Select(lists.SelectEntry.list_loc - 1)
			inquiry.InquiryScrollStop = true
			go inquiryScroll(*key, lists.SelectEntry.list_loc)
		case fyne.KeyLeft: //for inquiry list
			inquiry.InquiryTabs.SelectIndex(0)
		case fyne.KeyRight: //for inquiry detail
			inquiry.InquiryTabs.SelectIndex(1)
		case fyne.KeyPageDown:
			lists.List.ScrollToBottom()
		case fyne.KeyPageUp:
			lists.List.ScrollToTop()
		case fyne.KeyEscape: //for inquiry escape (unfocused)
			w.Close()
		}
	} else if key.Name == fyne.KeyEscape { //for all views
		w.Close()
	}
}

func deskup(key *fyne.KeyEvent) {
	switch key.Name { //for inquiry
	case fyne.KeyLeft:
		fallthrough
	case fyne.KeyRight:
		fallthrough
	case fyne.KeyDown:
		fallthrough
	case fyne.KeyUp: //for inquiry list navigation
		inquiry.InquiryScrollStop = false
	}
}

// custom scrolling behavior applied to inquiry list
func inquiryScroll(key fyne.KeyEvent, loc int) {
	time.Sleep(500 * time.Millisecond)
	for inquiry.InquiryScrollStop {
		time.Sleep(50 * time.Millisecond)
		switch key.Name {
		case fyne.KeyDown: //for inquiry list navigation
			loc += 1
			if loc > len(lists.ShowData.strlist)-1 {
				inquiry.InquiryScrollStop = false
				break
			}
			lists.List.Select(loc)
		case fyne.KeyUp: //for inquiry list navigation
			loc -= 1
			if loc < 0 {
				inquiry.InquiryScrollStop = false
				break
			}
			lists.List.Select(loc)
		}
	}
}

func (l *userList) GenListFromMap(key string) []string {
	var res []string
	var searchstr bytes.Buffer
	inquiry.SearchMap = make(map[string]int)
	f, toremove := "", 0

	for idx, val := range l.Data[key] {
		res = append(res, val.Name)
		if idx == len(l.Data[key])-1 {
			f = fmt.Sprintf("%s %s %s", val.Name, strconv.Itoa(val.Rating), val.Tags)
			toremove = len(f)
		} else {
			f = fmt.Sprintf("%s %s %s\n", val.Name, strconv.Itoa(val.Rating), val.Tags)
			toremove = len(f) - 1
		}

		searchstr.WriteString(f)
		inquiry.SearchMap[f[:toremove]] = idx //generate regex search map, no linefeed
	}
	inquiry.FilterList = searchstr.String() //generate regex  search string
	return res
}

// generates list of list names in alphabetical order
func (l *userList) GetOrderedListNames() []string {
	keys := make([]string, len(l.Data))
	i := 0
	for k := range l.Data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func (l *userList) ListExists(listname string) bool {
	for k := range l.Data {
		if k == listname {
			return true
		}
	}
	return false
}

// regex search, and create linkage to original list datastructure
func (l *userList) RegexSearch(input string) {
	rep := regexp.MustCompile("(?im)^.*" + input + `.*$`)
	res := rep.FindAllString(inquiry.FilterList, -1)

	if state.alphasort.enabled { //order for alphabetical sort if enabled
		switch state.alphasort.order {
		case 0:
			sort.Strings(res)
		case 1:
			sort.Sort(sort.Reverse(sort.StringSlice(res)))
		}
	}

	rescnt, tmp := 0, []string{}
	tmplinkage := make(map[int]int)
	//generate linkage to orginal data mapping
	for idx, v := range res {
		if val, ok := inquiry.SearchMap[v]; ok {
			tmplinkage[idx] = val
			tmp = append(tmp, l.Data[state.currentList][val].Name)
			rescnt += 1
		} else {
			fmt.Println("v: ", v, "idx: ", idx)
			dialog.ShowError(fmt.Errorf("Results do not match data linkage,\nplease check your regular expression"), w)
			return
		}
	}

	inquiry.LinkageMap = tmplinkage
	if input == "" {
		inquiry.InqIntro.SetText("Type your regex query here,\nuse the enter key to filter your list:\n" + state.currentList + ", size: " + strconv.Itoa(len(l.Data[state.currentList])))
	} else {
		inquiry.InqIntro.SetText("Querying List: " + state.currentList + ", query: " + input + "\nresult size: " + strconv.Itoa(len(res)))
	}
	l.ShowData.strlist = tmp
	l.SelectEntry.list_loc = 0
	l.ShowData.data.Reload()
	l.List.Select(0) // ??why doesnt this call onselected??
	inquiryIndexAndExpand(0)
}

func (l *userList) removeElement(key string, index int) {
	l.Data[key] = append(l.Data[key][:index], l.Data[key][index+1:]...)
}

func (l *userList) RemoveElementByName(name string) bool {
	foundCount, idx := 0, 0
	for i := range l.Data[state.currentList] {
		if l.Data[state.currentList][i].Name == name {
			foundCount += 1
			idx = i
		}
	}

	if foundCount == 1 {
		l.removeElement(state.currentList, idx)
		return true
	}
	return false
}

// inquiry list item selection behavior
func inquiryIndexAndExpand(index int) {
	if index < 0 {
		index = 0
	} else if index > lists.List.Length()-1 {
		index = lists.List.Length() - 1
	}
	lists.SelectEntry.list_loc = index

	var item ListItem
	if inquiry.LinkageMap == nil {
		item = lists.Data[state.currentList][index]
	} else { //use the linkage
		item = lists.Data[state.currentList][inquiry.LinkageMap[index]]
	}

	inquiry.ExpandL1.SetText("Name: \n" + item.Name)
	inquiry.ExpandL2.SetText("Rating: \n" + strconv.Itoa(item.Rating))
	inquiry.ExpandL3.SetText("Tags: \n" + item.Tags)
	// reset focus since we favor keyboard navigation of the list
	w.Canvas().Focus(lists.SelectEntry)
}
