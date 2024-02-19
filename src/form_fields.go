package main

import (
	"fyne.io/fyne/v2"
)

// entry used to submit with enter key
func NewSubmitEntry() *submitEntry {
	entry := &submitEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

// entry func used to submit the form on enter key
func (s *submitEntry) KeyUp(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeyReturn:
		s.currFormFunc()
		w.Canvas().Unfocus()
	}
}
