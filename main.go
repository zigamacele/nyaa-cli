package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func openFile(filename string) error {
	cmd := exec.Command("open", filename) // windows:  use "start" instead of "open"
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(filepath string, url string) (err error) {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	openFile(filepath)
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "nyaa-cli")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

type Torrent struct {
	category      string
	view_link     string
	name          string
	download_link string
	magnet_link   string
	size          string
	date          string
	seeders       string
	leechers      string
	downloads     string
}

func main() {

	//GUI

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	//WEB SCRAPING

	c := colly.NewCollector()

	var Torrents []Torrent

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		var torrent []string
		e.ForEach("td", func(_ int, td *colly.HTMLElement) {
			td.ForEach("a", func(_ int, a *colly.HTMLElement) {
				link := a.Attr("href")
				if strings.Contains(link, "/?c=") {
					link = strings.ReplaceAll(link, "/?c=", "")
				}
				if strings.Contains(link, "#comments") {
					return
				}

				//fmt.Printf("Link: %s\n", link)
				torrent = append(torrent, link)
			})

			var tdText = strings.ReplaceAll(td.Text, "\n", "")
			tdText = strings.ReplaceAll(tdText, "\t", "")
			if tdText != "" {
				torrent = append(torrent, tdText)
				//fmt.Printf("%q\n", tdText)
			}

		})

		Torrents = append(Torrents, Torrent{
			category:      torrent[0],
			view_link:     torrent[1],
			name:          torrent[2],
			download_link: torrent[3],
			magnet_link:   torrent[4],
			size:          torrent[5],
			date:          torrent[6],
			seeders:       torrent[7],
			leechers:      torrent[8],
			downloads:     torrent[9],
		})

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://nyaa.si")

	//TORRENT DOWNLOAD

	var url string = "https://nyaa.si" + Torrents[0].download_link
	var fileName string = Torrents[0].name + ".torrent"

	downloadFile(fileName, url)
	openFile(fileName)

}
