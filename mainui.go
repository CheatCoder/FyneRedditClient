package main

import (
	"fmt"
	"image/color"
	"net/url"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"github.com/CheatCoder/geddit"
)

type mwin struct {
	win     fyne.Window
	vbox    *fyne.Container
	scroll  *container.Scroll
	sidebar *widget.List

	nowsub string
	sort   geddit.PopularitySort
}

type Post struct {
	Title     string
	Text      string
	Thumbnail string
	Url       string
}

var m mwin

func mainwinstart() {
	m.sort = geddit.NewSubmissions
	app := fyne.CurrentApp()
	m.win = app.NewWindow("Reddit")
	sorting := fyne.NewMenuItem("Sorting", nil)
	menu := fyne.NewMenu("Subreddit", fyne.NewMenuItem("Unsubscribe", func() {
		defer func() {
			recover()
		}()
		o.Unsubscribe(m.nowsub)
		m.setsidebar()
	}),
		sorting)
	sorting.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("New", func() { m.sort = geddit.NewSubmissions; reload() }),
		fyne.NewMenuItem("Hot", func() { m.sort = geddit.HotSubmissions; reload() }),
		fyne.NewMenuItem("Top", func() { m.sort = geddit.TopSubmissions; reload() }),
		fyne.NewMenuItem("Rising", func() { m.sort = geddit.RisingSubmissions; reload() }),
	)

	usermenu := fyne.NewMenu("User", fyne.NewMenuItem("View Saved", func() {
		Savedviewer()
	}))
	m.win.SetMainMenu(
		fyne.NewMainMenu(
			fyne.NewMenu("Menü", fyne.NewMenuItem("Credits", func() {
				CreditsWindow(fyne.CurrentApp(), fyne.Size{Width: 300, Height: 500}).Show()
			})),
			menu,
			usermenu,
		),
	)

	m.setsidebar()

	m.vbox = container.New(layout.NewVBoxLayout())
	m.scroll = container.NewScroll(m.vbox)
	m.vbox.Add(widget.NewCard("Placeholder", "Substring", canvas.NewLine(color.Black)))
	front, err := o.Frontpage(m.sort, geddit.ListingOptions{Limit: 20})
	if err != nil {
		m.vbox.Add(widget.NewCard("Placeholder", "Substring", canvas.NewLine(color.Black)))
	} else {
		m.setPosts(front)
	}

	search := xwidget.NewCompletionEntry([]string{})
	search.Entry.PlaceHolder = "Search"
	search.Entry.OnChanged = func(s string) {
		search.HideCompletion()
		completed, err := o.Autocompelete(s, true, false)
		if err != nil {
			dialog.NewError(err, m.win)
			return
		}
		search.SetOptions(completed)
		search.ShowCompletion()
	}
	search.Entry.OnSubmitted = func(s string) {
		p, err := o.SubredditSubmissions(s, m.sort, geddit.ListingOptions{Limit: 20})
		defer func() {
			recover()
		}()
		if err != nil {
			dialog.NewError(err, m.win)
			return
		}
		m.nowsub = s
		m.setPosts(p)
	}

	m.win.SetContent(container.NewBorder(search, nil, m.sidebar, nil, m.scroll))

	m.win.ShowAndRun()
}

func (m *mwin) setsidebar() {
	sub, err := o.MySubreddits(100)
	if err != nil {
		return
	}
	subname := make([]string, 0)

	for _, v := range sub {
		subname = append(subname, v.Name)
	}

	sort.Strings(subname)

	m.sidebar = widget.NewList(
		func() int {
			return len(subname)
		},
		func() fyne.CanvasObject {
			dummy := strings.Repeat("*", len(subname))
			label := widget.NewLabel(dummy)
			label.Alignment = fyne.TextAlignCenter
			return label
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(subname[id])
		},
	)

	m.sidebar.OnSelected = func(id widget.ListItemID) {
		//fmt.Println(subname[id], "was clicked -> TODO: set content")
		posts, err := o.SubredditSubmissions(subname[id], m.sort, geddit.ListingOptions{Limit: 50})
		m.nowsub = subname[id]
		if err != nil {
			dialog.NewError(err, m.win)
		}
		m.setPosts(posts)
	}
}

func (m *mwin) setPosts(p []*geddit.Submission) {
	cbox := container.NewVBox()
	m.vbox.Remove(m.vbox.Objects[0])
	m.scroll.ScrollToTop()
	for _, v := range p {
		go func(v *geddit.Submission) {
			res, err := fyne.LoadResourceFromURLString(v.ThumbnailURL)
			if err != nil {
				fmt.Println("Returned")
				return
			}
			img := canvas.NewImageFromResource(res)
			img.FillMode = canvas.ImageFillOriginal
			var label *widget.Label

			opbtn := widget.NewButton("Open", func() {
				purl, _ := url.Parse(v.URL)
				fyne.CurrentApp().OpenURL(purl)
			})
			savebtn := widget.NewButton("Save", func() {
				err := o.Save(v, "")
				if err != nil {
					fmt.Println(err)
					return
				}
				label.SetText("Saved")
			})
			if v.IsSaved {
				label = widget.NewLabel("Is Saved")
				savebtn.Disable()
			} else {
				label = widget.NewLabel("Is not Saved")
			}
			hbox := container.NewHBox(opbtn, savebtn)
			box := container.NewVBox(img, label, hbox)
			card := widget.NewCard(v.Title, v.Subreddit, box)
			cbox.Add(card)
		}(v)
	}
	m.vbox.Add(cbox)
}

func reload() {
	var posts []*geddit.Submission
	var err error
	if m.nowsub == "" {
		posts, err = o.Frontpage(m.sort, geddit.ListingOptions{Limit: 20})
		if err != nil {
			return
		}
	} else {
		posts, err = o.SubredditSubmissions(m.nowsub, m.sort, geddit.ListingOptions{Limit: 50})
		if err != nil {
			return
		}
	}
	m.setPosts(posts)
}
