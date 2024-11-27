package album

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type album struct {
	widget.BaseWidget

	page  int
	image *canvas.Image
}

func NewAlbum() *album {
	return &album{
		page:  1,
		image: &canvas.Image{FillMode: canvas.ImageFillOriginal},
	}
}

func (a *album) downloadImage(url string) {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	file, err := ioutil.TempFile(os.TempDir(), "xkcd.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	a.image.File = file.Name()
	canvas.Refresh(a.image)
	fmt.Println("Image Refreshed ", a.page)
}

func (a *album) Paging() {
	endpoint := fmt.Sprintf("https://xkcd.com/%d/info.0.json", a.page)
	resp, err := http.Get(endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	type xkcdResp struct {
		Img string `json:"img"`
	}
	if resp.StatusCode == http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		x := &xkcdResp{}
		json.Unmarshal(data, x)
		go a.downloadImage(x.Img)
	} else {
		fmt.Println("Error getting ID", a.page, resp.Status, resp.StatusCode)
	}
}

func Show(win fyne.Window) fyne.CanvasObject {
	album := NewAlbum()

	lastPageBtn := widget.NewButton("<", func() {
		album.page--
		album.Paging()
	})
	nextPageBtn := widget.NewButton(">", func() {
		album.page++
		album.Paging()
	})

	border := container.NewBorder(nil, nil, lastPageBtn, nextPageBtn, album.image)
	go album.Paging()

	return border
}
