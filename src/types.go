package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// misc structures
type MenuPageLink struct {
	View func(w fyne.Window) fyne.CanvasObject
}

// form structures
type submitEntry struct {
	widget.Entry
	currFormFunc func()
}

// inquiry structures
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

type ListsSummary struct {
	Name                  string
	TotalContentCount     int
	ContentCountPerRating map[int]int
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

// application state structures
type AppState struct {
	currentList       string
	currentMenuItem   string
	noList            bool
	alphasort         AlphaSort
	currentThemeAlias string
}

type AlphaSort struct {
	enabled bool
	order   int //0 asc, 1 desc
}

// Default Word Cloud Configuration
type ConfImg struct {
	FontMaxSize     int          `json:"font_max_size"`
	FontMinSize     int          `json:"font_min_size"`
	RandomPlacement bool         `json:"random_placement"`
	FontFile        string       `json:"font_file"`
	Colors          []color.RGBA `json:"colors"`
	BackgroundColor color.RGBA   `yaml:"background_color"`
	Width           int          `json:"width"`
	Height          int          `json:"height"`
	Mask            MaskConf     `json:"mask"`
}

type MaskConf struct {
	File  string     `json:"file"`
	Color color.RGBA `json:"color"`
}
