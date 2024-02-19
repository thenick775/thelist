package main

import (
	"image"
	"image/color"
	"strings"

	"github.com/psykhi/wordclouds" //=> replaced with github.com/thenick775/wordclouds
)

var lightModeColors = []color.RGBA{
	{0x1b, 0x1b, 0x1b, 0xff},
	{0x48, 0x48, 0x4B, 0xff},
	{0x3a, 0x1b, 0xd1, 0xff},
	{0x10, 0x74, 0xe6, 0xff},
	{0xc7, 0x26, 0x41, 0xff},
	{0x70, 0xD6, 0xBF, 0xff},
}
var lightModeBackground = color.RGBA{255, 255, 255, 255}

var darkModeColors = []color.RGBA{
	{0x99, 0x6d, 0xb5, 0xff},
	{0xab, 0x1d, 0x13, 0xff},
	{0x96, 0x82, 0xfa, 0xff},
	{0x65, 0xCD, 0xFA, 0xff},
	{0xff, 0xff, 0xff, 0xff},
	{0x70, 0xD6, 0xBF, 0xff},
}
var darkModeBackground = color.RGBA{48, 48, 48, 255}

var defaultConfImg = ConfImg{
	FontMaxSize:     600,
	FontMinSize:     15,
	RandomPlacement: false,
	FontFile:        "/Roboto-Regular.ttf", //prepended with fontLoc once initialized
	Colors:          darkModeColors,        //dark is default
	BackgroundColor: darkModeBackground,    //dark is default
	Width:           2048,
	Height:          2048,
	Mask: MaskConf{"", color.RGBA{ //no masking by default
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	}},
}

// generate image object containing word cloud
func genWordCloudImg() (image.Image, map[string]int) {
	confImg := defaultConfImg
	confImg.FontFile = fontLoc + confImg.FontFile
	//exclusion zones if present
	var boxes []*wordclouds.Box
	if confImg.Mask.File != "" {
		boxes = wordclouds.Mask(
			confImg.Mask.File,
			confImg.Width,
			confImg.Height,
			confImg.Mask.Color)
	}
	//word colors
	if strings.EqualFold(state.currentThemeAlias, "light") {
		confImg.Colors = lightModeColors
		confImg.BackgroundColor = lightModeBackground
	}

	colors := make([]color.Color, 0)
	for _, c := range confImg.Colors {
		colors = append(colors, c)
	}
	//data processing
	wordCounts := make(map[string]int)
	for _, item := range lists.Data[state.currentList] {
		splitTags := strings.Fields(item.Tags)

		for _, tag := range splitTags {
			wordCounts[tag] += 1
		}
	}

	cloud := wordclouds.NewWordcloud(wordCounts,
		wordclouds.FontFile(confImg.FontFile),
		wordclouds.FontMaxSize(confImg.FontMaxSize),
		wordclouds.FontMinSize(confImg.FontMinSize),
		wordclouds.Colors(colors),
		wordclouds.BackgroundColor(confImg.BackgroundColor),
		wordclouds.MaskBoxes(boxes),
		wordclouds.Height(confImg.Height),
		wordclouds.Width(confImg.Width),
		wordclouds.RandomPlacement(confImg.RandomPlacement))
	//image generation
	img := cloud.Draw()
	return img, wordCounts
}

func genStats() []ListsSummary {
	var ret []ListsSummary
	keys := lists.GetOrderedListNames()
	for _, key := range keys {
		listsum := ListsSummary{
			Name:                  key,                             //name of list
			TotalContentCount:     len(lists.Data[key]),            //total number of items in list
			ContentCountPerRating: getRatingCount(lists.Data[key]), //count of items in list per rating
		}
		ret = append(ret, listsum)
	}

	return ret
}

func getRatingCount(list []ListItem) map[int]int {
	var ret = make(map[int]int)
	for _, listitem := range list {
		ret[listitem.Rating] += 1
	}
	return ret
}
