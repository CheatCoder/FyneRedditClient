package main

import (
	"fmt"
	"image/color"
	"net/url"
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

var mwin fyne.Window
var vbox *fyne.Container
var scroll *container.Scroll

var nowsub string

type Post struct {
	Title     string
	Text      string
	Thumbnail string
	Url       string
}

func mainwinstart() {
	app := fyne.CurrentApp()

	mwin = app.NewWindow("Reddit")
	menu := fyne.NewMenu("Subreddit", fyne.NewMenuItem("Unsubscribe", func() {
		defer func() {
			recover()
		}()
		o.Unsubscribe(nowsub)
	}))
	mwin.SetMainMenu(
		fyne.NewMainMenu(
			fyne.NewMenu("Menü", fyne.NewMenuItem("Credits", func() {
				CreditsWindow(fyne.CurrentApp(), fyne.Size{Width: 300, Height: 500}).Show()
			})),
			menu,
		),
	)

	nav := sidebar()

	vbox = container.New(layout.NewVBoxLayout())
	scroll = container.NewScroll(vbox)
	vbox.Add(widget.NewCard("Placeholder", "Substring", canvas.NewLine(color.Black)))
	front, err := o.Frontpage(geddit.NewSubmissions, geddit.ListingOptions{Limit: 20})
	if err != nil {
		vbox.Add(widget.NewCard("Placeholder", "Substring", canvas.NewLine(color.Black)))
	} else {
		setPosts(front)
	}

	search := xwidget.NewCompletionEntry([]string{})
	search.Entry.PlaceHolder = "Search"
	search.Entry.OnChanged = func(s string) {
		search.HideCompletion()
		completed, err := o.Autocompelete(s, true, false)
		if err != nil {
			dialog.NewError(err, mwin)
			return
		}
		search.SetOptions(completed)
		search.ShowCompletion()
	}
	search.Entry.OnSubmitted = func(s string) {
		p, err := o.SubredditSubmissions(s, geddit.NewSubmissions, geddit.ListingOptions{Limit: 20})
		defer func() {
			recover()
		}()
		if err != nil {
			dialog.NewError(err, mwin)
			return
		}
		setPosts(p)
	}

	mwin.SetContent(container.NewBorder(search, nil, nav, nil, scroll))

	mwin.ShowAndRun()
}

func sidebar() *widget.List {
	sub, err := o.MySubreddits()
	if err != nil {
		return nil
	}
	subname := make([]string, 0)

	for _, v := range sub {
		subname = append(subname, v.Name)
	}

	list := widget.NewList(
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

	list.OnSelected = func(id widget.ListItemID) {
		//fmt.Println(subname[id], "was clicked -> TODO: set content")
		posts, err := o.SubredditSubmissions(subname[id], geddit.NewSubmissions, geddit.ListingOptions{Limit: 50})
		nowsub = subname[id]
		if err != nil {
			dialog.NewError(err, mwin)
		}
		setPosts(posts)
	}

	return list
}

func setPosts(p []*geddit.Submission) {
	cbox := container.NewVBox()
	vbox.Remove(vbox.Objects[0])
	scroll.ScrollToTop()
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
	vbox.Add(cbox)
}
