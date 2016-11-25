package notebook

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

// Notebook extends access to gtk.Notebook.
type Notebook struct {
	Book *gtk.Notebook
}

// New creates a new Notebook object.
func New() (*Notebook, error) {
	notebook := &Notebook{}
	book, err := gtk.NotebookNew()
	if err != nil {
		return nil, errors.Wrap(err, "gtk.NotebookNew")
	}

	notebook.Book = book
	notebook.Book.Connect("destroy", func() {
		notebook.CloseAll()
	})

	return notebook, nil
}

// OpenTab opens a new notebook tab with a ScrolledWindow, TextBuffer and a
// label attaches it to our notebook and returns the TextBuffer so it can be
// written to.
func (nb *Notebook) OpenTab(label string) (*gtk.TextBuffer, error) {
	// Create a widget and set it up to be a close button using the stock gtk icons.
	btn, err := gtk.ButtonNew()
	if err != nil {
		return nil, errors.Wrap(err, "gtk.ButtonNew")
	}

	btn.SetFocusOnClick(false)
	btn.SetName("close-tab-button")
	image, err := gtk.ImageNewFromIconName("window-close", gtk.ICON_SIZE_MENU)
	if err != nil {
		return nil, errors.Wrap(err, "gtk.ImageNewFromIconName")
	}
	btn.Add(image)

	// Also create a widget to hold the text label for our new tab.
	tab, err := newTab(label)
	if err != nil {
		return nil, errors.Wrap(err, "NewTab")
	}

	// Generate a grid and place both our tab label and close button into it.
	tabGrid, err := gtk.GridNew()
	if err != nil {
		return nil, errors.Wrap(err, "gtk.GridNew")
	}
	tabGrid.Attach(tab.Label, 0, 0, 1, 1)
	tabGrid.AttachNextTo(btn, tab.Label, gtk.POS_RIGHT, 1, 1)
	tabGrid.ShowAll()

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "gtk.ScrolledWindowNew")
	}
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.SetBorderWidth(5)
	// Make sure the scrolledWindow is set to fill the entire non-occupied area
	// of the window.
	scrolledWindow.SetHExpand(true)
	scrolledWindow.SetHAlign(gtk.ALIGN_FILL)
	scrolledWindow.SetVExpand(true)
	scrolledWindow.SetVAlign(gtk.ALIGN_FILL)

	textView, err := gtk.TextViewNew()
	if err != nil {
		return nil, errors.Wrap(err, "gtk.TextViewNew")
	}
	textView.SetEditable(false)
	textView.SetWrapMode(gtk.WRAP_WORD)
	textView.SetHExpand(true)
	textView.SetHAlign(gtk.ALIGN_FILL)
	textView.SetVExpand(true)
	textView.SetVAlign(gtk.ALIGN_FILL)
	textView.SetCursorVisible(false)
	scrolledWindow.Add(textView)

	buffer, err := textView.GetBuffer()
	if err != nil {
		return nil, errors.Wrap(err, "textView.GetBuffer")
	}

	scrolledWindow.ShowAll()

	// Now that we have both the tab and the window created go ahead and connect
	// the close button from the tab to the window.
	btn.Connect("clicked", func() {
		nb.CloseTab()
	})

	nb.Book.AppendPage(scrolledWindow, tabGrid)
	nb.Book.SetTabReorderable(scrolledWindow, true)

	return buffer, nil
}

func (nb *Notebook) closeTab(tab int) {
	nb.Book.RemovePage(tab)
}

// CloseTab closes the selected tab.
func (nb *Notebook) CloseTab() {
	// TODO make this work for tabs that aren't the "current" page...
	// Currently the tab with focus is closed...
	nb.closeTab(nb.Book.GetCurrentPage())
}

// CloseAll closes all open tabs.
func (nb *Notebook) CloseAll() {
	for i := nb.Book.GetNPages(); i > 0; i-- {
		nb.closeTab(i - 1)
	}
}
