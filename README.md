# TheList
An app to make a fast, searchable, in memory set of lists.

This is a generic listing utility capable of managing multiple lists. Fyne was used for the application display and packaging.

I built this because sadly, I watch so many movies and read so many books that it became impossible for me to keep track of my favorites in my head.

The search utilizes regular expressions, which is the primary way I wanted to be able to filter and search by names/multiple item tags.

This is the currnt version of this software, the terminal TUI utility was retired in favor of the graphical display, where keyboard shortcuts for this version were included for ease of use.

This application has been tested and packaged on OSX, where the script included places the configuration file next to the executable.

#### Features
- Search list items by regular expression
- Support for multiple lists
- Sort alphabetically
- Add/Remove/Edit/Delete list items
- Add/Remove/Edit/Delete lists

#### Application File System Structure
TheList.app
└── Contents
    ├── Info.plist
    ├── MacOS
    │   ├── conf.json
    │   └── the-list
    └── Resources
        └── icon.icns

#### Demo

![](demo_fyne_v1.gif)
