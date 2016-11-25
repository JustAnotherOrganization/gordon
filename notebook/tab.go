package notebook

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

// Tab provides access to setting a tab label externally.
type Tab struct {
	Label *gtk.Label
}

func newTab(str string) (*Tab, error) {
	if str == "" {
		str = "New Chat"
	}

	var err error
	tab := &Tab{}
	tab.Label, err = gtk.LabelNew(str)
	if err != nil {
		return nil, errors.Wrap(err, "gtk.LabelNew")
	}

	return tab, nil
}

// SetLabel allows us to reset the label of a Tab externally.
func (t *Tab) SetLabel(str string) {
	t.Label.SetLabel(str)
}
