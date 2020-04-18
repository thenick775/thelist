package main

import (
	"bufio"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

//This program was made just because I watch and read an immense amount of movies/books etc,
//and could not keep track of them all. I also wanted the ability to filter a list
//by regular expression, and be able to add descriptive tags or other fields as needed.
//This was made quick and dirty, as I just wanted an easy way for me to store all
//this information. Any recommendations are welcome.

var mode = "add"               //default mode upon start is add
var fname = "themovielist.csv" //file to store list in
var CurrentList []string       //current in memory list
var names []string

type searchFlags struct {
	alphabetical bool //sort alphabetically by name
	quickscroll  bool
}

//adds a movie to the CurrentList, to be saved on exit
//format: name, rating, tags, descrip... (must have name, and rating as number)
func addToList(input string, history *tui.Box, list int) {
	vars := strings.Split(input, ",")
	if input == "" || len(vars) < 2 || vars[0] == "" || vars[1] == "" {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("invalid input for add to list, type fullcmd for help"))))
	} else {
		ratCheck := regexp.MustCompile(`^[0-5]$`) //check for valid rating
		valid := ratCheck.MatchString(vars[1])
		if !valid {
			history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("invalid ratinng, must contain rating [0,5]"))))
			return
		}
		//add to history command window
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("adding '%s' to the list '%s'", vars[0], names[list]))),
			tui.NewSpacer(),
		))

		//add to movie to list
		CurrentList[list] += strings.Join(vars, ",") + "\n"
	}
}

//filters movie by regular expression (or just movie name/tag), returns entire roow for match
func filterList(input string, history *tui.Box, f searchFlags, list int) {
	matched, err := regexp.MatchString(input+`.*\n`, CurrentList[list])
	if !matched {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel(time.Now().Format("15:04")+" querying '"+input+"'\nNothing found during filter"))))
	} else if err != nil {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Errr during filter: "+err.Error()))))
	} else {
		rep := regexp.MustCompile(".*" + input + `.*\n`)
		res := rep.FindAllString(CurrentList[list], -1)

		if f.alphabetical {
			sort.StringSlice(res).Sort()
		}

		//add to history command window
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("querying '%s' list '%s' result size=%d", input, names[list], len(res)))),
			tui.NewSpacer(),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("\n%s\n", strings.Join(res, "\n")))),
		))
	}
}

//removes a movie to the CurrentList, to be saved on exit
func removeFromList(input string, history *tui.Box, list int) {
	matched, err := regexp.MatchString(input+`.*\n`, CurrentList[list])
	if !matched {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Nothing found during search"))))
	} else if input == "" {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("invalid input for removing from list, type fullcmd for help"))))
	} else if err != nil {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Errr during search: "+err.Error()))))
	} else {
		//add to history command window
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("removing '%s' from the list '%s'", input, names[list]))),
			tui.NewSpacer(),
		))
		//remove
		rep := regexp.MustCompile(input + `.*\n`)
		CurrentList[list] = rep.ReplaceAllString(CurrentList[list], "")
	}
}

func switchList(input string, history *tui.Box,list int) int{
	var err error
	newind, err := strconv.Atoi(input)
	if err != nil {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Invalid input for list to switch to"))))
		return list
	}
	if newind >= 0 && newind < len(CurrentList)-1 {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel(fmt.Sprintf("Switched to list at index %d", newind)))))
		return newind
	} else {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("invalid switch index number"))))
		return list
	}
}

//saves data on program exit
func saveData(history *tui.Box) {
	history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("saving data back to disk"))))
	f, err := os.Create(fname)

	for i, name := range names {
		_, err = f.WriteString(CurrentList[i])
		if err != nil {
			f.Close()
			log.Fatal(err)
		}
		_, err = f.WriteString("#" + name + "\n")
		if err != nil {
			f.Close()
			log.Fatal(err)
		}
	}

	err = f.Close()
	if err != nil {
		f.Close()
		log.Fatal(err)
	}
}

//reads existing data upon program start
func readData() {
	file, err := os.Open(fname)
	if err != nil {
		file.Close()
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	namescnt := 0
	CurrentList = append(CurrentList, "")
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '#' {
			names = append(names, line[1:])
			CurrentList = append(CurrentList, "")
			namescnt += 1
		} else {
			CurrentList[namescnt] += scanner.Text() + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		file.Close()
		log.Fatal(err)
	}
}

func main() {
	f := searchFlags{
		alphabetical: false,
		quickscroll:  false,
	}
	//setup initial side labels
	sidebar := tui.NewVBox(
		tui.NewLabel("Type:\nsearch\nadd\nremove\nswitch\nor createlist\nto switch mode\n\nUse right arrow\nto toggle scroll"),
		tui.NewLabel(""),
		tui.NewLabel("In Search\nspecify:\nregex to find\nname/tag/tags...\nuse semicolon to\ndelimit multiple\nstatements"),
		tui.NewLabel(""),
		tui.NewLabel("Type:\nfullcmd\nto list all\ncommands"),
		tui.NewLabel(""),
		tui.NewLabel("Use Esc\nto quit"),
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()
	history.Append(tui.NewHBox(
		tui.NewLabel(time.Now().Format("15:04")),
		tui.NewLabel("hello, themovielist has started"),
	))
	history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("currently in mode: "+mode))))

	//create data file if it does not exist
	if _, err := os.Stat(fname); err != nil {
		if os.IsNotExist(err) {
			history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("datafile does not exist, creating fresh one"))))
			os.Create(fname)
		}
	} else {
		readData()
	}

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	currentList := 0
	input.OnSubmit(func(e *tui.Entry) {
		inputtxt := strings.Split(e.Text(), ";")
		for _, x := range inputtxt {
			if x == "add" || x == "search" || x == "remove" || x == "switch" || x == "createlist" {
				mode = x
				history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("switching to mode: "+mode))))

				if x == "switch" { //need to write scheme to add list, remove list
					listnames := "list indices and names:\n"
					for i, val := range names {
						listnames += fmt.Sprintf("%d %s\n\n", i, val)
					}

					history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel(listnames+"enter index of list to switch to: "))))
				} else if x == "createlist" {
					history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("enter name of new list"))))
				}
			} else if x == "fullcmd" {
				history.Append(tui.NewHBox(tui.NewLabel("command list:\n\nsearch: type regex to search in names/tags\nuse this format for multiple items:\nex. sci fi or comedy\n(sci fi|comedy)\n\nadd: enter name,rating,tag/tags,...\n\nremove: enter name of movie to remove, or regex matching any other field\n\nswitch: switch to another list by index\n\ncreatelist: create a list with provided name\n\nUse right arrow for quick scroll toggle\nUse Tab key to toggle alphabetical sort\n")))
			} else {
				switch mode {
				case "add":
					addToList(x, history, currentList)
				case "search":
					filterList(x, history, f, currentList)
				case "remove":
					removeFromList(x, history, currentList)
				case "switch":
					currentList=switchList(x,history,currentList)
				case "createlist": //maybe add removelist in future
					if x != "" {
						CurrentList = append(CurrentList, "")
						names = append(names, x)
					} else {
						history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("invalid input for new list name, type fullcmd for help"))))
					}
				}
			}
		}

		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { saveData(history); ui.Quit() })
	ui.SetKeybinding("Up", func() { historyScroll.Scroll(0, -5) }) //both of these are for scroll mode
	ui.SetKeybinding("Down", func() { historyScroll.Scroll(0, 5) })
	ui.SetKeybinding("Right", func() { historyScroll.SetAutoscrollToBottom(f.quickscroll); f.quickscroll = !f.quickscroll }) //for quick scroll hotkey
	ui.SetKeybinding("Tab", func() {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel(fmt.Sprintf("\n%s,%t\n", "toggling alphabetical sort", !f.alphabetical)))))
		f.alphabetical = !f.alphabetical
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
