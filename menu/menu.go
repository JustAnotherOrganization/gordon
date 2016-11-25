package menu

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"

	"github.com/JustAnotherOrganization/gordon/notebook"
)

// SetupMenu returns a menubar with all our menu options attached.
func SetupMenu(mainWindow *gtk.Window, nb *notebook.Notebook) (*gtk.MenuBar, error) {
	menuBar, err := gtk.MenuBarNew()
	if err != nil {
		return nil, errors.Wrap(err, "gtk.MenuBarNew")
	}

	fileMenu, err := setupFileMenu(mainWindow, nb)
	if err != nil {
		return nil, errors.Wrap(err, "setupFileMenu")
	}
	menuBar.Add(fileMenu)

	helpMenu, err := setupHelpMenu(mainWindow)
	if err != nil {
		return nil, errors.Wrap(err, "setupHelpMenu")
	}
	menuBar.Add(helpMenu)

	return menuBar, nil
}

func setupFileMenu(mainWindow *gtk.Window, nb *notebook.Notebook) (*gtk.MenuItem, error) {
	itemFile, err := gtk.MenuItemNewWithLabel("File")
	if err != nil {
		return nil, errors.Wrap(err, "gtk.MenuItemNewWithLabel")
	}
	fileMenu, err := gtk.MenuNew()
	if err != nil {
		return nil, errors.Wrap(err, "gtk.MenuNew")
	}
	itemFile.SetSubmenu(fileMenu)
	itemCloseAll, err := gtk.MenuItemNewWithLabel("Close All")
	if err != nil {
		return nil, errors.Wrap(err, "gtk.MenuItemNewWithLabel")
	}
	fileMenu.Add(itemCloseAll)
	itemCloseAll.Connect("activate", func() {
		nb.CloseAll()
	})

	itemQuit, err := gtk.MenuItemNewWithLabel("Quit")
	if err != nil {
		return nil, errors.Wrap(err, "gtk.MenuItemNewWithLabel")
	}
	fileMenu.Add(itemQuit)
	itemQuit.Connect("activate", func() {
		mainWindow.Destroy()
	})

	return itemFile, nil
}

func setupHelpMenu(mainWindow *gtk.Window) (*gtk.MenuItem, error) {
	itemHelp, err := gtk.MenuItemNewWithLabel("Help")
	if err != nil {
		return nil, errors.Wrap(err, "gtk.MenuItemNewWithLabel")
	}
	helpMenu, err := gtk.MenuNew()
	if err != nil {
		return nil, errors.Wrap(err, "gtk.MenuNew")
	}
	itemHelp.SetSubmenu(helpMenu)
	itemAbout, err := gtk.MenuItemNewWithLabel("About")
	if err != nil {
		return nil, errors.Wrap(err, "gtk.MenuItemNewWithLabel")
	}
	helpMenu.Add(itemAbout)
	itemAbout.Connect("activate", func() {
		if err = postAbout(mainWindow); err != nil {
			// TODO post errors in dialog
			log.Println(errors.Wrap(err, "menu.PostAbout"))
		}
	})

	return itemHelp, nil
}

func postAbout(parent *gtk.Window) error {
	dialog, err := gtk.AboutDialogNew()
	if err != nil {
		return errors.Wrap(err, "gtk.AboutDialogNew")
	}

	dialog.SetTransientFor(parent)
	dialog.SetDestroyWithParent(true)
	dialog.SetModal(true)
	dialog.SetAuthors([]string{"Nathan Bass", "Grim"})
	dialog.SetProgramName("Gordon")
	dialog.SetComments("A JIM client designed for Tiberious")
	dialog.SetVersion("0.0.1")
	dialog.SetCopyright("JustAnotherOrganization")
	dialog.SetLicenseType(gtk.LICENSE_CUSTOM)

	license := `Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`
	dialog.SetLicense(license)

	dialog.Present()
	return nil
}
