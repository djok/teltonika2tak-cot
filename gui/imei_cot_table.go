package gui

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type IMEICoTTable struct {
	window    *gtk.Window
	listStore *gtk.ListStore
}

func NewIMEICoTTable() *IMEICoTTable {
	gtk.Init(nil)

	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	window.SetTitle("IMEI to CoT Conversion Table")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}
	window.Add(scrolledWindow)

	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}

	treeView, err := gtk.TreeViewNewWithModel(listStore)
	if err != nil {
		log.Fatal("Unable to create tree view:", err)
	}
	scrolledWindow.Add(treeView)

	for i, title := range []string{"IMEI", "CoT ID"} {
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

	return &IMEICoTTable{
		window:    window,
		listStore: listStore,
	}
}

func (t *IMEICoTTable) Update(imei string, cotID string) {
	iter := t.listStore.Append()

	err := t.listStore.Set(iter,
		[]int{0, 1},
		[]interface{}{
			imei,
			cotID,
		})
	if err != nil {
		log.Println("Failed to update IMEI to CoT table:", err)
	}
}

func (t *IMEICoTTable) Run() {
	t.window.ShowAll()
	gtk.Main()
}
