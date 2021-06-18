package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/CheatCoder/geddit"
)

var mainwin fyne.Window
var after string
var before string
var page = 0
var rwmu sync.RWMutex

func startfyne() {
	app.NewWithID("Reddit Saved View")
}

func Savedviewer() {
	app := fyne.CurrentApp()

	mainwin = app.NewWindow("Saved Posts")

	mainwin.Resize(fyne.Size{Height: 600, Width: 300})

	re, err := o.Me()
	if err != nil {
		os.Exit(-100)
	}

	con := content(page)
	mainwin.SetContent(con)
	dialog.NewInformation("Welcome", "You are Loged in as "+re.Name, mainwin).Show()

	mainwin.ShowAndRun()
}

func content(nspage int) *container.Scroll {
	list := getUserList(nspage)
	//con := container.New(layout.NewMaxLayout())

	vbox := container.NewVBox()

	scroll := container.NewVScroll(vbox)
	//con.Add(scroll)

	for _, v := range list {
		time.Sleep(1 * time.Millisecond)
		go func(v *geddit.Submission) {
			rwmu.Lock()
			defer rwmu.Unlock()
			ivbox := container.NewVBox()
			if strings.HasSuffix(v.URL, ".jpg") || strings.HasSuffix(v.URL, ".png") {
				res, err := fyne.LoadResourceFromURLString(v.ThumbnailURL)
				if err != nil {
					return
				}
				img := canvas.NewImageFromResource(res)
				img.SetMinSize(fyne.Size{Height: 350})
				img.FillMode = canvas.ImageFillStretch
				img.ScaleMode = canvas.ImageScaleSmooth

				obtn := widget.NewButton("Open", func() {
					go func(upurl string) {
						purl, _ := url.Parse(upurl)
						fyne.CurrentApp().OpenURL(purl)
					}(v.URL)
				})
				ivbox.Add(img)
				ivbox.Add(obtn)
				vbox.Add(ivbox)
				vbox.Add(layout.NewSpacer())

			} else {
				fmt.Println(v.URL)
			}
		}(v)
	}

	before = list[0].FullID
	after = list[len(list)-1].FullID

	if page == 0 {
		nextbtn := widget.NewButton("Next", func() {
			page += 1
			lbox := content(page + 1)
			mainwin.SetContent(lbox)
		})
		vbox.Add(container.NewVBox(layout.NewSpacer(), nextbtn))
	} else {
		prevbtn := widget.NewButton("Previous", func() {
			page -= 1
			lbox := content(page + 1)
			mainwin.SetContent(lbox)
		})
		nextbtn := widget.NewButton("Next", func() {
			page += 1
			lbox := content(page + 1)
			mainwin.SetContent(lbox)
		})
		vbox.Add(container.NewVBox(prevbtn, layout.NewSpacer(), nextbtn))
	}

	return scroll
}
