package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

//misc structures
type MenuPageLink struct {
	View func(w fyne.Window) fyne.CanvasObject
}

//form structures
type submitEntry struct {
	widget.Entry
}

//inquiry structures
type inquiryEntry struct {
	widget.Entry
	list_loc int //move this to the listdata struct??
}

type listData struct {
	data    binding.ExternalStringList
	strlist []string
}

type ListItem struct {
	Name   string
	Rating int
	Tags   string
}

type userList struct {
	Data         map[string][]ListItem
	List         *widget.List
	SelectEntry  *inquiryEntry
	ShowData     listData
	ListModified bool
}

type Inquiry struct {
	FilterList        string
	SearchMap         map[string]int
	LinkageMap        map[int]int
	ExpandL1          *widget.Label
	ExpandL2          *widget.Label
	ExpandL3          *widget.Label
	InquiryTabs       *container.AppTabs
	InquiryScrollStop bool
	InqTitle          *widget.Label
	InqIntro          *widget.Label
}

//application state structures
type AppState struct {
	currentList     string
	currentMenuItem string
	noList          bool
	alphasort       AlphaSort
}

type AlphaSort struct {
	enabled bool
	order   int //0 asc, 1 desc
}
