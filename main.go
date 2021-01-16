package main

import (
	"bufio"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/marcusolsson/tui-go/wordwrap"
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

var fname = "themovielist.csv" //file to store list in
var CurrentList []string       //current in memory list
var names []string             //names of lists in CurrentList

type searchFlags struct {
	alphabetical bool //sort alphabetically by name
	quickscroll  bool
}

//adds a movie to the CurrentList, to be saved on exit
//format: name, rating, tags, descrip... (must have name, and rating as number)
func addToList(input string, history *tui.Box, list int) {
	vars := strings.Split(input, ",")
	if input == "" || len(vars) < 2 || vars[0] == "" || vars[1] == "" {
		queryerr := tui.NewLabel("invalid input for add to list, type fullcmd for help")
		queryerr.SetStyleName("err")
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
	} else {
		ratCheck := regexp.MustCompile(`^[0-5]$`) //check for valid rating
		valid := ratCheck.MatchString(vars[1])
		if !valid {
			queryerr := tui.NewLabel("invalid ratinng, must contain rating [0,5]")
			queryerr.SetStyleName("err")
			history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
			return
		}
		//add to history command window
		queryres := tui.NewLabel(fmt.Sprintf(time.Now().Format("15:04")+" adding '%s' to the list '%s'", vars[0], names[list]))
		queryres.SetStyleName("res")
		history.Append(tui.NewHBox(
			tui.NewPadder(0, 0, queryres),
			tui.NewSpacer(),
		))

		//add to movie to list
		CurrentList[list] += strings.Join(vars, ",") + "\n"
	}
}

//filters movie by regular expression (or just movie name/tag), returns entire roow for match
func filterList(input string, history *tui.Box, historyBox *tui.Box, f searchFlags, list int) {
	matched, err := regexp.MatchString(input+`.*\n`, CurrentList[list])
	if !matched {
		queryres := tui.NewLabel(time.Now().Format("15:04") + " querying '" + input + "'\nNothing found during filter")
		queryres.SetStyleName("res")
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryres)))
	} else if err != nil {
		queryerr := tui.NewLabel("Error during filter: " + err.Error())
		queryerr.SetStyleName("err")
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
	} else {
		rep := regexp.MustCompile(".*" + input + `.*\n`)
		res := rep.FindAllString(CurrentList[list], -1)

		if f.alphabetical {
			sort.StringSlice(res).Sort()
		}

		queryres := tui.NewLabel(fmt.Sprintf(time.Now().Format("15:04")+" querying '%s' list '%s' result size=%d", input, names[list], len(res)))
		queryres.SetStyleName("res")
		//add to history command window
		history.Append(tui.NewHBox(
			tui.NewPadder(0, 0, queryres),
			tui.NewSpacer(),
		))

		history.Append(tui.NewHBox(tui.NewPadder(10, 0, tui.NewLabel(wordwrap.WrapString(fmt.Sprintf("\n%s\n", strings.Join(res, "\n")), historyBox.Size().X-20)))))
	}
}

//counts current list and appends result to the command window
func countList(history *tui.Box, list int) {
	rep := regexp.MustCompile(`.+\n`)
	lenlist := len(rep.FindAllString(CurrentList[list], -1))
	res := tui.NewLabel(fmt.Sprintf(time.Now().Format("15:04")+" Count for list '%s' result size=%d", names[list], lenlist))
	res.SetStyleName("res")
	//add to history command window
	history.Append(tui.NewHBox(
		tui.NewPadder(0, 1, res),
		tui.NewSpacer(),
	))
}

//removes a movie to the CurrentList, to be saved on exit
func removeFromList(input string, history *tui.Box, list int) {
	matched, err := regexp.MatchString(input+`.*\n`, CurrentList[list])
	if !matched {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Nothing found during search"))))
	} else if input == "" {
		queryerr := tui.NewLabel("invalid input for removing from list, type fullcmd for help")
		queryerr.SetStyleName("err")
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
	} else if err != nil {
		queryerr := tui.NewLabel("Error during search: " + err.Error())
		queryerr.SetStyleName("err")
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
	} else {
		queryres := tui.NewLabel(fmt.Sprintf(time.Now().Format("15:04")+" removing '%s' from the list '%s'", input, names[list]))
		queryres.SetStyleName("res")
		//add to history command window
		history.Append(tui.NewHBox(
			tui.NewPadder(0, 0, queryres),
			tui.NewSpacer(),
		))
		//remove
		rep := regexp.MustCompile(input + `.*\n`)
		CurrentList[list] = rep.ReplaceAllString(CurrentList[list], "")
	}
}

func switchList(input string, history *tui.Box, list int) int {
	var err error
	newind, err := strconv.Atoi(input)
	if err != nil {
		queryerr := tui.NewLabel("Invalid input for list to switch to")
		queryerr.SetStyleName("err")
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
	} else if newind >= 0 && newind < len(CurrentList)-1 {
		queryres := tui.NewLabel(fmt.Sprintf("Switched to list at index %d, '%s'", newind, names[newind]))
		queryres.SetStyleName("res")
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryres)))
		return newind
	} else {
		queryerr := tui.NewLabel("invalid switch index number")
		queryerr.SetStyleName("err")
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
	}
	return list
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

func addTag(input string, history *tui.Box, list int) {
	if strings.Contains(input, " newtag:") || input != "" {
		params := strings.Split(input, " newtag:")
		matched, err := regexp.MatchString(params[0]+`.*\n`, CurrentList[list])
		if !matched {
			history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Nothing found during search"))))
		} else if err != nil {
			queryerr := tui.NewLabel("Error during search: " + err.Error())
			queryerr.SetStyleName("err")
			history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
		} else if params[0] == "" {
			queryerr := tui.NewLabel("invalid input for adding tag, type fullcmd for help")
			queryerr.SetStyleName("err")
			history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
		} else {
			params[1] = strings.TrimSpace(params[1])
			rep := regexp.MustCompile(".*" + params[0] + `.*\n`)
			res := rep.FindAllString(CurrentList[list], -1)

			if len(res) == 1 {
				history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("oldtags:\n"+res[0][0:len(res[0])-1]))))
				queryres := tui.NewLabel("success")
				queryres.SetStyleName("res")
				history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryres)))
				history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("newtags:\n"+res[0][0:len(res[0])-1]+" "+params[1]))))

				CurrentList[list] = rep.ReplaceAllString(CurrentList[list], res[0][0:len(res[0])-1]+" "+params[1]+"\n")
			} else {
				queryerr := tui.NewLabel("too many search results, add more of the line you would like to add a tag to")
				queryerr.SetStyleName("err")
				history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
			}
		}
	} else {
		queryerr := tui.NewLabel("No new tag specified")
		queryerr.SetStyleName("err")
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, queryerr)))
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
			CurrentList[namescnt] += line + "\n"
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
	var chstore [50]string
	currentchloc, tmpchloc := 0, 0
	mode := "add" //default mode upon start is add

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
			if x == "add" || x == "search" || x == "remove" || x == "switch" || x == "createlist" || x == "addtag" {
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
				} else if x == "addtag" {
					history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("enter name of item to add tag"))))
				}
			} else if x == "fullcmd" {
				history.Append(tui.NewHBox(tui.NewLabel("command list:\n\nsearch: type regex to search in names/tags\nuse this format for multiple items:\nex. sci fi or comedy\n(sci fi|comedy)\n\nadd: enter name,rating,tag/tags,...\n\nremove: enter name of movie to remove, or regex matching any other field\n\nswitch: switch to another list by index\n\ncreatelist: create a list with provided name\n\naddtag: specify string to serach for to identify single item,\nthen add the identifier 'newtag:' followed by your new tag\n\nUse right arrow for quick scroll toggle\nUse Tab key to toggle alphabetical sort\n")))
			} else if x == "count()" {
				countList(history, currentList)
			} else {
				switch mode {
				case "add":
					addToList(x, history, currentList)
				case "addtag":
					addTag(x, history, currentList)
				case "search":
					filterList(x, history, historyBox, f, currentList)
				case "remove":
					removeFromList(x, history, currentList)
				case "switch":
					currentList = switchList(x, history, currentList)
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
		chstore[(currentchloc)%len(chstore)] = e.Text()
		currentchloc = currentchloc + 1
		tmpchloc = currentchloc
		input.SetText("")
	})

	//styles
	t := tui.NewTheme()
	t.SetStyle("label.res", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorBlue})
	t.SetStyle("label.err", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorRed})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetTheme(t)
	ui.SetKeybinding("Esc", func() { saveData(history); ui.Quit() })
	ui.SetKeybinding("Up", func() {
		if f.quickscroll == true {
			historyScroll.Scroll(0, -5)
		} else if tmpchloc-1 >= 0 {
			input.SetText(chstore[tmpchloc-1])
			tmpchloc = tmpchloc - 1
		}
	}) //both of these are for scroll mode
	ui.SetKeybinding("Down", func() {
		if f.quickscroll == true {
			historyScroll.Scroll(0, 5)
		} else if tmpchloc+1 < len(chstore) && tmpchloc+1 <= currentchloc%len(chstore)-1 {
			input.SetText(chstore[tmpchloc+1])
			tmpchloc = tmpchloc + 1
		} else if tmpchloc == currentchloc%len(chstore)-1 {
			tmpchloc = currentchloc % len(chstore)
			input.SetText("")
		}
	})
	ui.SetKeybinding("Right", func() { historyScroll.SetAutoscrollToBottom(f.quickscroll); f.quickscroll = !f.quickscroll }) //for quick scroll hotkey
	ui.SetKeybinding("Tab", func() {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel(fmt.Sprintf("\n%s,%t\n", "toggling alphabetical sort", !f.alphabetical)))))
		f.alphabetical = !f.alphabetical
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
