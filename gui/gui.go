package gui

import (
	"log"
	"strconv"
	"time"

	"teltonika2tak-cot/parser"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Gui struct {
	window    *gtk.Window
	listStore *gtk.ListStore
}

func NewGui() *Gui {
	gtk.Init(nil)

	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	window.SetTitle("GPS Tracker Data")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}
	window.Add(scrolledWindow)

	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}

	treeView, err := gtk.TreeViewNewWithModel(listStore)
	if err != nil {
		log.Fatal("Unable to create tree view:", err)
	}
	scrolledWindow.Add(treeView)

	for i, title := range []string{"IMEI", "Time", "Coordinates", "Direction", "Speed"} {
		cellRenderer, err := gtk.CellRendererTextNew()
		if err != nil {
			log.Fatal("Unable to create cell renderer:", err)
		}

		column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", i)
		if err != nil {
			log.Fatal("Unable to create column:", err)
		}

		treeView.AppendColumn(column)
	}

	return &Gui{
		window:    window,
		listStore: listStore,
	}
}

func (g *Gui) Update(data *parser.TeltParsedData) {
	iter := g.listStore.Append()

	err := g.listStore.Set(iter,
		[]int{0, 1, 2, 3, 4},
		[]interface{}{
			data.IMEI,
			time.Unix(int64(data.Time), 0).Format(time.RFC3339),
			strconv.FormatFloat(data.Lat, 'f', 6, 64) + ", " + strconv.FormatFloat(data.Lon, 'f', 6, 64),
			strconv.FormatFloat(data.Heading, 'f', 2, 64),
			strconv.FormatFloat(data.Speed, 'f', 2, 64),
		})
	if err != nil {
		log.Println("Failed to update GUI:", err)
	}
}

func (g *Gui) Run() {
	g.window.ShowAll()
	gtk.Main()
}
