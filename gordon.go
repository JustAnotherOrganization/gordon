package main

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"

	"github.com/JustAnotherOrganization/gordon/jim"
	"github.com/JustAnotherOrganization/gordon/menu"
	"github.com/JustAnotherOrganization/gordon/notebook"
)

var (
	mainWindow *gtk.Window
	nb         *notebook.Notebook
)

func setupWindow() (*gtk.Grid, error) {
	var err error
	nb, err = notebook.New()
	if err != nil {
		return nil, errors.Wrap(err, "notebook.New")
	}

	menuBar, err := menu.SetupMenu(mainWindow, nb)
	if err != nil {
		return nil, errors.Wrap(err, "menu.SetupMenu")
	}

	grid, err := gtk.GridNew()
	if err != nil {
		return nil, errors.Wrap(err, "gtk.GridNew")
	}

	grid.Attach(menuBar, 0, 0, 1, 1)
	grid.AttachNextTo(nb.Book, menuBar, gtk.POS_BOTTOM, 1, 1)
	grid.ShowAll()

	return grid, nil
}

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	var err error
	mainWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(err)
	}
	mainWindow.SetTitle("Gordon")
	mainWindow.SetDefaultSize(800, 600)
	mainWindow.Connect("destroy", func() {
		if nb != nil {
			nb.Book.Destroy()
		}
		gtk.MainQuit()
	})

	// Create a grid window with our notebook structure built in.
	gWindow, err := setupWindow()
	if err != nil {
		log.Fatal(err)
	}

	mainWindow.Add(gWindow)

	// Recursively show all widgets contained in this window.
	mainWindow.ShowAll()

	// Create a tab for handling Tiberious server messages (conversations will
	// get individual tabs).
	tiberiousWindow, err := nb.OpenTab("Tiberious")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := jim.ConnectSocket("ws://localhost:4002/ws")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func(conn *websocket.Conn) {
		for {
			byt, err := jim.ReadSocketMessage(conn)
			if err != nil {
				// TODO properly format this and display it as a connection error
				tiberiousWindow.InsertAtCursor(err.Error())
				log.Print(err)
				break
			}

			tiberiousWindow.InsertAtCursor(string(byt) + "\n")
		}
	}(conn)

	// Begin executing the GTK main loop. This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
