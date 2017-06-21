package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func main() {
	gtk.Init(nil)
	var e error
	var win *gtk.Window
	win, e = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if e != nil {
		log.Fatal("Unable to create window:", e)
	}
	win.SetTitle("Control de sincronización")
	win.Connect("destroy", gtk.MainQuit)

	// Create a new label widget to show in the window.
	var l *gtk.Label
	l, e = gtk.LabelNew("¡Hola gotk3!")
	if e != nil {
		log.Fatal("Unable to create label:", e)
	}

	// Add the label to the window.
	win.Add(l)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
