package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	state   AppState
	lists   userList
	conf    map[string]interface{}
	inquiry Inquiry
	w       fyne.Window
	a       fyne.App
	tree    *widget.Tree
	confLoc = filepath.FromSlash("/conf.json") //required conf location, appended to executable location
)

func main() {
	path, err := os.Executable() //get path of current executable
	if err != nil {
		panic(err)
	}
	fmt.Println("execu path: ", path)
	execupath := filepath.Dir(path)
	confLoc = execupath + confLoc //append configuration location to executable path (in same dir)
	fmt.Println("conf loc:", confLoc)

	//get configuration
	conf = make(map[string]interface{})
	//read configuration file
	conf_file, err := ioutil.ReadFile(confLoc)
	if err != nil {
		panic(err)
	}

	// Decode config json into our map
	err = json.Unmarshal(conf_file, &conf)
	if err != nil {
		panic("config err:" + err.Error())
		return
	}

	defaultSelected := conf["configuration"].(map[string]interface{})["default selected"].(string)
	defaultBranch := conf["configuration"].(map[string]interface{})["default branches open"].([]interface{})[0].(string)
	defaultTheme := conf["configuration"].(map[string]interface{})["default theme"].(string)
	local_item_file := conf["configuration"].(map[string]interface{})["local item file"].(string)
	state.currentList = conf["configuration"].(map[string]interface{})["default list"].(string)
	state.noList = false
	state.alphasort.enabled = false
	state.alphasort.order = 1

	a = app.NewWithID("com.vancise.thelist")
	a.SetIcon(theme.FyneLogo())
	w = a.NewWindow("TheList Utility")
	//setup app tree menu
	setupSystemMenu(w, a)
	//decode list data
	if local_item_file != "" { //how to handle this situation for users??
		byteValue, err := ioutil.ReadFile(local_item_file)
		if err != nil {
			fmt.Println("local item error")
			panic(err)
		}
		err = json.Unmarshal(byteValue, &lists.Data)
		if err != nil {
			panic(err)
		}
	} else {
		state.noList = true
		dialog.ShowInformation("Information", "No list file,\nPlease select a file location using:\nConfiguration > Defaults", w)
		lists.Data = make(map[string][]ListItem)
	}

	//intialize lists and inquiry
	inquiry.Initialize()
	lists.Initialize()

	mainView := container.NewMax() //placeholder that will take up max size of panel
	tree = menuTree(w, mainView, defaultBranch, defaultSelected)

	//set theme
	if strings.EqualFold(defaultTheme, "light") {
		a.Settings().SetTheme(theme.LightTheme())
	} else if strings.EqualFold(defaultTheme, "dark") {
		a.Settings().SetTheme(theme.DarkTheme())
	}

	if fyne.CurrentDevice().IsMobile() {
		panic("mobile not yet supported")
	} else {
		split := container.NewHSplit(container.NewBorder(container.NewVBox(
			widget.NewLabel("Main Menu"), widget.NewSeparator()), nil, nil, nil, tree), mainView)
		split.Offset = 0.1
		w.SetContent(split)
	}

	if deskCanvas, ok := w.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(deskdown) //for monitoring navigation of the list in inquire mode
		deskCanvas.SetOnKeyUp(deskup)
	} else {
		panic("mobile not yet supported")
	}

	w.Resize(fyne.NewSize(940, 660))
	w.Canvas().Focus(lists.SelectEntry)
	w.SetOnClosed(func() {
		if lists.ListModified {
			fmt.Println("list modified, saving to file")
			file2, err := json.MarshalIndent(lists.Data, "", " ")
			if err != nil {
				panic(err)
			}
			local_item_file := conf["configuration"].(map[string]interface{})["local item file"].(string)
			err = ioutil.WriteFile(local_item_file, file2, 0644)
			if err != nil {
				panic(err)
			}
		}
	})
	//desktop shortcuts
	setupDesktopShortcuts(w)

	lists.List.Select(0)
	w.ShowAndRun()
}
