package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)
	var e error

	// lpe: Pendientes Label
	// lpr: Propuestos Label
	// lre: Registro   Label
	var lpe, lpr, lre *gtk.Label
	// mpe: Pendientes ListStore
	// mpr: Propuestos ListStore
	var mpe, mpr *gtk.ListStore
	// lspe: Pendientes List
	// lspr: Propuestos List
	var lspe, lspr *gtk.TreeView
	// btpe: Pendientes Button (Proponer)
	// btpr: Propuestos Button (Revertir)
	var btpe, btpr *gtk.Button
	// ape: Pendientes ActionBar
	// apr: Propuestos ActionBar
	// are: Registro   ActionBar
	var ape, apr, are *gtk.ActionBar
	// bpe: Pendientes Box
	// bpr: Propuestos Box
	// bre: Registro	 Box
	var bpe, bpr, bre *gtk.Box
	if e == nil {
		lpe, e = gtk.LabelNew("Pendientes")
	}
	if e == nil {
		mpe, e = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	}
	if e == nil {
		lspe, e = gtk.TreeViewNewWithModel(mpe)
	}
	var cpe *gtk.TreeViewColumn
	if e == nil {
		//cell renderer
		cpe, e = gtk.TreeViewColumnNewWithAttribute("coco")
	}
	if e == nil {
		lspe.AppendColumn(cpe)
	}
	if e == nil {
		btpe, e = gtk.ButtonNewWithLabel("Proponer")
	}
	if e == nil {
		ape, e = gtk.ActionBarNew()
	}
	if e == nil {
		bpe, e = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	}
	// { Pendientes widget's initialized }
	if e == nil {
		lpr, e = gtk.LabelNew("Propuestos")
	}
	if e == nil {
		mpr, e = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	}
	if e == nil {
		lspr, e = gtk.TreeViewNewWithModel(mpr)
	}
	if e == nil {
		btpr, e = gtk.ButtonNewWithLabel("Revertir")
	}
	if e == nil {
		apr, e = gtk.ActionBarNew()
	}
	if e == nil {
		bpr, e = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	}
	// { Propuestos widget's initialized }
	if e == nil {
		lre, e = gtk.LabelNew("Registro")
	}
	if e == nil {
		are, e = gtk.ActionBarNew()
	}
	if e == nil {
		bre, e = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	}
	// { Registro widget's initialized }

	var l *gtk.Notebook
	if e == nil {
		l, e = gtk.NotebookNew()
	}
	if e == nil {
		bpe.PackStart(lspe, true, true, 1)
		ape.PackStart(btpe)
		bpe.PackStart(ape, true, true, 1)
		l.AppendPage(bpe, lpe)
		bpr.PackStart(lspr, true, true, 1)
		apr.PackStart(btpr)
		bpr.PackStart(apr, true, true, 1)
		l.AppendPage(bpr, lpr)
		bre.PackStart(are, true, true, 1)
		l.AppendPage(bre, lre)
	}
	// { l is NoteBook with Pendientes, Propuestos and Record }

	var win *gtk.Window
	if e == nil {
		win, e = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	}

	if e == nil {
		win.Add(l)
		win.SetTitle("Control de sincronización")
		win.SetIconName("emblem-synchronizing")
		win.Connect("destroy", gtk.MainQuit)
		win.SetDefaultSize(800, 600)
		win.ShowAll()
		gtk.Main()
	}
}
