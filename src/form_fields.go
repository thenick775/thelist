package main

import (
	"fyne.io/fyne/v2"
)

//used to associate a current form's last field to its enter key action
var currFormFunc func()

//entry used to submit with enter key
func NewSubmitEntry() *submitEntry {
	entry := &submitEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

//entry func used to submit the form on enter key
func (s *submitEntry) KeyUp(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeyReturn:
		currFormFunc()
		w.Canvas().Unfocus()
	}
}
