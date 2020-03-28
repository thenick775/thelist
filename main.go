package main

import (
	"bufio"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

//This program was made just because I watch an immense amount of movies,
//and could not keep track of them all. I also wanted the ability to filter
//by regular expression, and be able to tag movies or add descriptions or other fields as needed.
//This was made quick and dirty, as I just wanted an easy way for me to store all
//this information. Any recommendations are welcome. Eventually I plan to write the
//file in binary, or to a compressed format.

var mode = "add"               //default mode upon start is add
var fname = "themovielist.csv" //file to store movielist in, in the future this may not be human readable
var CurrentList string         //current in memory movie list

//adds a movie to the CurrentList, to be saved on exit
//format: name, rating, tags, descrip... (must have name, and rating as number)
func addMovie(input string, history *tui.Box) {
	vars := strings.Split(input, ",")
	if input == "" || len(vars) < 2 || vars[0] == "" || vars[1] == "" {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("invalid input for adding movie, type fullcmd for help"))))
	} else {
		ratCheck := regexp.MustCompile(`^[0-5]$`) //check for valid rating
		valid := ratCheck.MatchString(vars[1])
		if !valid {
			history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("invalid ratinng, movie must contain rating [0,5]"))))
			return
		}
		//add to history command window
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("adding movie '%s' to the list", vars[0]))),
			tui.NewSpacer(),
		))

		//add to movie to list
		CurrentList += strings.Join(vars, ",") + "\n"
	}
}

//filters movie by regular expression (or just movie name/tag), returns entire roow for match
func filterMovie(input string, history *tui.Box) {
	matched, err := regexp.MatchString(input+`.*\n`, CurrentList)
	if !matched {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Nothing found during filter"))))
	} else if err != nil {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Errr during filter: "+err.Error()))))
	} else {
		rep := regexp.MustCompile(".*" + input + `.*\n`)
		res := rep.FindAllString(CurrentList, -1)
		//add to history command window
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("querying '%s' result size=%d", input, len(res)))),
			tui.NewSpacer(),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("\n%s\n", strings.Join(res,"\n")))),
		))
	}
}

//removes a movie to the CurrentList, to be saved on exit
func removeMovie(input string, history *tui.Box) {
	matched, err := regexp.MatchString(input+`.*\n`, CurrentList)
	if !matched {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Nothing found during search"))))
	} else if input == "" {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("invalid input for removing movie, type fullcmd for help"))))
	} else if err != nil {
		history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("Errr during search: "+err.Error()))))
	} else {
		//add to history command window
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("removing '%s' from the list", input))),
			tui.NewSpacer(),
		))
		//remove
		rep := regexp.MustCompile(input + `.*\n`)
		CurrentList = rep.ReplaceAllString(CurrentList, "")
	}
}

//saves data on program exit
func saveData(history *tui.Box) {
	history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("saving data back to disk"))))
	f, err := os.Create(fname)

	_, err = f.WriteString(CurrentList)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

//reads existing data upon program start
func readData() {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		CurrentList += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	//setup initial side labels
	sidebar := tui.NewVBox(
		tui.NewLabel("Type:\nsearch\nadd\nremove\nor scroll\nto switch mode"),
		tui.NewLabel(""),
		tui.NewLabel("In Search\nspecify:\nregex to find\nname/tag/tags..."),
		//tui.NewLabel("Use cmd:\nswitch list->...\nto switch lists"),
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

	input.OnSubmit(func(e *tui.Entry) {
		if e.Text() == "add" || e.Text() == "search" || e.Text() == "remove" || e.Text() == "scroll" {
			if mode == "scroll" {
				historyScroll.SetAutoscrollToBottom(true)
			} else if e.Text() == "scroll" {
				historyScroll.SetAutoscrollToBottom(false)
			}
			mode = e.Text()
			history.Append(tui.NewHBox(tui.NewPadder(0, 0, tui.NewLabel("switching to mode: "+mode))))
		} else if e.Text() == "fullcmd" {
			history.Append(tui.NewHBox(tui.NewLabel("command list:\n\nsearch: type regex to search in names/tags\n\nadd: enter name,rating,tag/tags,...\n\nremove: enter name of movie to remove, or regex matching any other field\n\nscroll: type scroll, then use arrow keys to movie view up or down")))
		} else {
			switch mode {
			case "add":
				addMovie(e.Text(), history)
			case "search":
				filterMovie(e.Text(), history)
			case "remove":
				removeMovie(e.Text(), history)
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

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
