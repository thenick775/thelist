package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"

	"fyne.io/fyne/v2/dialog"
)

func write_conf() {
	conf_rewrite, err := json.MarshalIndent(conf, "", " ")
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to marshal list:\n"+err.Error()), w)
	}
	err = os.WriteFile(confLoc, conf_rewrite, 0644)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to save configuration"), w)
	}
}

func write_csv(fullexport bool, fname string) {
	var buf string

	if fullexport {
		buf += "name,rating,tags\n"
		if inquiry.LinkageMap == nil {
			for _, val := range lists.Data[state.currentList] {
				buf += fmt.Sprintf("%s,%s,%s\n", val.Name, strconv.Itoa(val.Rating), val.Tags)
			}
		} else {
			for k := range inquiry.LinkageMap {
				val := lists.Data[state.currentList][inquiry.LinkageMap[k]]
				buf += fmt.Sprintf("%s,%s,%s\n", val.Name, strconv.Itoa(val.Rating), val.Tags)
			}
		}
	} else {
		buf += "name\n"
		for _, v := range lists.ShowData.strlist {
			buf += v + "\n"
		}
	}

	err := os.WriteFile(fname, []byte(buf), 0644)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to create export file:\n"+err.Error()), w)
	}
}

func write_json(fullexport bool, fname string) {
	var (
		buf []byte
		err error
	)
	if fullexport {
		if inquiry.LinkageMap == nil {
			buf, err = json.MarshalIndent(lists.Data[state.currentList], "", " ")
			if err != nil {
				dialog.ShowError(fmt.Errorf("Failed to marshal list:\n"+err.Error()), w)
			}
		} else {
			var items []ListItem
			for k := range inquiry.LinkageMap {
				item := lists.Data[state.currentList][inquiry.LinkageMap[k]]
				items = append(items, item)
			}
			buf, err = json.MarshalIndent(items, "", " ")
			if err != nil {
				dialog.ShowError(fmt.Errorf("Failed to marshal list:\n"+err.Error()), w)
			}
		}
	} else {
		buf, err = json.MarshalIndent(lists.ShowData.strlist, "", " ")
		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to create export file:\n"+err.Error()), w)
		}
	}

	err = os.WriteFile(fname, buf, 0644)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to create export file:\n"+err.Error()), w)
	}
}

func write_png(img image.Image, fname string) {
	f, err := os.Create(fname)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to open save file:\n"+err.Error()), w)
	} else {
		defer f.Close()
		if err = png.Encode(f, img); err != nil {
			dialog.ShowError(fmt.Errorf("Failed to encode file:\n"+err.Error()), w)
		}
	}
}
