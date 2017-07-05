package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/lamg/tesis"
	"strconv"
)

const (
	syncImgPth = "synchronizing-emblem"
	addImgPth  = "list-add"
	delImgPth  = "list-remove"

	// Contexts
	pending  = "Pendientes"
	proposed = "Propuestos"
	record   = "Registro"
)

var (
	// pl: Page Label
	// tl: Total Label
	pl, tl *gtk.Label
	// pb: Previous Button
	// nb: Next Button
	// ab: Contextual action Button
	pb, nb, ab *gtk.Button
	// sl: Show contextual elements ListBox
	sl *gtk.ListBox
)

func initWindow(dm tesis.DBManager, n error) (e error) {
	gtk.Init(nil)

	var bs []*gtk.Button
	var lbs []string
	var cbs []func(interface{})
	var n int
	n, lbs, cbs = 0,
		[]string{pending, proposed, record},
		[]func(interface{}){pendClk, propClk, recrClk}

	bs = make([]*gtk.Button, len(lbs))
	for e == nil && n != len(lbs) {
		var b *gtk.Button
		b, e = gtk.ButtonNewWithLabel(lbs[n])
		if e == nil {
			b.Connect("clicked", cbs[n])
			bs[n] = b
		}
	}

	var sc *gtk.Stack
	if e == nil {
		n = 0
		sc, e = gtk.StackNew()
	}

	for e == nil && n != len(bs) {
		sc.Add(bs[i])
	}
	if e == nil {
		setCtx(0)
	}
	// TODO StackSwitcher, ListBox update according context
	// buttons and labels updates according context
	var win *gtk.Window
	if e == nil {
		win, e = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	}
	if e == nil {
		win.Add(l)
		win.SetTitle("Control de sincronización")
		if n != nil {
			win.SetTitle(n.Error())
			e = n
		}
		win.SetIconName("emblem-synchronizing")
		win.Connect("destroy", gtk.MainQuit)
		win.SetDefaultSize(800, 600)
		win.ShowAll()
		gtk.Main()
	}
	// { win set up }
	return
}

func initbpe() (b *gtk.Box, e error) {
	var l *gtk.ListBox
	var a *gtk.ActionBar
	if e == nil {
		l, a, e = initlbpe()
	}

	if e == nil {
		b, e = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	}
	if e == nil {
		b.PackStart(l, true, true, 1)
		b.PackStart(a, true, true, 1)
	}
	// { Pendientes widget's initialized }
	return
}

func initbpr() (b *gtk.Box, e error) {
	var l *gtk.ListBox
	var a *gtk.ActionBar
	if e == nil {
		l, e = initlbpr()
	}
	if e == nil {
		a, e = gtk.ActionBarNew()
	}
	if e == nil {
		b, e = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	}
	if e == nil {
		b.PackStart(l, true, true, 1)
		b.PackStart(a, true, true, 1)
	}
	// { Propuestos widget's initialized }
	return
}

func initbre() (b *gtk.Box, e error) {
	var l *gtk.ListBox
	if e == nil {
		l, e = initlbre()
	}
	var a *gtk.ActionBar
	if e == nil {
		a, e = gtk.ActionBarNew()
	}
	if e == nil {
		b, e = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	}
	if e == nil {
		b.PackStart(l, true, true, 1)
		b.PackStart(a, true, true, 1)
	}
	// { Registro widget's initialized }
	return
}

func initlbpe(dm tesis.DBManager) (l *gtk.ListBox, a *gtk.ActionBar, e error) {
	if e == nil {
		l, e = gtk.ListBoxNew()
	}
	// pe0: page 0 of all Pending
	var pe0 *tesis.PageD
	if e == nil {
		pe0, e = dm.Pending(0)
	}
	// TODO a's components interact with l
	var lp, lt *gtk.Label
	var bp, bn *gtk.Button
	if e == nil {
		lp, e = gtk.LabelNew("")
	}
	if e == nil {
		lt, e = gtk.LabelNew("")
	}
	if e == nil {
		bp, e = gtk.ButtonNewFromIconName("go-previous",
			gtk.ICON_SIZE_BUTTON)
	}
	if e == nil {
		bn, e = gtk.ButtonNewFromIconName("go-next",
			gtk.ICON_SIZE_BUTTON)
	}
	if e == nil {
		e = addPage(pe0, l, lp, lt, bp, bn)
	}

}

func addPage(p *tesis.PageD, l *gtk.ListBox,
	lp, lt *gtk.Label, bp, bn *gtk.Button) (e error) {

	lp.SetText(strconv.Itoa(p.PageN))
	lt.SetText(strconv.Itoa(p.Total))
	bp.SetProperty("sensitive", p.PageN > 0)
	bn.SetProperty("sensitive", p.PageN < p.Total-1)
	var n int
	n = 0
	for e == nil && n != len(p.DiffP) {
		var b *gtk.Box
		b, e = boxp(p[n])
		if e == nil {
			l.Add(b)
		}
	}
	return
}

func boxp(d *tesis.Diff) (b *gtk.Box, e error) {
	b, e = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
	var s *gtk.CheckButton
	if e == nil {
		s, e = gtk.CheckButtonNew()
	}
	if e == nil {
		b.PackStart(s, true, true, 0)
	}
	// { component 0 from b added }
	var vbx *gtk.Box
	if e == nil {
		vbx, e = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	}
	var ipt, txt string
	if e == nil {
		if d.Exists && d.Mismatch {
			ipt, txt = syncImgPth, "Actualización"
		} else if d.Exists && !d.Mismatch {
			ipt, txt = delImgPth, "Eliminación"
		} else if !d.Exists {
			ipt, txt = addImgPth, "Adición"
		}
	}
	var img *gtk.Image
	if e == nil {
		img, e = gtk.ImageNewFromFile(ipt)
	}
	var lbl *gtk.Label
	if e == nil {
		lbl, e = gtk.LabelNew(txt)
	}
	if e == nil {
		vbx.PackStart(img, true, true, 0)
		vbx.PackStart(lbl, true, true, 0)
		b.PackStart(vbx, true, true, 0)
	}
	// { component 1 from b added }
	var hbx *gtk.Box
	if e == nil {
		hbx, e = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
	}
	var rl *tesis.DBRecord
	var n int
	rl, n = []*tesis.DBRecord{&d.DBRec, &d.LDAPRec}, 0
	for e == nil && n != len(rl) {
		var rbx *gtk.Box
		rbx, e = boxRecord(j)
		if e == nil {
			hbx.PackStart(rbx, true, true, 0)
		}
		n++
	}
	if e == nil {
		b.PackStart(hbx, true, true, 0)
		// { component 2 from b added  }
	}
	return
}

func boxRecord(r *tesis.DBRecord) (b *gtk.Box, e error) {
	// Id, IN, Name, Addr, Tel
	var fs, ls []string
	var n int
	fs, ls, n =
		[]string{r.Id, r.IN, r.Name, r.Addr, r.Tel},
		[]string{"Id en la base de datos", "Número de identidad",
			"Nombre", "Dirección", "Teléfono"},
		0
	for e == nil && n != len(fs) {
		var lb, vl *gtk.Label
		var hb *gtk.Box
		lb, e = gtk.LabelNew(fs[n])
		if e == nil {
			vl, e = gtk.LabelNew(ls[n])
		}
		if e == nil {
			hb, e = gtk.BoxNew(ORIENTATION_HORIZONTAL, 2)
		}
		if e == nil {
			hb.PackStart(lb, true, true, 0)
			hb.PackStart(vl, true, true, 0)
			b.PackStart(hb, true, true, 0)
		}
		n++
	}
	return
}

func initlbre() (l *gtk.ListBox, e error) {
	//TODO
	return
}

func initape() (a *gtk.ActionBar, e error) {
	//TODO
	a, e = gtk.ActionBarNew()
	// btpe: Pendientes Button (Proponer)
	var btpe *gtk.Button
	if e == nil {
		btpe, e = gtk.ButtonNewWithLabel("Proponer")
	}
	if e == nil {
		a.Add(btpe)
	}
	return
}

func initapr() (a *gtk.ActionBar, e error) {
	//TODO
	a, e = gtk.ActionBarNew()
	// btpr: Propuestos Button (Revertir)
	var btpr *gtk.Button
	if e == nil {
		btpr, e = gtk.ButtonNewWithLabel("Revertir")
	}
	if e == nil {
		a.Add(btpr)
	}
	return
}
